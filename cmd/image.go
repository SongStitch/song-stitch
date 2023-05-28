package main

import (
	"bytes"
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
	ArtistName bool
	AlbumName  bool
	PlayCount  bool
	Compress   bool
	Resize     bool
	Width      uint
	Height     uint
}

type Drawable interface {
	GetImage() *image.Image
	GetParameters() map[string]string
}

const (
	fontFile           = "./assets/NotoSans-Regular.ttf"
	compressionQuality = 70
)

var textLocation = [3]int{20, 35, 50}

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
	i := 0
	parameters := drawable.GetParameters()
	if displayOptions.ArtistName {
		// Add shadow
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(parameters["artist"], float64(x+10)+1, float64(y+textLocation[i])+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(parameters["artist"], float64(x+10), float64(y+textLocation[i]), 0, 0)
		i++
	}
	if displayOptions.AlbumName {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(parameters["album"], float64(x+10)+1, float64(y+textLocation[i])+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(parameters["album"], float64(x+10), float64(y+textLocation[i]), 0, 0)
		i++
	}
	playcount := parameters["playcount"]
	if displayOptions.PlayCount && len(playcount) > 0 {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", playcount), float64(x+10)+1, float64(y+textLocation[i])+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", playcount), float64(x+10), float64(y+textLocation[i]), 0, 0)
	}
}

func resizeImage(img *image.Image, width uint, height uint) image.Image {
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

func createCollage[T Drawable](albums []T, rows int, columns int, imageDimension int, fontSize float64, displayOptions DisplayOptions) (image.Image, error) {

	collageWidth := imageDimension * columns
	collageHeight := imageDimension * rows
	dc := gg.NewContext(collageWidth, collageHeight)
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace(fontFile, fontSize)

	for i, album := range albums {
		x := (i % columns) * imageDimension
		y := (i / columns) * imageDimension
		if *album.GetImage() != nil {
			dc.DrawImage(*album.GetImage(), x, y)
		}
		placeText(dc, album, displayOptions, x, y)
	}
	collage := dc.Image()

	if displayOptions.Resize {
		collage = resizeImage(&collage, displayOptions.Width, displayOptions.Height)
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
