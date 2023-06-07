package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type SpotifyAuthRequest struct {
	grant_type    string
	client_id     string
	client_secret string
}

type SpotifyAuthResponse struct {
	access_token string
	token_type   string
	expires_in   int
}

type SpotifyClient struct {
	token string
}

func NewSpotifyClient() (*SpotifyClient, error) {
	client_id := os.Getenv("SPOTIFY_CLIENT_ID")
	client_secret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if client_id == "" || client_secret == "" {
		return nil, errors.New("spotify credentials not set")
	}

	url := "https://accounts.spotify.com/api/token"
	requestBody := SpotifyAuthRequest{
		grant_type:    "client_credentials",
		client_id:     client_id,
		client_secret: client_secret,
	}

	marshalled, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalled))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("spotify authentication failed")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response SpotifyAuthResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	return &SpotifyClient{token: response.access_token}, nil
}
