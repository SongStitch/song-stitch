package collages

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"slices"
	"time"

	"github.com/SongStitch/go-webp/encoder"
	"github.com/SongStitch/go-webp/webp"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/fogleman/gg"

	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
)

type DisplayOptions struct {
	TextLocation   lastfm.TextLocation
	Height         uint
	Width          uint
	ImageDimension int
	Columns        int
	Rows           int
	FontSize       float64
	PlayCount      bool
	Resize         bool
	BoldFont       bool
	Grayscale      bool
	Compress       bool
	ArtistName     bool
	TrackName      bool
	Webp           bool
	AlbumName      bool
}

type CollageElement struct {
	Image      image.Image
	Parameters map[string]string
}

const (
	fontFileRegular    = "./assets/NotoSans-Regular.ttf"
	fontFileBold       = "./assets/NotoSans-Bold.ttf"
	compressionQuality = 70
)

func getTextOffset(dc *gg.Context, text string, displayOptions DisplayOptions) (float64, float64) {
	width, height := dc.MeasureString(text)
	imageSize := float64(dc.Width()/displayOptions.Columns - 20)
	switch displayOptions.TextLocation {
	case lastfm.LocationTopLeft:
		return 0, 0
	case lastfm.LocationTopCentre:
		return imageSize/2 - width/2, 0
	case lastfm.LocationTopRight:
		return imageSize - width, 0
	case lastfm.LocationBottomLeft:
		return 0, imageSize - height
	case lastfm.LocationBottomCentre:
		return imageSize/2 - width/2, imageSize - height
	case lastfm.LocationBottomRight:
		return imageSize - width, imageSize - height
	default:
		return 0, 0
	}
}

func drawText(
	dc *gg.Context,
	text string,
	x float64,
	y float64,
	displayOptions DisplayOptions,
) float64 {
	x_offset, y_offset := getTextOffset(dc, text, displayOptions)
	x, y = x+x_offset, y+y_offset
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, (x)+1, (y)+1, 0, 0)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, (x), (y), 0, 0)
	return (3 + displayOptions.FontSize)
}

func placeText(
	dc *gg.Context,
	drawable CollageElement,
	displayOptions DisplayOptions,
	x float64,
	y float64,
) {
	parameters := drawable.Parameters
	textToDraw := []string{}
	if val, ok := parameters["track"]; ok && displayOptions.TrackName && len(val) > 0 {
		textToDraw = append(textToDraw, val)
	}
	if val, ok := parameters["artist"]; ok && displayOptions.ArtistName && len(val) > 0 {
		textToDraw = append(textToDraw, val)
	}
	if val, ok := parameters["album"]; ok && displayOptions.AlbumName && len(val) > 0 {
		textToDraw = append(textToDraw, val)
	}
	if val, ok := parameters["playcount"]; ok && displayOptions.PlayCount && len(val) > 0 {
		textToDraw = append(textToDraw, val)
	}

	if !displayOptions.TextLocation.IsTop() {
		slices.Reverse(textToDraw)

	}
	textLocation := (8 + displayOptions.FontSize)
	for _, text := range textToDraw {
		newOffset := drawText(dc, text, x+10, y+textLocation, displayOptions)
		if displayOptions.TextLocation.IsTop() {
			textLocation += newOffset
		} else {
			textLocation -= newOffset
		}
	}
}

func resizeImage(ctx context.Context, img image.Image, width uint, height uint) image.Image {
	if width == 0 && height == 0 {
		zerolog.Ctx(ctx).Info().Msg("Unable to resize image, both width and height are 0")
		return img
	} else if int(width) == (img).Bounds().Dx() && int(height) == (img).Bounds().Dy() /* #nosec G115 */ {
		return img
	} else if height == 0 {
		height = uint(float64(width) * float64((img).Bounds().Dy()) / float64((img).Bounds().Dx()))
	} else if width == 0 {
		width = uint(float64(height) * float64((img).Bounds().Dx()) / float64((img).Bounds().Dy()))
	}
	result := resize.Resize(width, height, img, resize.Lanczos3)
	return result
}

func webpEncode(buf *bytes.Buffer, collage image.Image, quality float32) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return err
	}
	options.LowMemory = true

	err = webp.Encode(buf, collage, options)
	return err
}

func convertToGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	draw.Draw(grayImg, grayImg.Bounds(), img, img.Bounds().Min, draw.Src)
	return grayImg
}

func CreateCollage(
	ctx context.Context,
	collageElements []CollageElement,
	displayOptions DisplayOptions,
) (image.Image, *bytes.Buffer, error) {
	start := time.Now()
	logger := zerolog.Ctx(ctx)

	collageWidth := displayOptions.ImageDimension * displayOptions.Columns
	collageHeight := displayOptions.ImageDimension * displayOptions.Rows
	dc := gg.NewContext(collageWidth, collageHeight)
	dc.SetRGB(0, 0, 0)
	fontFile := fontFileRegular
	if displayOptions.BoldFont {
		fontFile = fontFileBold
	}
	dc.LoadFontFace(fontFile, displayOptions.FontSize)

	for i, collageElement := range collageElements {
		x := (i % displayOptions.Columns) * displayOptions.ImageDimension
		y := (i / displayOptions.Columns) * displayOptions.ImageDimension
		img := collageElement.Image
		if img != nil {
			img = resizeImage(
				ctx,
				img,
				uint(displayOptions.ImageDimension), // #nosec G115
				uint(displayOptions.ImageDimension), // #nosec G115
			)
			dc.DrawImage(img, x, y)
		}
		placeText(dc, collageElement, displayOptions, float64(x), float64(y))
	}
	collage := dc.Image()

	if displayOptions.Resize {
		collage = resizeImage(ctx, collage, displayOptions.Width, displayOptions.Height)
	}

	if displayOptions.Grayscale {
		collage = convertToGrayscale(collage)
	}

	collageBuffer := new(bytes.Buffer)

	if displayOptions.Webp && !displayOptions.Grayscale {
		logger.Info().Msg("Converting to Webp image")
		err := webpEncode(collageBuffer, collage, compressionQuality)
		if err != nil {
			logger.Err(err).Msg("Unable to create Webp image")
		}
	}

	logger.Info().
		Dur("duration", time.Since(start)).
		Int("rows", displayOptions.Rows).
		Int("columns", displayOptions.Columns).
		Msg("Collage created")
	return collage, collageBuffer, nil
}
