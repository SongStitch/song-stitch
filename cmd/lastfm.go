package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type LastFMImage struct {
	Size string `json:"size"`
	Link string `json:"#text"`
}

type LastFMUser struct {
	User       string `json:"user"`
	TotalPages string `json:"totalPages"`
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	Total      string `json:"total"`
}

type LastFMAlbum struct {
	Artist struct {
		URL        string `json:"url"`
		ArtistName string `json:"name"`
		Mbid       string `json:"mbid"`
	} `json:"artist"`
	Image     []LastFMImage `json:"image"`
	Mbid      string        `json:"mbid"`
	URL       string        `json:"url"`
	Playcount string        `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	AlbumName string `json:"name"`
}

type LastFMArtist struct {
	Image     []LastFMImage `json:"image"`
	Mbid      string        `json:"mbid"`
	URL       string        `json:"url"`
	Playcount string        `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Name string `json:"name"`
}

type LastFMResponse struct {
	TopAlbums struct {
		Album []LastFMAlbum `json:"album"`
		Attr  LastFMUser    `json:"@attr"`
	} `json:"topalbums"`

	TopArtists struct {
		Artists []LastFMArtist `json:"artist"`
		Attr    LastFMUser     `json:"@attr"`
	} `json:"topartists"`
}

var ErrUserNotFound = errors.New("user not found")

func getMethodForCollageType(collageType CollageType) string {
	switch collageType {
	case ALBUM:
		return "gettopalbums"
	case ARTIST:
		return "gettopartists"
	default:
		return ""
	}
}

func getAlbums(collageType CollageType, username string, period Period, count int, imageSize string) ([]*Album, error) {
	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")

	// Image URLs stop getting returned by the API at around 500
	const maxPerPage = 500
	var totalFetched = 0
	var page = 1
	var albums []*Album

	for count > totalFetched {
		// Determine the limit for this request
		limit := count - totalFetched
		if limit > maxPerPage {
			limit = maxPerPage
		}

		method := getMethodForCollageType(collageType)
		url := fmt.Sprintf("%s?method=user.%s&user=%s&period=%s&limit=%d&page=%d&api_key=%s&format=json", endpoint, method, username, period, limit, page, key)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return nil, ErrUserNotFound
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var lastFMResponse LastFMResponse
		err = json.Unmarshal([]byte(body), &lastFMResponse)
		if err != nil {
			return nil, err
		}

		totalPages, err := strconv.Atoi(lastFMResponse.TopAlbums.Attr.TotalPages)
		if err != nil {
			return nil, err
		}

		// No more pages to fetch
		if page > totalPages {
			break
		}

		for _, album := range lastFMResponse.TopAlbums.Album {
			newAlbum := &Album{
				Name:      album.AlbumName,
				Artist:    album.Artist.ArtistName,
				Playcount: album.Playcount,
			}

			for _, image := range album.Image {
				if image.Size == imageSize {
					newAlbum.ImageUrl = image.Link
				}
			}

			albums = append(albums, newAlbum)

			totalFetched += 1
			if totalFetched >= count {
				return albums, nil
			}
		}

		// Move to next page
		page++
	}

	return albums, nil
}
