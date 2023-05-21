package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

func get_albums(username string, period string, count int) []string {

	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")
	url := fmt.Sprintf("%s?method=user.gettopalbums&user=%s&period=%s&limit=%d&api_key=%s&format=json", endpoint, username, period, count, key)

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
	imageArray = make([]string, count)
	for i, album := range lastFMResponse.TopAlbums.Album {
		for _, image := range album.Image {
			if image.Size == "extralarge" {
				imageArray[i] = image.Link
			}
		}
	}

	return imageArray
}
