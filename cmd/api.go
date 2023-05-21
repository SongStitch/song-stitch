package main

import (
	"fmt"
	"github.com/dyninc/qstring"
	"net/http"
	"net/url"
)

type CollageRequest struct {
	Width    int    `url:"width"`
	Height   int    `url:"height"`
	Username string `url:"username"`
	Period   string `url:"period"`
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

	fmt.Fprintf(w, "Width: %d\n", request.Width)
	fmt.Fprintf(w, "Height: %d\n", request.Height)
	fmt.Fprintf(w, "Username: %s\n", request.Username)
	fmt.Fprintf(w, "Period: %s\n", request.Period)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
