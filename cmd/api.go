package main

import (
	"fmt"
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

type CollageRequest struct {
	Rows             int    `in:"query=rows;default=3" validate:"required,gte=1,lte=10"`
	Columns          int    `in:"query=columns;default=3" validate:"required,gte=1,lte=10"`
	Username         string `in:"query=username;required" validate:"required"`
	Period           string `in:"query=period;default=7day" validate:"required,validatePeriod"`
	DisplayArtist    bool   `in:"query=artist;default=false"`
	DisplayAlbum     bool   `in:"query=album;default=false"`
	DisplayPlaycount bool   `in:"query=playcount;default=false"`
}

func getCollage(request *CollageRequest) image.Image {
	count := request.Rows * request.Columns

	period := getPeriodFromStr(request.Period)
	albums := getAlbums(request.Username, period, count)

	err := downloadImagesForAlbums(albums)
	if err != nil {
		log.Println(err)
	}

	collage, _ := createCollage(albums, request.Rows, request.Columns, request.DisplayAlbum, request.DisplayArtist, request.DisplayPlaycount)
	return collage
}

func collage(w http.ResponseWriter, r *http.Request) {
	request := r.Context().Value(httpin.Input).(*CollageRequest)

	validate := validator.New()
	validate.RegisterValidation("validatePeriod", validatePeriod)

	err := validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := getCollage(request)
	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API running")
}
