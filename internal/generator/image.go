package generator

import (
	"bytes"
	"context"
	"fmt"

	"image"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/fogleman/gg"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"

	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
)

type DisplayOptions struct {
	ArtistName     bool
	AlbumName      bool
	TrackName      bool
	PlayCount      bool
	Compress       bool
	Resize         bool
	Width          uint
	Height         uint
	FontSize       float64
	BoldFont       bool
	Rows           int
	Columns        int
	ImageDimension int
	Webp           bool
	TextLocation   constants.TextLocation
}

const (
	fontFileRegular    = "./assets/NotoSans-Regular.ttf"
	fontFileBold       = "./assets/NotoSans-Bold.ttf"
	compressionQuality = 70
)

func getExtension(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	// Split the path component of the URL into a slice of path elements
	pathElements := strings.Split(parsedURL.Path, "/")

	// The last element of the path should be the filename
	fileName := pathElements[len(pathElements)-1]

	// Extract the file extension from the filename
	ext := filepath.Ext(fileName)
	return ext, nil
}

func getTextDrawLocation(dc *gg.Context, text string, x float64, y float64, displayOptions DisplayOptions) (float64, float64) {
	width, height := dc.MeasureString(text)
	imageWidth := float64(300 - 20)
	imageHeight := float64(300 - 20)
	textLength, _ := dc.MeasureString(text)
	fmt.Println(displayOptions.TextLocation)
	switch displayOptions.TextLocation {
	case constants.TOP_LEFT:
		return x, y
	case constants.TOP_CENTRE:
		return x + imageWidth/2 - width/2, y
	case constants.TOP_RIGHT:
		return x + imageWidth - textLength, y
	case constants.BOTTOM_LEFT:
		return x, y + imageHeight - height
	case constants.BOTTOM_CENTRE:
		return x + imageWidth/2 - width/2, y + imageHeight - height
	case constants.BOTTOM_RIGHT:
		return x + imageWidth - width, y + imageHeight - height
	default:
		return x, y
	}
}

func drawText(dc *gg.Context, text string, x float64, y float64, displayOptions DisplayOptions) float64 {
	x, y = getTextDrawLocation(dc, text, x, y, displayOptions)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, (x)+1, (y)+1, 0, 0)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, (x), (y), 0, 0)
	return (3 + displayOptions.FontSize)
}

func placeText[T Drawable](dc *gg.Context, drawable T, displayOptions DisplayOptions, x float64, y float64) {
	parameters := drawable.GetParameters()
	textLocation := (8 + displayOptions.FontSize)
	if val, ok := parameters["track"]; ok && displayOptions.TrackName && len(val) > 0 {
		// Add shadow
		textLocation += drawText(dc, val, x+10, y+textLocation, displayOptions)
	}
	if val, ok := parameters["artist"]; ok && displayOptions.ArtistName && len(val) > 0 {
		// Add shadow
		textLocation += drawText(dc, val, x+10, y+textLocation, displayOptions)
	}
	if val, ok := parameters["album"]; ok && displayOptions.AlbumName && len(val) > 0 {
		textLocation += drawText(dc, val, x+10, y+textLocation, displayOptions)
	}
	if val, ok := parameters["playcount"]; ok && displayOptions.PlayCount && len(val) > 0 {
		textLocation += drawText(dc, val, x+10, y+textLocation, displayOptions)
	}
}

func resizeImage(ctx context.Context, img *image.Image, width uint, height uint) *image.Image {
	if width == 0 && height == 0 {
		zerolog.Ctx(ctx).Info().Msg("Unable to resize image, both width and height are 0")
		return img
	} else if int(width) == (*img).Bounds().Dx() && int(height) == (*img).Bounds().Dy() {
		return img
	} else if height == 0 {

		height = uint(float64(width) * float64((*img).Bounds().Dy()) / float64((*img).Bounds().Dx()))
	} else if width == 0 {
		width = uint(float64(height) * float64((*img).Bounds().Dx()) / float64((*img).Bounds().Dy()))
	}
	result := resize.Resize(width, height, *img, resize.Lanczos3)
	return &result
}

func webpEncode(buf *bytes.Buffer, collage *image.Image, quality float32) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return err
	}

	err = webp.Encode(buf, *collage, options)
	return err
}

func CreateCollage[T Drawable](ctx context.Context, collageElements []T, displayOptions DisplayOptions) (*image.Image, *bytes.Buffer, error) {
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
		img := collageElement.GetImage()
		if *img != nil {
			img = resizeImage(ctx, img, uint(displayOptions.ImageDimension), uint(displayOptions.ImageDimension))
			dc.DrawImage(*img, x, y)
		}
		placeText(dc, collageElement, displayOptions, float64(x), float64(y))
	}
	collage := dc.Image()

	if displayOptions.Resize {
		collage = *resizeImage(ctx, &collage, displayOptions.Width, displayOptions.Height)
	}

	collageBuffer := new(bytes.Buffer)

	if displayOptions.Webp {
		err := webpEncode(collageBuffer, &collage, compressionQuality)
		if err != nil {
			logger.Err(err).Msg("Unable to create Webp image")
		}
	}

	logger.Info().Dur("duration", time.Since(start)).Int("rows", displayOptions.Rows).Int("columns", displayOptions.Columns).Msg("Collage created")
	return &collage, collageBuffer, nil
}
