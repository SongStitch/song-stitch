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
	Rows     int    `url:"rows"`
	Columns  int    `url:"columns"`
	Username string `url:"username"`
	Period   string `url:"period"`
}

func get_collage(request *CollageRequest) image.Image {
	limit := request.Rows * request.Columns
	albums := get_albums(request.Username, request.Period, limit)

	err := downloadImagesForAlbums(albums)
	if err != nil {
		log.Println(err)
	}

	collage, _ := create_collage(albums, request.Rows, request.Columns)
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
