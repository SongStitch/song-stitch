package api

import (
	"context"
	"errors"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/collages"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/SongStitch/song-stitch/internal/generator"
)

type CollageRequest struct {
	Rows          int    `in:"query=rows;default=3" validate:"required,gte=1,lte=15"`
	Columns       int    `in:"query=columns;default=3" validate:"required,gte=1,lte=15"`
	Username      string `in:"query=username;required" validate:"required"`
	Period        string `in:"query=period;default=7day" validate:"required,validatePeriod"`
	DisplayArtist bool   `in:"query=artist;default=false"`
	DisplayAlbum  bool   `in:"query=album;default=false"`
	DisplayTrack  bool   `in:"query=track;default=false"`
	PlayCount     bool   `in:"query=playcount;default=false"`
	Compress      bool   `in:"query=compress;default=false"`
	Width         uint   `in:"query=width;default=0" validate:"gte=0,lte=3000"`
	Height        uint   `in:"query=height;default=0" validate:"gte=0,lte=3000"`
	Method        string `in:"query=method;default=album" validate:"required,oneof=album artist track"`
	FontSize      int    `in:"query=fontsize;default=12" validate:"gte=8,lte=30"`
}

func generateCollage(ctx context.Context, request *CollageRequest) (*image.Image, error) {
	count := request.Rows * request.Columns
	imageSize := "extralarge"
	imageDimension := 300
	if count > 100 && count <= 1000 {
		imageSize = "large"
		imageDimension = 174
	} else if count > 1000 && count <= 2000 {
		imageSize = "medium"
		imageDimension = 64
	} else if count > 2000 {
		imageSize = "small"
		imageDimension = 3
	}

	displayOptions := generator.DisplayOptions{
		ArtistName:     request.DisplayArtist,
		AlbumName:      request.DisplayAlbum,
		TrackName:      request.DisplayTrack,
		PlayCount:      request.PlayCount,
		Resize:         request.Width > 0 || request.Height > 0,
		Width:          request.Width,
		Height:         request.Height,
		Compress:       request.Compress,
		ImageDimension: imageDimension,
		FontSize:       float64(request.FontSize),
		Rows:           request.Rows,
		Columns:        request.Columns,
	}

	period := constants.GetPeriodFromStr(request.Period)
	method := constants.GetCollageTypeFromStr(request.Method)
	switch method {
	case constants.ALBUM:
		return collages.GenerateCollageForAlbum(ctx, request.Username, period, count, imageSize, displayOptions)
	case constants.ARTIST:
		return collages.GenerateCollageForArtist(ctx, request.Username, period, count, imageSize, displayOptions)
	case constants.TRACK:
		return collages.GenerateCollageForTrack(ctx, request.Username, period, count, imageSize, displayOptions)
	default:
		return nil, errors.New("invalid collage type")
	}
}

func Collage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := ctx.Value(httpin.Input).(*CollageRequest)
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Received request")

	validate := validator.New()
	validate.RegisterValidation("validatePeriod", validatePeriod)

	err := validate.Struct(request)
	if err != nil {
		logger.Warn().Err(err).Msg("Request was invalid")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info().
		Str("username", request.Username).
		Int("rows", request.Rows).
		Int("columns", request.Columns).
		Str("period", request.Period).
		Bool("artist", request.DisplayArtist).
		Bool("album", request.DisplayAlbum).
		Bool("track", request.DisplayTrack).
		Bool("playcount", request.PlayCount).
		Bool("compress", request.Compress).
		Uint("width", request.Width).
		Uint("height", request.Height).
		Str("method", request.Method).
		Int("fontsize", request.FontSize).
		Msg("Generating collage")

	image, err := generateCollage(ctx, request)
	if ctx.Err() != nil {
		logger.Warn().Err(ctx.Err()).Msg("Context cancelled")
		// 499 is the http status code for client closed request
		http.Error(w, "Context cancelled", 499)
		return
	}
	if err != nil {
		switch {
		case err == constants.ErrUserNotFound:
			logger.Warn().Err(err).Str("username", request.Username).Msg("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
		case err == constants.ErrTooManyImages:
			logger.Warn().Err(err).Str("method", request.Method).Int("rows", request.Rows).Int("columns", request.Columns).Msg("Too many images requested for the collage type")
			http.Error(w, "Requested collage size is too large for the collage type", http.StatusBadRequest)
		default:
			logger.Error().Err(err).Msg("Error occurred generating collage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, *image, nil)
	if ctx.Err() != nil {
		logger.Warn().Err(ctx.Err()).Msg("Context cancelled")
		// 499 is the http status code for client closed request
		http.Error(w, "Context cancelled", 499)
		return
	}
	if err != nil {
		logger.Error().Err(err).Msg("Error occurred encoding collage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
