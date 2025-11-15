package lastfm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/clients"
	"github.com/SongStitch/song-stitch/internal/config"
)

type LastfmImage struct {
	Size string `json:"size"`
	Link string `json:"#text"`
}

type LastfmUser struct {
	User       string `json:"user"`
	TotalPages string `json:"totalPages"`
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	Total      string `json:"total"`
}

func getMethodForCollageType(collageType Method) string {
	switch collageType {
	case MethodAlbum:
		return "user.gettopalbums"
	case MethodArtist:
		return "user.gettopartists"
	case MethodTrack:
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

func GetLastFmResponse(
	ctx context.Context,
	collageType Method,
	username string,
	period Period,
	count int,
	handler func(data []byte) (int, int, error),
) error {
	config := config.GetConfig()
	endpoint := config.Lastfm.Endpoint
	apiKey := config.Lastfm.APIKey

	// Image URLs stop getting returned by the API at around 500
	const maxPerPage = 500
	totalFetched := 0
	page := 1

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
		limit := min(count-totalFetched, maxPerPage)
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
				return nil, ErrUserNotFound
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
			return err
		}

		fetched, totalPages, err := handler(body)
		if err != nil {
			return err
		}
		if totalPages == page || totalPages == 0 {
			break
		}
		totalFetched = fetched
		page++
	}
	return nil // No more pages to fetch
}

type GetTrackInfoResponse struct {
	Track struct {
		Album struct {
			AlbumName string        `json:"title"`
			Images    []LastfmImage `json:"image"`
		} `json:"Album"`
	} `json:"track"`
}

func GetTrackInfo(
	trackName string,
	artistName string,
	imageSize string,
) (clients.TrackInfo, error) {
	config := config.GetConfig()
	endpoint := config.Lastfm.Endpoint
	apiKey := config.Lastfm.APIKey

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
		return clients.TrackInfo{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// ensure sensitive information is not returned in error message
		return clients.TrackInfo{}, cleanError(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return clients.TrackInfo{}, errors.New("track not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return clients.TrackInfo{}, err
	}

	var response GetTrackInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return clients.TrackInfo{}, err
	}

	for _, image := range response.Track.Album.Images {
		if image.Size == imageSize {
			return clients.TrackInfo{
				AlbumName: response.Track.Album.AlbumName,
				ImageUrl:  image.Link,
			}, nil
		}
	}
	return clients.TrackInfo{}, errors.New("no image found")
}

const maxRetries = 3

var (
	backoffSchedule = []time.Duration{
		200 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
)

func GetImageIdForArtistWithRetry(ctx context.Context, artistUrl string) (string, error) {
	var e error
	for i := range maxRetries {
		url, err := GetImageIdForArtist(ctx, artistUrl)
		if err == nil {
			return url, nil
		}

		e = err
		elem := min(len(backoffSchedule)-1, i)
		delay := backoffSchedule[elem]
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
			continue
		}
	}

	return "", fmt.Errorf("failed to get artist image after %d retries: %w", maxRetries, e)

}

func GetImageIdForArtist(ctx context.Context, artistUrl string) (string, error) {
	url := strings.TrimRight(artistUrl, "/") + "/+images"
	logger := zerolog.Ctx(ctx)

	logger.Info().Str("artistUrl", url).Msg("Getting image for artist")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:145.0) Gecko/20100101 Firefox/145.0",
	)

	if dumpReq, err := httputil.DumpRequestOut(req, false); err == nil {
		logger.Debug().Msgf("Outgoing request:\n%s", dumpReq)
	} else {
		logger.Error().Err(err).Msg("Failed to dump outgoing request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if dumpResp, err := httputil.DumpResponse(resp, true); err == nil {
		logger.Debug().Msgf("Incoming response:\n%s", dumpResp)
	} else {
		logger.Error().Err(err).Msg("Failed to dump incoming response")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	href := doc.Find(".image-list-item-wrapper").First().Find("a").First().AttrOr("href", "")
	if href == "" {
		return "", errors.New("no image found")
	}
	return path.Base(href), nil
}
