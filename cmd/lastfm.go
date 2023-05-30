package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
		return "gettopalbums"
	case ARTIST:
		return "gettopartists"
	case TRACK:
		return "gettoptracks"
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
