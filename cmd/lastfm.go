package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
			Playcount string `json:"playcount"`
			AlbumName string `json:"name"`
			Image     []struct {
				Size string `json:"size"`
				Link string `json:"#text"`
			} `json:"image"`
		} `json:"album"`
	} `json:"topalbums"`
}

func getAlbums(username string, period Period, count int, imageSize string) ([]Album, error) {
	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")
	url := fmt.Sprintf("%s?method=user.gettopalbums&user=%s&period=%s&limit=%d&api_key=%s&format=json", endpoint, username, period, count, key)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var lastFMResponse LastFMResponse
	err = json.Unmarshal([]byte(body), &lastFMResponse)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// No albums to return
	if len(lastFMResponse.TopAlbums.Album) == 0 {
		log.Println("no albums!")
		return nil, errors.New("no albums found! is the username correct?")
	}

	albums := make([]Album, count)
	for i, album := range lastFMResponse.TopAlbums.Album {
		albums[i].Name = album.AlbumName
		albums[i].Artist = album.Artist.ArtistName
		albums[i].Playcount = album.Playcount
		for _, image := range album.Image {
			if image.Size == imageSize {
				albums[i].ImageUrl = image.Link
			}
		}
	}

	return albums, nil
}
