package lastfm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/clients"
	"github.com/SongStitch/song-stitch/internal/config"
	"github.com/SongStitch/song-stitch/internal/constants"
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

func getMethodForCollageType(collageType constants.CollageType) string {
	switch collageType {
	case constants.ALBUM:
		return "user.gettopalbums"
	case constants.ARTIST:
		return "user.gettopartists"
	case constants.TRACK:
		return "user.gettoptracks"
	default:
		return ""
	}
}

type CleanError struct {
	errStr string
}

func (e CleanError) Error() string {
	return e.errStr
}

// strip sensitive information from error message
func cleanError(err error) error {
	errStr := err.Error()
	pattern := `(&|\?)api_key=[^&]+(&|\b)`
	regex := regexp.MustCompile(pattern)
	modifiedString := regex.ReplaceAllString(errStr, "$1")
	return CleanError{errStr: modifiedString}
}

func GetLastFmResponse[T LastFMResponse](
	ctx context.Context,
	collageType constants.CollageType,
	username string,
	period constants.Period,
	count int,
) (*T, error) {
	config := config.GetConfig()
	endpoint := config.LastFM.Endpoint
	apiKey := config.LastFM.APIKey

	// Image URLs stop getting returned by the API at around 500
	const maxPerPage = 500
	totalFetched := 0
	page := 1

	var result T
	initialised := false

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Fetching LastFM data")
	method := getMethodForCollageType(collageType)
	for count > totalFetched {
		logger.Info().
			Int("page", page).
			Int("totalFetched", totalFetched).
			Int("count", count).
			Msg("Fetching page")
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
		q.Set("api_key", apiKey)
		q.Set("format", "json")
		u.RawQuery = q.Encode()

		body, err := func() ([]byte, error) {
			req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
			if err != nil {
				return nil, err
			}

			start := time.Now()
			res, err := http.DefaultClient.Do(req)
			if res != nil {
				defer res.Body.Close()
			}
			logger.Info().
				Dur("duration", time.Since(start)).
				Str("method", method).
				Msg("Last.fm request completed")
			if err != nil {
				// ensure sensitive information is not returned in error message
				return nil, cleanError(err)
			}

			if res.StatusCode == http.StatusNotFound {
				return nil, constants.ErrUserNotFound
			}

			if res.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			return body, nil
		}()
		if err != nil {
			return nil, err
		}

		var response T
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if !initialised {
			result = response
			initialised = true
		} else {
			err = result.Append(response)
			if err != nil {
				return nil, err
			}
		}
		totalFetched = result.TotalFetched()
		totalPages := result.TotalPages()
		if totalPages == page || totalPages == 0 {
			break
		}
		page++
	}
	return &result, nil // No more pages to fetch
}

type GetTrackInfoResponse struct {
	Track struct {
		Album struct {
			AlbumName string        `json:"title"`
			Images    []LastFMImage `json:"image"`
		} `json:"Album"`
	} `json:"track"`
}

func GetTrackInfo(
	trackName string,
	artistName string,
	imageSize string,
) (*clients.TrackInfo, error) {
	config := config.GetConfig()
	endpoint := config.LastFM.Endpoint
	apiKey := config.LastFM.APIKey

	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("track", trackName)
	q.Set("artist", artistName)
	q.Set("method", "track.getInfo")
	q.Set("api_key", apiKey)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// ensure sensitive information is not returned in error message
		return nil, cleanError(err)
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
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, image := range response.Track.Album.Images {
		if image.Size == imageSize {
			return &clients.TrackInfo{
				AlbumName: response.Track.Album.AlbumName,
				ImageUrl:  image.Link,
			}, nil
		}
	}
	return nil, errors.New("no image found")
}

func GetImageIdForArtist(ctx context.Context, artistUrl string) (string, error) {
	url := artistUrl + "/+images"
	zerolog.Ctx(ctx).Info().Str("artistUrl", url).Msg("Getting image for artist")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
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
