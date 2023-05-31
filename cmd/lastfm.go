package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
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
	log.Println("Fetching last.fm data with method", method, "for username", username, "period", period, "count", count, "imageSize", imageSize)

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

type ArtistSearchResponse struct {
	Results struct {
		ArtistMatches struct {
			Artists []struct {
				Name string `json:"name"`
				Mbid string `json:"mbid"`
			} `json:"artist"`
		} `json:"artistmatches"`
	} `json:"results"`
}

func searchArtist(name string) (*ArtistSearchResponse, error) {

	endpoint := os.Getenv("LASTFM_ENDPOINT")
	key := os.Getenv("LASTFM_API_KEY")
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("method", "artist.search")
	q.Set("artist", name)
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

	var response ArtistSearchResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
