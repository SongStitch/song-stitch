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

type Downloadable interface {
	GetImageUrl() string
	SetImage(*image.Image)
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

	log.Println("Downloading images for", len(entities), "entities")

	for i := range entities {
		entity := &entities[i]
		// download each image in a separate goroutine
		go func(entity *T) {
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
