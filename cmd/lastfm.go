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
			Artist struct {
				ArtistName string `json:"name"`
			}
			Playcount int    `json:"playcount"`
			AlbumName string `json:"name"`
			Image     []struct {
				Size string `json:"size"`
				Link string `json:"#text"`
			} `json:"image"`
		} `json:"album"`
	} `json:"topalbums"`
}

func getAlbums(username string, period Period, count int) []Album {
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

	var albums []Album
	albums = make([]Album, count)
	for i, album := range lastFMResponse.TopAlbums.Album {
		albums[i].Name = album.AlbumName
		albums[i].Artist = album.Artist.ArtistName
		albums[i].Playcount = album.Playcount
		for _, image := range album.Image {
			if image.Size == "extralarge" {
				albums[i].ImageUrl = image.Link
			}
		}
	}

	return albums
}
