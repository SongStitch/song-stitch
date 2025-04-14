package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Lastfm struct {
		Endpoint string `env:"LASTFM_ENDPOINT,required"`
		APIKey   string `env:"LASTFM_API_KEY,required"`
	}
	Spotify struct {
		ClientId     string `env:"SPOTIFY_CLIENT_ID"`
		ClientSecret string `env:"SPOTIFY_CLIENT_SECRET"`
	}
	MaxImages struct {
		Albums  int `env:"MAX_ALBUM_IMAGES,default=400"`
		Artists int `env:"MAX_ARTIST_IMAGES,default=400"`
		Tracks  int `env:"MAX_TRACK_IMAGES,default=100"`
	}
	ImageSizeCutoffs struct {
		ExtraLarge int `env:"IMAGE_SIZE_CUTOFF_EXTRA_LARGE,default=100"`
		Large      int `env:"IMAGE_SIZE_CUTOFF_LARGE,default=1000"`
		Medium     int `env:"IMAGE_SIZE_CUTOFF_MEDIUM,default=2000"`
	}
}

var config *Config

func GetConfig() *Config {
	return config
}

func InitConfig() error {
	ctx := context.Background()

	config = &Config{}

	if err := envconfig.Process(ctx, config); err != nil {
		return err
	}
	return nil
}
