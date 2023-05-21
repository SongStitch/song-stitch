package main

import (
	"fmt"
	"net/http"
)

func collage(w http.ResponseWriter, r *http.Request) {

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
