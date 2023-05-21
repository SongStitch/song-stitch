package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
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

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	ioBody := resp.Body
	defer resp.Body.Close()

	extension, err := getExtension(url)
	if err != nil {
		return nil, err
	}

	fmt.Println(extension)
	if strings.ToLower(extension) == jpgFileType {
		img, err := jpeg.Decode(ioBody)
		if err != nil {
			fmt.Println(err)
		}
		return img, err
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		img, _, err := image.Decode(bytes.NewReader(body))
		return img, err
	}
}

func downloadImages(imageUrls []string) ([]image.Image, error) {
	var wg sync.WaitGroup
	wg.Add(len(imageUrls))

	fmt.Println(imageUrls)
	images := make([]image.Image, len(imageUrls))

	var mux sync.Mutex

	for i, url := range imageUrls {
		// download each image in a separate goroutine
		go func(i int, url string) {
			defer wg.Done()
			img, err := downloadImage(url)
			if err != nil {
				fmt.Println("Error downloading image:", err)
				img = fallbackImage
				//	return
			}
			// Protect writing to the images slice with a mutex
			mux.Lock()
			images[i] = img
			mux.Unlock()
		}(i, url)
	}

	// wait for all downloads to finish
	wg.Wait()

	fmt.Println(images)
	return images, nil
}

func create_collage(images []image.Image, rows int, columns int) (image.Image, error) {

	// create a new blank image with dimensions to fit all the images
	collage := imaging.New(imageCoverWidth*columns, imageCoverHeight*rows, image.Transparent)

	// add each image to the collage
	for i, img := range images {
		x := (i % columns) * imageCoverWidth
		y := (i / columns) * imageCoverHeight
		collage = imaging.Paste(collage, img, image.Pt(x, y))
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
