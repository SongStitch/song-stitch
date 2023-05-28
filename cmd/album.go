package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	jpgFileType = ".jpg"
	gifFileType = ".gif"
)

type Album struct {
	Name      string
	Artist    string
	Playcount string
	ImageUrl  string
	Image     image.Image
}

type Artist struct {
	Name      string
	Playcount string
	Image     image.Image
	ImageUrl  string
}

type Downloadable interface {
	GetImageUrl() string
	GetImage() *image.Image
	SetImage(*image.Image)
}

func (a *Album) GetImageUrl() string {
	return a.ImageUrl
}

func (a *Album) GetImage() *image.Image {

	return &a.Image
}

func (a *Album) SetImage(img *image.Image) {
	a.Image = *img
}

func downloadImage[T Downloadable](a T) error {
	url := a.GetImageUrl()
	if len(url) == 0 {
		// Skip album art if it doesn't exist
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	ioBody := resp.Body
	defer resp.Body.Close()

	extension, err := getExtension(url)
	if err != nil {
		return err
	}

	if strings.ToLower(extension) == jpgFileType {
		img, err := jpeg.Decode(ioBody)
		if err != nil {
			return err
		}
		a.SetImage(&img)
		return err
	} else if strings.ToLower(extension) == gifFileType {
		img, err := gif.Decode(ioBody)
		if err != nil {
			return err
		}
		a.SetImage(&img)
		return err
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		img, _, err := image.Decode(bytes.NewReader(body))
		a.SetImage(&img)
		return err

	}
}

func downloadImages[T Downloadable](entities []T) error {

	var wg sync.WaitGroup
	wg.Add(len(entities))

	for i := range entities {
		entity := &entities[i]
		// download each image in a separate goroutine
		func(entity *T) {
			defer wg.Done()
			err := downloadImage(*entity)
			if err != nil {
				log.Println(err)
			}
		}(entity)
	}

	// wait for all downloads to finish
	wg.Wait()

	return nil
}
