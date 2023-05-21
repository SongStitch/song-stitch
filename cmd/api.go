package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

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
	fmt.Fprintf(w, "API running")
}
