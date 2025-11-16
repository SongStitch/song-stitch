package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Lastfm struct {
		Endpoint string
		APIKey   string
	}
	Fanart struct {
		APIKey string
	}
	Spotify struct {
		ClientId     string
		ClientSecret string
	}
	MaxImages struct {
		Albums  int
		Artists int
		Tracks  int
	}
	ImageSizeCutoffs struct {
		ExtraLarge int
		Large      int
		Medium     int
	}
}

var config *Config

func GetConfig() *Config {
	return config
}

func parseIntWithDefault(configField *int, name string, d int) error {
	v := os.Getenv(name)
	if v == "" {
		*configField = d
	} else {
		value, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid '%s': %w", name, err)
		}
		*configField = value
	}
	return nil
}

func Init() error {
	c := Config{}

	c.Lastfm.Endpoint = os.Getenv("LASTFM_ENDPOINT")
	if c.Lastfm.Endpoint == "" {
		return fmt.Errorf("'%s' is required", "LASTFM_ENDPOINT")
	}
	c.Lastfm.APIKey = os.Getenv("LASTFM_API_KEY")
	if c.Lastfm.APIKey == "" {
		return fmt.Errorf("'%s' is required", "LASTFM_API_KEY")
	}

	c.Fanart.APIKey = os.Getenv("FANART_API_KEY")
	if c.Fanart.APIKey == "" {
		return fmt.Errorf("'%s' is required", "FANART_API_KEY")
	}

	c.Spotify.ClientId = os.Getenv("SPOTIFY_CLIENT_ID")
	c.Spotify.ClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	if err := parseIntWithDefault(&c.MaxImages.Albums, "MAX_ALBUM_IMAGES", 400); err != nil {
		return err
	}
	if err := parseIntWithDefault(&c.MaxImages.Artists, "MAX_ARTIST_IMAGES", 400); err != nil {
		return err
	}
	if err := parseIntWithDefault(&c.MaxImages.Tracks, "MAX_TRACK_IMAGES", 400); err != nil {
		return err
	}
	if err := parseIntWithDefault(&c.ImageSizeCutoffs.ExtraLarge, "IMAGE_SIZE_CUTOFF_EXTRA_LARGE", 100); err != nil {
		return err
	}
	if err := parseIntWithDefault(&c.ImageSizeCutoffs.Large, "IMAGE_SIZE_CUTOFF_LARGE", 1000); err != nil {
		return err
	}
	if err := parseIntWithDefault(&c.ImageSizeCutoffs.Medium, "IMAGE_SIZE_CUTOFF_MEDIUM", 2000); err != nil {
		return err
	}

	config = &c

	return nil
}
