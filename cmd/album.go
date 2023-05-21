package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Album struct {
	Name      string
	Artist    string
	Playcount int
	ImageUrl  string
	Image     image.Image
}

func (a *Album) DownloadImage() error {
	resp, err := http.Get(a.ImageUrl)
	if err != nil {
		return err
	}
	ioBody := resp.Body
	defer resp.Body.Close()

	extension, err := getExtension(a.ImageUrl)
	if err != nil {
		return err
	}

	fmt.Println(extension)
	if strings.ToLower(extension) == jpgFileType {
		img, err := jpeg.Decode(ioBody)
		if err != nil {
			fmt.Println(err)
		}
		a.Image = img
		return err
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
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
				fmt.Println("Error downloading image:", err)
				album.Image = fallbackImage
				//	return
			}
		}(album)
	}

	// wait for all downloads to finish
	wg.Wait()

	return nil
}
