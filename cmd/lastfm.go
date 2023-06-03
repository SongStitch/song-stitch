package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/anaskhan96/soup"
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

type LastFMResponse interface {
	Append(l LastFMResponse)
	GetTotalPages() int
	GetTotalFetched() int
}

var ErrUserNotFound = errors.New("user not found")

func getMethodForCollageType(collageType CollageType) string {
	switch collageType {
	case ALBUM:
		return "user.gettopalbums"
	case ARTIST:
		return "user.gettopartists"
	case TRACK:
		return "user.gettoptracks"
	default:
		return ""
	}
}

func getLastFmResponse[T LastFMResponse](collageType CollageType, username string, period Period, count int, imageSize string) (*T, error) {
	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")

	// Image URLs stop getting returned by the API at around 500
	const maxPerPage = 500
	var totalFetched = 0
	var page = 1

	var result T
	initialised := false

	method := getMethodForCollageType(collageType)
	for count > totalFetched {
		// Determine the limit for this request
		limit := count - totalFetched
		if limit > maxPerPage {
			limit = maxPerPage
		}
		u, err := url.Parse(endpoint)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("user", username)
		q.Set("method", method)
		q.Set("period", string(period))
		q.Set("limit", strconv.Itoa(limit))
		q.Set("page", strconv.Itoa(page))
		q.Set("api_key", key)
		q.Set("format", "json")
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
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

		var response T
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			return nil, err
		}
		if !initialised {
			result = response
			initialised = true
		} else {
			result.Append(response)
		}
		totalFetched = result.GetTotalFetched()
		totalPages := result.GetTotalPages()
		if totalPages == page {
			break
		}
		page++
	}
	return &result, nil // No more pages to fetch
}

type GetTrackInfoResponse struct {
	Track struct {
		Album struct {
			Images []LastFMImage `json:"image"`
		} `json:"Album"`
	} `json:"track"`
}

func getImageUrlForTrack(trackName string, artistName string, imageSize string) (string, error) {

	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("track", trackName)
	q.Set("artist", artistName)
	q.Set("method", "track.getInfo")
	q.Set("api_key", key)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return "", ErrUserNotFound
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response GetTrackInfoResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", err
	}

	for _, image := range response.Track.Album.Images {
		if image.Size == imageSize {
			return image.Link, nil
		}
	}
	log.Println("No image found for track ", trackName, " and artist ", artistName, " and size ", imageSize)
	return "", nil

}

func getImageUrlForArtist(artistUrl string) (string, error) {
	url := artistUrl + "/+images"
	log.Println("Getting image for artist ", url)
	resp, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc := soup.HTMLParse(resp)
	elements := doc.FindAll("class", "image-list-item-wrapper")
	if len(elements) == 0 {
		log.Fatal("No elements with class image-list-item-wrapper found")
	}

	links := elements[0].FindAll("a")
	if len(links) == 0 {
		log.Fatal("No links found in the first element with class image-list-item-wrapper")
	}

	fmt.Println(links[0].Attrs()["href"])
	return links[0].Attrs()["href"], nil

}
