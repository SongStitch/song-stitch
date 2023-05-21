package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	jpgFileType      = ".jpg"
	imageCoverWidth  = 300
	imageCoverHeight = 300
)

var fallbackImage image.Image

func init() {
	var err error
	fallbackImage, err = ReadImage("./fallback.jpg")
	if err != nil {
		log.Fatal(err)
	}
}

func ReadImage(path string) (image.Image, error) {
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

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255} // black
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func create_collage(albums []Album, rows int, columns int) (image.Image, error) {

	// create a new blank image with dimensions to fit all the images
	collage := imaging.New(imageCoverWidth*columns, imageCoverHeight*rows, image.Transparent)

	// add each image to the collage
	for i, album := range albums {
		imgRGBA := image.NewRGBA(album.Image.Bounds())
		draw.Draw(imgRGBA, imgRGBA.Bounds(), album.Image, album.Image.Bounds().Min, draw.Src)
		addLabel(imgRGBA, 10, 10, album.Artist)
		addLabel(imgRGBA, 10, 25, album.Name)
		x := (i % columns) * imageCoverWidth
		y := (i / columns) * imageCoverHeight
		collage = imaging.Paste(collage, imgRGBA, image.Pt(x, y))
	}

	/* No need to save? */
	// save the collage to file
	//	err := imaging.Save(collage, "collage.jpg")
	//	if err != nil {
	//		fmt.Println(err)
	//				return
	//	}

	return collage, nil
}
