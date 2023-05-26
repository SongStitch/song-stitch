package main

import (
	"fmt"
	"image"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type DisplayOptions struct {
	ArtistName bool
	AlbumName  bool
	PlayCount  bool
	Resize     bool
	Width      uint
	Height     uint
}

const (
	fontFile = "./assets/Hack-Regular.ttf"
)

var textLocation = [3]int{20, 35, 50}

func readImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

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

func placeText(dc *gg.Context, album *Album, displayOptions DisplayOptions, x int, y int) {
	i := 0
	if displayOptions.ArtistName {
		// Add shadow
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(album.Artist, float64(x+10)+1, float64(y+textLocation[i])+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(album.Artist, float64(x+10), float64(y+textLocation[i]), 0, 0)
		i++
	}
	if displayOptions.AlbumName {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(album.Name, float64(x+10)+1, float64(y+textLocation[i])+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(album.Name, float64(x+10), float64(y+textLocation[i]), 0, 0)
		i++
	}
	if displayOptions.PlayCount && len(album.Playcount) > 0 {
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", album.Playcount), float64(x+10)+1, float64(y+textLocation[i])+1, 0, 0)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", album.Playcount), float64(x+10), float64(y+textLocation[i]), 0, 0)
	}
}

func resizeImage(img image.Image, width uint, height uint) image.Image {
	if width == 0 && height == 0 {
		log.Println("Unable to resize image, both width and height are 0")
		return img
	} else if height == 0 {
		height = uint(float64(width) * float64(img.Bounds().Dy()) / float64(img.Bounds().Dx()))
	} else if width == 0 {
		width = uint(float64(height) * float64(img.Bounds().Dx()) / float64(img.Bounds().Dy()))
	}
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func createCollage(albums []Album, rows int, columns int, imageDimension int, fontSize float64, displayOptions DisplayOptions) (image.Image, error) {

	collageWidth := imageDimension * columns
	collageHeight := imageDimension * rows
	dc := gg.NewContext(collageWidth, collageHeight)
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace(fontFile, fontSize)
	err := dc.LoadFontFace(fontFile, fontSize)
	if err != nil {
		panic(err)
	}

	for i, album := range albums {
		x := (i % columns) * imageDimension
		y := (i / columns) * imageDimension
		if album.Image != nil {
			dc.DrawImage(album.Image, x, y)
		}
		placeText(dc, &album, displayOptions, x, y)
	}
	collage := dc.Image()

	if displayOptions.Resize {
		collage = resizeImage(collage, displayOptions.Width, displayOptions.Height)
	}

	return collage, nil

}
