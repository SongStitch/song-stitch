package main

import (
	"errors"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-playground/validator/v10"
)

type Period string

const (
	OVERALL       Period = "overall"
	SEVEN_DAYS    Period = "7day"
	ONE_MONTH     Period = "1month"
	THREE_MONTHS  Period = "3month"
	SIX_MONTHS    Period = "6month"
	TWELVE_MONTHS Period = "12month"
)

type CollageType string

const (
	ALBUM  CollageType = "album"
	ARTIST CollageType = "artist"
	TRACK  CollageType = "track"
)

// Seems stupid but was generated by copilot so ¯\_(ツ)_/¯
func getPeriodFromStr(s string) Period {
	switch s {
	case "overall":
		return OVERALL
	case "7day":
		return SEVEN_DAYS
	case "1month":
		return ONE_MONTH
	case "3month":
		return THREE_MONTHS
	case "6month":
		return SIX_MONTHS
	case "12month":
		return TWELVE_MONTHS
	default:
		return OVERALL
	}
}

func validatePeriod(fl validator.FieldLevel) bool {
	period := Period(fl.Field().String())
	switch period {
	case OVERALL, SEVEN_DAYS, ONE_MONTH, THREE_MONTHS, SIX_MONTHS, TWELVE_MONTHS:
		return true
	default:
		return false
	}
}

func getCollageTypeFromStr(s string) CollageType {
	switch s {
	case "album":
		return ALBUM
	case "artist":
		return ARTIST
	case "track":
		return TRACK
	default:
		return ALBUM
	}
}

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
}

func generateCollageForAlbum(username string, period Period, count int, imageSize string, displayOptions DisplayOptions) (image.Image, error) {
	albums, err := getAlbums(username, period, count, imageSize)
	if err != nil {
		return nil, err
	}

	downloadImages(albums)

	return createCollage(albums, displayOptions)
}

func generateCollageForArtist(username string, period Period, count int, imageSize string, displayOptions DisplayOptions) (image.Image, error) {
	displayOptions.AlbumName = false
	artists, err := getArtists(username, period, count, imageSize)
	if err != nil {
		return nil, err
	}

	downloadImages(artists)

	return createCollage(artists, displayOptions)
}

func generateCollageForTrack(username string, period Period, count int, imageSize string, displayOptions DisplayOptions) (image.Image, error) {
	displayOptions.AlbumName = false
	tracks, err := getTracks(username, period, count, imageSize)
	if err != nil {
		return nil, err
	}

	downloadImages(tracks)

	return createCollage(tracks, displayOptions)
}
func generateCollage(request *CollageRequest) (image.Image, error) {
	count := request.Rows * request.Columns
	imageSize := "extralarge"
	imageDimension := 300
	var fontSize float64 = 12
	if count > 100 && count <= 1000 {
		imageSize = "large"
		imageDimension = 174
		fontSize = 8
	} else if count > 1000 && count <= 2000 {
		imageSize = "medium"
		imageDimension = 64
		fontSize = 6
	} else if count > 2000 {
		imageSize = "small"
		imageDimension = 34
		fontSize = 2
	}

	displayOptions := DisplayOptions{
		ArtistName:     request.DisplayArtist,
		AlbumName:      request.DisplayAlbum,
		TrackName:      request.DisplayTrack,
		PlayCount:      request.PlayCount,
		Resize:         request.Width > 0 || request.Height > 0,
		Width:          request.Width,
		Height:         request.Height,
		Compress:       request.Compress,
		ImageDimension: imageDimension,
		FontSize:       fontSize,
		Rows:           request.Rows,
		Columns:        request.Columns,
	}

	period := getPeriodFromStr(request.Period)
	method := getCollageTypeFromStr(request.Method)
	switch method {
	case ALBUM:
		return generateCollageForAlbum(request.Username, period, count, imageSize, displayOptions)
	case ARTIST:
		return generateCollageForArtist(request.Username, period, count, imageSize, displayOptions)
	case TRACK:
		return generateCollageForTrack(request.Username, period, count, imageSize, displayOptions)
	default:
		return nil, errors.New("invalid collage type")
	}
}

func collage(w http.ResponseWriter, r *http.Request) {
	request := r.Context().Value(httpin.Input).(*CollageRequest)

	validate := validator.New()
	validate.RegisterValidation("validatePeriod", validatePeriod)

	err := validate.Struct(request)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := generateCollage(request)
	if err != nil {
		log.Println(err)
		switch {
		case err == ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, response, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
