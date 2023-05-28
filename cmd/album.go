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

type TopResponse interface {
	GetName() string
	GetPlaycount() string
	GetImageUrl() string
	GetImage() *image.Image
	SetImage(*image.Image)
}

func (a *Album) GetName() string {
	return a.Name
}

func (a *Album) GetPlaycount() string {
	return a.Playcount
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

func DownloadImage(a TopResponse) error {
	if len(a.GetImageUrl()) == 0 {
		// Skip album art if it doesn't exist
		return nil
	}
	resp, err := http.Get(a.GetImageUrl())
	if err != nil {
		return err
	}
	ioBody := resp.Body
	defer resp.Body.Close()

	extension, err := getExtension(a.GetImageUrl())
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

func downloadImagesForAlbums(albums []Album) error {
	var wg sync.WaitGroup
	wg.Add(len(albums))

	for i := range albums {
		album := &albums[i]
		// download each image in a separate goroutine
		go func(album *Album) {
			defer wg.Done()
			err := DownloadImage(album)
			if err != nil {
				log.Println(err)
			}
		}(album)
	}

	// wait for all downloads to finish
	wg.Wait()

	return nil
}
