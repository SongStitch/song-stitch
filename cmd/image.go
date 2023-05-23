package main

import (
	"fmt"
	"image"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

type DisplayOptions struct {
	ArtistName bool
	AlbumName  bool
	PlayCount  bool
}

const (
	jpgFileType      = ".jpg"
	imageCoverWidth  = 300
	imageCoverHeight = 300
	fontfile         = "./assets/Hack-Regular.ttf"
	size             = 12
	dpi              = 72
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

func placeText(dc *gg.Context, album Album, displayOptions DisplayOptions, x int, y int) {
	i := 0
	if displayOptions.ArtistName {
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(album.Artist, float64(x+10), float64(y+textLocation[i]), 0, 0)
		i++
	}
	if displayOptions.AlbumName {
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(album.Name, float64(x+10), float64(y+textLocation[i]), 0, 0)
		i++
	}
	if displayOptions.PlayCount {
		if len(album.Playcount) > 0 {
			dc.SetRGB(1, 1, 1)
			dc.DrawStringAnchored(fmt.Sprintf("Plays: %s", album.Playcount), float64(x+10), float64(y+textLocation[i]), 0, 0)
		}
	}
}

func createCollage(albums []Album, rows int, columns int, displayOptions DisplayOptions) (image.Image, error) {

	dc := gg.NewContext(imageCoverWidth*columns, imageCoverHeight*rows)
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace(fontfile, 12)

	for i, album := range albums {
		x := (i % columns) * imageCoverWidth
		y := (i / columns) * imageCoverHeight
		if album.Image != nil {
			dc.DrawImage(album.Image, x, y)
		}
		placeText(dc, album, displayOptions, x, y)
	}
	collage := dc.Image()

	return collage, nil

}
