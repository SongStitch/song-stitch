package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/SongStitch/song-stitch/internal/models"
	"github.com/rs/zerolog"
)

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

type Token struct {
	AccessToken string
	ExpiresIn   int
	client      *http.Client
	endpoint    string
}

func (t *Token) Refresh() error {
	client_id := os.Getenv("SPOTIFY_CLIENT_ID")
	client_secret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if client_id == "" || client_secret == "" {
		return errors.New("spotify credentials not set")
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)

	req, err := http.NewRequest(http.MethodPost, t.endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("spotify authentication failed")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response SpotifyAuthResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return err
	}
	t.AccessToken = response.AccessToken
	t.ExpiresIn = response.ExpiresIn
	return nil
}

func (t *Token) KeepAlive(log zerolog.Logger) {
	for {
		log.Info().Msg("refreshing spotify token...")
		err := t.Refresh()
		if err != nil {
			log.Error().Err(err).Msg("failed to refresh spotify token")
		}
		// Wait until 5 minutes before expiration and then refresh again
		time.Sleep(time.Duration(t.ExpiresIn-5*60) * time.Second)
	}
}

type SpotifyClient struct {
	token    *Token
	endpoint string
	client   *http.Client
}

var spotifyClient *SpotifyClient

func GetSpotifyClient() (*SpotifyClient, error) {
	if spotifyClient == nil {
		return nil, constants.ErrSpotifyClientNotInitialized
	}
	return spotifyClient, nil
}

func InitSpotifyClient(log zerolog.Logger) {
	token := &Token{
		client:   http.DefaultClient,
		endpoint: "https://accounts.spotify.com/api/token",
	}
	err := token.Refresh()
	if err != nil {
		log.Error().Err(err).Msg("failed to get spotify token, not using spotify client")
		return
	}
	go token.KeepAlive(log)
	client := &SpotifyClient{token: token, endpoint: "https://api.spotify.com/v1/search", client: http.DefaultClient}
	spotifyClient = client
}

var spotifyMarkets = [5]string{"US", "AU", "CA", "GB", "JP"}

func (c *SpotifyClient) doRequest(ctx context.Context, requestType string, queryParams map[string]string, market string) ([]byte, error) {
	logger := zerolog.Ctx(ctx)
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}
	queryAsString := ""
	for key, value := range queryParams {
		queryAsString += key + ": " + value + " "
	}
	q := u.Query()
	q.Set("q", queryAsString)
	q.Set("type", requestType)
	q.Set("market", market)
	q.Set("limit", "10")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Warn().Int("status", res.StatusCode).Msg("Spotify returned non-200 status")
		return nil, errors.New("unexpected status code: " + strconv.Itoa(res.StatusCode))

	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *SpotifyClient) doTrackRequest(ctx context.Context, trackName string, artistName string, market string) (*models.TrackInfo, error) {
	body, err := c.doRequest(ctx, "track", map[string]string{"track": trackName, "artist": artistName}, market)
	if err != nil {
		return nil, err
	}
	var response TracksResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	for _, item := range response.SearchResult.Items {
		if strings.EqualFold(item.Artists[0].Name, artistName) {
			for _, image := range item.Album.Images {
				if image.Height == 300 {
					return &models.TrackInfo{ImageUrl: image.URL, AlbumName: item.Album.Name}, nil
				}
			}
			// if no images 300x300, just return the first image
			if len(item.Album.Images) > 0 {
				return &models.TrackInfo{ImageUrl: item.Album.Images[0].URL, AlbumName: item.Album.Name}, nil
			}
		}
	}
	return nil, errors.New("track not found in market")
}

func (c *SpotifyClient) GetTrackInfo(ctx context.Context, trackName string, artistName string) (*models.TrackInfo, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Str("track", trackName).Str("artist", artistName).Msg("Fetching Spotify data")
	for _, market := range spotifyMarkets {
		track, err := c.doTrackRequest(ctx, trackName, artistName, market)
		if err != nil {
			logger.Warn().Err(err).Str("track", trackName).Str("artist", artistName).Str("market", market).Msg("Error fetching track info in market")
			continue
		}
		if track.ImageUrl != "" {
			return track, nil
		}
	}
	return nil, errors.New("track not found in any market")
}

func (c *SpotifyClient) doAlbumRequest(ctx context.Context, albumName string, artistName string, market string) (*models.AlbumInfo, error) {
	body, err := c.doRequest(ctx, "album", map[string]string{"album": albumName, "artist": artistName}, market)
	if err != nil {
		return nil, err
	}
	var response AlbumResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	for _, item := range response.SearchResult.Items {
		if strings.EqualFold(item.Artists[0].Name, artistName) {
			for _, image := range item.Images {
				if image.Height == 300 {
					return &models.AlbumInfo{ImageUrl: image.URL}, nil
				}
			}
			// if no images 300x300, just return the first image
			if len(item.Images) > 0 {
				return &models.AlbumInfo{ImageUrl: item.Images[0].URL}, nil
			}
		}
	}
	return nil, errors.New("album not found in market")
}

func (c *SpotifyClient) GetAlbumInfo(ctx context.Context, albumName string, artistName string) (*models.AlbumInfo, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Str("album", albumName).Str("artist", artistName).Msg("Fetching Spotify data")
	for _, market := range spotifyMarkets {
		album, err := c.doAlbumRequest(ctx, albumName, artistName, market)
		if err != nil {
			logger.Warn().Err(err).Str("album", albumName).Str("artist", artistName).Str("market", market).Msg("Error fetching album info in market")
			continue
		}
		if album.ImageUrl != "" {
			return album, nil
		}
	}
	return nil, errors.New("album not found in any market")
}
