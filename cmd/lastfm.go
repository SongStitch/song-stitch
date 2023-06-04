package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/PuerkitoBio/goquery"
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
			Images    []LastFMImage `json:"image"`
			AlbumName string        `json:"title"`
		} `json:"Album"`
	} `json:"track"`
}

type TrackInfo struct {
	AlbumName string
	ImageUrl  string
}

func getTrackInfo(trackName string, artistName string, imageSize string) (*TrackInfo, error) {

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
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("track not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response GetTrackInfoResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	for _, image := range response.Track.Album.Images {
		if image.Size == imageSize {
			return &TrackInfo{response.Track.Album.AlbumName, image.Link}, nil
		}
	}
	return nil, errors.New("no image found")

}

func getImageIdForArtist(artistUrl string) (string, error) {
	url := artistUrl + "/+images"
	log.Println("Getting image for artist ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	href := (doc.Find(".image-list-item-wrapper").First().Find("a").First().AttrOr("href", ""))
	if href == "" {
		return "", errors.New("no image found")
	}
	return path.Base(href), nil

}
