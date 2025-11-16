package api

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/collages"
	"github.com/SongStitch/song-stitch/internal/config"
)

func generateCollage(
	ctx context.Context,
	request *CollageRequest,
) (image.Image, *bytes.Buffer, error) {
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

	displayOptions := collages.DisplayOptions{
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
		TextLocation:   request.TextLocation,
	}

	jobChan := make(chan collages.CollageElement, 100)
  logger := zerolog.Ctx(ctx)
	go func() {
    var err error
		switch request.Method {
		case lastfm.MethodAlbum:
			err = collages.GetElementsForAlbum(
				ctx,
				request.Username,
				request.Period,
				count,
				imageSize,
				displayOptions,
        jobChan,
			)
		case lastfm.MethodArtist:
			err = collages.GetElementsForArtist(
				ctx,
				request.Username,
				request.Period,
				count,
				imageSize,
				displayOptions,
        jobChan,
			)
		case lastfm.MethodTrack:
			err = collages.GetElementsForTrack(
				ctx,
				request.Username,
				request.Period,
				count,
				imageSize,
				displayOptions,
				jobChan,
			)
		}
    if err != nil {
      logger.Error().Err(err).Msg("failed to fetch image data")
    }
		close(jobChan)
	}()
	return collages.CreateCollage(ctx, displayOptions, jobChan)
}

func Collage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Received request")

	request, err := ParseQueryValues(r.URL.Query())
	if err != nil {
		logger.Warn().Err(err).Msg("Request was invalid")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info().
		Str("username", request.Username).
		Int("rows", request.Rows).
		Int("columns", request.Columns).
		Str("period", string(request.Period)).
		Bool("artist", request.DisplayArtist).
		Bool("album", request.DisplayAlbum).
		Bool("track", request.DisplayTrack).
		Bool("playcount", request.PlayCount).
		Uint("width", request.Width).
		Uint("height", request.Height).
		Str("method", string(request.Method)).
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
		case lastfm.ErrUserNotFound:
			logger.Warn().Err(err).Str("username", request.Username).Msg("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
		case lastfm.ErrTooManyImages:
			logger.Warn().
				Err(err).
				Str("method", string(request.Method)).
				Int("rows", request.Rows).
				Int("columns", request.Columns).
				Msg("Too many images requested for the collage type")
			http.Error(
				w,
				"Requested collage size is too large for the collage type",
				http.StatusBadRequest,
			)
		default:
			logger.Error().Err(err).Msg("Error occurred generating collage")
			http.Error(
				w,
				"An error occurred processing your request",
				http.StatusInternalServerError,
			)
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
		err = jpeg.Encode(w, image, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Error occurred encoding collage")
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}
	}
}
