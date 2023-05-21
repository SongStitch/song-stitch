package main

import (
	"encoding/json"
	"fmt"
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

type CollageResponse struct {
	Images []string `json:"images"`
}

func get_collage(request *CollageRequest) CollageResponse {
	return CollageResponse{
		Images: []string{"https://i.imgur.com/3jO3l4l.jpg", "https://i.imgur.com/3jO3l4l.jpg"},
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
	responseJson, err := json.Marshal(response)
	fmt.Fprintf(w, string(responseJson))
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Api running")
}
