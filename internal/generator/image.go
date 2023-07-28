package generator

import (
	"bytes"
	"context"
	"image"
	"sync"
	"time"

	"github.com/SongStitch/go-webp/encoder"
	"github.com/SongStitch/go-webp/webp"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/fogleman/gg"

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

func getTextOffset(dc *gg.Context, text string, displayOptions DisplayOptions) (float64, float64) {
	width, height := dc.MeasureString(text)
	imageSize := float64(dc.Width() - 20)
	switch displayOptions.TextLocation {
	case constants.TOP_LEFT:
		return 0, 0
	case constants.TOP_CENTRE:
		return imageSize/2 - width/2, 0
	case constants.TOP_RIGHT:
		return imageSize - width, 0
	case constants.BOTTOM_LEFT:
		return 0, imageSize - height
	case constants.BOTTOM_CENTRE:
		return imageSize/2 - width/2, imageSize - height
	case constants.BOTTOM_RIGHT:
		return imageSize - width, imageSize - height
	default:
		return 0, 0
	}
}

func drawText(dc *gg.Context, text string, x float64, y float64, displayOptions DisplayOptions) float64 {
	x_offset, y_offset := getTextOffset(dc, text, displayOptions)
	x, y = x+x_offset, y+y_offset
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, (x)+1, (y)+1, 0, 0)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, (x), (y), 0, 0)
	return (3 + displayOptions.FontSize)
}

func placeText[T Drawable](dc *gg.Context, drawable T, displayOptions DisplayOptions, x float64, y float64) {
	parameters := drawable.GetParameters()
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
		reverse(textToDraw)
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

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
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

func WebpEncode(buf *bytes.Buffer, collage *image.Image) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, compressionQuality)
	if err != nil {
		return err
	}
	options.LowMemory = true

	err = webp.Encode(buf, *collage, options)
	return err
}

func CreateCollage[T Drawable](ctx context.Context, collageElements []T, displayOptions DisplayOptions) (*image.Image, error) {
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
			collageElement.ClearImage()
		}
		placeText(dc, collageElement, displayOptions, float64(x), float64(y))
	}
	collage := dc.Image()

	if displayOptions.Resize {
		collage = *resizeImage(ctx, &collage, displayOptions.Width, displayOptions.Height)
	}

	logger.Info().Dur("duration", time.Since(start)).Int("rows", displayOptions.Rows).Int("columns", displayOptions.Columns).Msg("Collage created")
	return &collage, nil
}

func CreateCollageEfficient[T Drawable](ctx context.Context, albums []T, displayOptions DisplayOptions) (*image.Image, error) {
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

	type albumImagePair struct {
		Album T
		Img   image.Image
	}

	// Create a channel to receive album-image pairs
	// resultChan := make(chan albumImagePair)

	// Use a wait group to keep track of goroutines
	var wg sync.WaitGroup
	maxconcurrent := 10
	sem := make(chan struct{}, maxconcurrent)
	images := make([]image.Image, maxconcurrent)

	// increment the wait group for each goroutine created
	wg.Add(len(albums))

	// create a goroutine for each album to fetch and process the image concurrently
	for i, album := range albums {
		x := (i % displayOptions.Columns) * displayOptions.ImageDimension
		y := (i / displayOptions.Columns) * displayOptions.ImageDimension
		sem <- struct{}{} // acquire a semaphore slot
		go func(album T, x int, y int, i int) {
			defer func() {
				<-sem // Release the semaphore slot when done
				wg.Done()
			}()

			img := images[i%maxconcurrent]
			err := GetImage(ctx, &img, album.GetImageUrl())
			if err != nil {
				logger.Error().Err(err).Msg("Failed to fetch image")
				return
			}

			if img != nil {
				dc.DrawImage(img, x, y)
				album.ClearImage()

				// Place text after drawing the image to avoid overlapping issues
				placeText(dc, album, displayOptions, float64(x), float64(y))
			}

			// Send the album-image pair to the channel
			// resultChan <- albumImagePair{Album: album, Img: img}
		}(album, x, y, i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the result channel after all goroutines have completed
	// close(resultChan)

	collage := dc.Image()

	if displayOptions.Resize {
		collage = *resizeImage(ctx, &collage, displayOptions.Width, displayOptions.Height)
	}

	logger.Info().Dur("duration", time.Since(start)).Int("rows", displayOptions.Rows).Int("columns", displayOptions.Columns).Msg("Collage created")
	return &collage, nil
}
