package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
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
	Rows           int
	Columns        int
	ImageDimension int
}

type Drawable interface {
	GetImage() *image.Image
	GetParameters() map[string]string
}

const (
	fontFile           = "./assets/NotoSans-Regular.ttf"
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

func placeText[T Drawable](dc *gg.Context, drawable T, displayOptions DisplayOptions, x int, y int) {
	parameters := drawable.GetParameters()
	textLocation := 8 + displayOptions.FontSize
	if val, ok := parameters["track"]; ok && displayOptions.TrackName && len(val) > 0 {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(val, float64(x+10)+1, float64(y)+textLocation+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(parameters["track"], float64(x+10), float64(y)+textLocation, 0, 0)
		textLocation += 3 + displayOptions.FontSize
	}
	if val, ok := parameters["artist"]; ok && displayOptions.ArtistName && len(val) > 0 {
		// Add shadow
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(val, float64(x+10)+1, float64(y)+textLocation+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(val, float64(x+10), float64(y)+textLocation, 0, 0)
		textLocation += 3 + displayOptions.FontSize
	}
	if val, ok := parameters["album"]; ok && displayOptions.AlbumName && len(val) > 0 {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(val, float64(x+10)+1, float64(y)+textLocation+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(val, float64(x+10), float64(y)+textLocation, 0, 0)
		textLocation += 3 + displayOptions.FontSize
	}
	if val, ok := parameters["playcount"]; ok && displayOptions.PlayCount && len(val) > 0 {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", val), float64(x+10)+1, float64(y)+textLocation+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", val), float64(x+10), float64(y)+textLocation, 0, 0)
		textLocation += 3 + displayOptions.FontSize
	}
}

func resizeImage(ctx context.Context, img *image.Image, width uint, height uint) image.Image {
	if width == 0 && height == 0 {
		log.Println("Unable to resize image, both width and height are 0")
		return *img
	} else if height == 0 {
		height = uint(float64(width) * float64((*img).Bounds().Dy()) / float64((*img).Bounds().Dx()))
	} else if width == 0 {
		width = uint(float64(height) * float64((*img).Bounds().Dx()) / float64((*img).Bounds().Dy()))
	}
	return resize.Resize(width, height, *img, resize.Lanczos3)
}

func compressImage(collage *image.Image, quality int) (image.Image, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, *collage, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(bytes.NewReader(buf.Bytes()))
}

func createCollage[T Drawable](ctx context.Context, albums []T, displayOptions DisplayOptions) (image.Image, error) {

	collageWidth := displayOptions.ImageDimension * displayOptions.Columns
	collageHeight := displayOptions.ImageDimension * displayOptions.Rows
	dc := gg.NewContext(collageWidth, collageHeight)
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace(fontFile, displayOptions.FontSize)

	for i, album := range albums {
		x := (i % displayOptions.Columns) * displayOptions.ImageDimension
		y := (i / displayOptions.Columns) * displayOptions.ImageDimension
		if *album.GetImage() != nil {
			dc.DrawImage(*album.GetImage(), x, y)
		}
		placeText(dc, album, displayOptions, x, y)
	}
	collage := dc.Image()

	if displayOptions.Resize {
		collage = resizeImage(ctx, &collage, displayOptions.Width, displayOptions.Height)
	}

	if displayOptions.Compress {
		collageCompressed, err := compressImage(&collage, compressionQuality)
		if err != nil {
			// Skip and just serve the non-compressed image
			log.Println(err)
		} else {
			collage = collageCompressed
		}
	}
	return collage, nil
}
