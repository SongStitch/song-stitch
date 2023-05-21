package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func run() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}
