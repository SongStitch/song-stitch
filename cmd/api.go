package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/url"

	"github.com/dyninc/qstring"
)

type CollageRequest struct {
	Width    int    `url:"width"`
	Height   int    `url:"height"`
	Username string `url:"username"`
	Period   string `url:"period"`
}

func get_collage(request *CollageRequest) image.Image {
	limit := request.Width * request.Height
	imageUrls := get_albums(request.Username, request.Period, limit)

	images, err := downloadImages(imageUrls)
	if err != nil {
		log.Println(err)
	}

	collage, _ := create_collage(images, request.Width, request.Height)
	return collage
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
	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API running")
}
