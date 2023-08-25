package api

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/collages"
	"github.com/SongStitch/song-stitch/internal/config"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/SongStitch/song-stitch/internal/generator"
)

type CollageRequest struct {
	Method        string `in:"query=method;default=album" validate:"required,oneof=album artist track"`
	TextLocation  string `in:"query=textlocation;default=topleft" validate:"validateTextLocation"`
	Username      string `in:"query=username;required" validate:"required"`
	Period        string `in:"query=period;default=7day" validate:"required,validatePeriod"`
	Height        uint   `in:"query=height;default=0" validate:"gte=0,lte=3000"`
	Width         uint   `in:"query=width;default=0" validate:"gte=0,lte=3000"`
	Rows          int    `in:"query=rows;default=3" validate:"required"`
	FontSize      int    `in:"query=fontsize;default=12" validate:"gte=8,lte=30"`
	Columns       int    `in:"query=columns;default=3" validate:"required"`
	DisplayAlbum  bool   `in:"query=album;default=false"`
	DisplayTrack  bool   `in:"query=track;default=false"`
	PlayCount     bool   `in:"query=playcount;default=false"`
	DisplayArtist bool   `in:"query=artist;default=false"`
	BoldFont      bool   `in:"query=boldfont;default=false"`
	Grayscale     bool   `in:"query=grayscale;default=false"`
	Webp          bool   `in:"query=webp;default=false"`
}

func generateCollage(ctx context.Context, request *CollageRequest) (*image.Image, *bytes.Buffer, error) {
	config := config.GetConfig()

	count := request.Rows * request.Columns

	var imageSize string
	var imageDimension int
	if count > config.ImageSizeCutoffs.Medium {
		imageSize = "small"
		imageDimension = 3
	} else if count > config.ImageSizeCutoffs.Large {
		imageSize = "medium"
		imageDimension = 64
	} else if count > config.ImageSizeCutoffs.ExtraLarge {
		imageSize = "large"
		imageDimension = 174
	} else {
		imageSize = "extralarge"
		imageDimension = 300
	}

	displayOptions := generator.DisplayOptions{
		ArtistName:     request.DisplayArtist,
		AlbumName:      request.DisplayAlbum,
		TrackName:      request.DisplayTrack,
		PlayCount:      request.PlayCount,
		Resize:         request.Width > 0 || request.Height > 0,
		Width:          request.Width,
		Height:         request.Height,
		ImageDimension: imageDimension,
		FontSize:       float64(request.FontSize),
		BoldFont:       request.BoldFont,
		Grayscale:      request.Grayscale,
		Webp:           request.Webp,
		Rows:           request.Rows,
		Columns:        request.Columns,
		TextLocation:   constants.GetTextLocationFromStr(request.TextLocation),
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
		return nil, nil, errors.New("invalid collage type")
	}
}

func Collage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := ctx.Value(httpin.Input).(*CollageRequest)
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Received request")

	validate := validator.New()
	validate.RegisterValidation("validatePeriod", validatePeriod)
	validate.RegisterValidation("validateTextLocation", validateTextLocation)

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
		Uint("width", request.Width).
		Uint("height", request.Height).
		Str("method", request.Method).
		Int("fontsize", request.FontSize).
		Bool("boldfont", request.BoldFont).
		Bool("grayscale", request.Grayscale).
		Bool("webp", request.Webp).
		Msg("Generating collage")

	image, buffer, err := generateCollage(ctx, request)
	if ctx.Err() != nil {
		logger.Warn().Err(ctx.Err()).Msg("Context cancelled")
		// 499 is the http status code for client closed request
		http.Error(w, "Context cancelled", 499)
		return
	}
	if err != nil {
		switch err {
		case constants.ErrUserNotFound:
			logger.Warn().Err(err).Str("username", request.Username).Msg("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
		case constants.ErrTooManyImages:
			logger.Warn().Err(err).Str("method", request.Method).Int("rows", request.Rows).Int("columns", request.Columns).Msg("Too many images requested for the collage type")
			http.Error(w, "Requested collage size is too large for the collage type", http.StatusBadRequest)
		default:
			logger.Error().Err(err).Msg("Error occurred generating collage")
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
		}
		return
	}

	if ctx.Err() != nil {
		logger.Warn().Err(ctx.Err()).Msg("Context cancelled")
		// 499 is the http status code for client closed request
		http.Error(w, "Context cancelled", 499)
		return
	}

	if request.Webp && request.Grayscale {
		request.Webp = false
	}

	if request.Webp {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(buffer.Bytes())
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, *image, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Error occurred encoding collage")
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}
	}
}
