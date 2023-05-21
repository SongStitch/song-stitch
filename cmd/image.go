package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
)

const (
	jpgFileType = ".jpg"
)

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
		fmt.Println("JPG FILE")
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

	images := make([]image.Image, len(imageUrls))

	var mux sync.Mutex

	for _, url := range imageUrls {
		// download each image in a separate goroutine
		go func(url string) {
			defer wg.Done()
			img, err := downloadImage(url)
			if err != nil {
				fmt.Println("Error downloading image:", err)
				return
			}
			// Protect writing to the images slice with a mutex
			mux.Lock()
			images = append(images, img)
			mux.Unlock()
		}(url)
	}

	// wait for all downloads to finish
	wg.Wait()

	return images, nil
}

func create_collage(images []image.Image, width int, height int) (image.Image, error) {
	//	// determine collage grid size (assume square grid for simplicity)
	//	gridSize := int(math.Sqrt(float64(len(images))))
	//
	//	// dimensions of the collage based on the first image
	//	width := images[0].Bounds().Size().X
	//	height := images[0].Bounds().Size().Y
	//
	//	// create a new blank image with dimensions to fit all the images
	//	collage := imaging.New(width*gridSize, height*gridSize, image.Transparent)
	//
	//	// add each image to the collage
	//	for i, img := range images {
	//		x := (i % gridSize) * width
	//		y := (i / gridSize) * heighg@l
	//		collage = imaging.Paste(collage, img, image.Pt(x, y))
	//	}
	//
	//	// save the collage to file
	//	err = imaging.Save(collage, "collage.jpg")
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	fmt.Println("Collage created successfully!")

	// TODO: Figure out why the first element of this array is nil - should only be one element anyways.
	return images[1], nil
}
