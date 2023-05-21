package main

import (
	"bytes"
	"encoding/json"
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

	"github.com/dyninc/qstring"
)

const (
	jpgFileType = ".jpg"
)

type CollageRequest struct {
	Width    int    `url:"width"`
	Height   int    `url:"height"`
	Username string `url:"username"`
	Period   string `url:"period"`
}

type CollageResponse struct {
	Images []string `json:"images"`
}

type LastFMResponse struct {
	TopAlbums struct {
		Album []struct {
			Image []struct {
				Size string `json:"size"`
				Link string `json:"#text"`
			} `json:"image"`
		} `json:"album"`
	} `json:"topalbums"`
}

func get_collage(request *CollageRequest) CollageResponse {

	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")

	limit := request.Width * request.Height

	url := fmt.Sprintf("%s?method=user.gettopalbums&user=%s&period=%s&limit=%d&api_key=%s&format=json", endpoint, request.Username, request.Period, limit, key)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var lastFMResponse LastFMResponse
	err = json.Unmarshal([]byte(body), &lastFMResponse)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var imageArray []string
	imageArray = make([]string, limit)
	for i, album := range lastFMResponse.TopAlbums.Album {
		for _, image := range album.Image {
			if image.Size == "small" {
				imageArray[i] = image.Link
			}
		}
	}

	return CollageResponse{
		Images: imageArray,
	}

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

func collage(w http.ResponseWriter, r *http.Request) {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request CollageRequest
	err = qstring.Unmarshal(queryParams, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := get_collage(&request)
	imageURLs := response.Images
	// use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(len(imageURLs))

	images := make([]image.Image, len(imageURLs))

	var mux sync.Mutex

	for _, url := range imageURLs {
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

	if len(images) == 0 {
		fmt.Println("No images downloaded successfully.")
		return
	}

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
	//		y := (i / gridSize) * height
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

	responseJson, err := json.Marshal(response)
	fmt.Fprintf(w, string(responseJson))
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API running")
}
