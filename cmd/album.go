package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
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

func (a *Album) DownloadImage() error {
	if len(a.ImageUrl) == 0 {
		// Skip album art if it doesn't exist
		return nil
	}
	resp, err := http.Get(a.ImageUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	ioBody := resp.Body
	defer resp.Body.Close()

	extension, err := getExtension(a.ImageUrl)
	if err != nil {
		log.Println(err)
		return err
	}

	if strings.ToLower(extension) == jpgFileType {
		img, err := jpeg.Decode(ioBody)
		if err != nil {
			log.Println(err)
		}
		a.Image = img
		return err
	} else if strings.ToLower(extension) == gifFileType {
		img, err := gif.Decode(ioBody)
		if err != nil {
			log.Println(err)
		}
		a.Image = img
		return err
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}
		img, _, err := image.Decode(bytes.NewReader(body))
		a.Image = img
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
			err := album.DownloadImage()
			if err != nil {
				log.Println(err)
			}
		}(album)
	}

	// wait for all downloads to finish
	wg.Wait()

	return nil
}
