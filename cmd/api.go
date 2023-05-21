package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/url"

	"github.com/dyninc/qstring"
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
	Rows     int    `url:"rows" validate:"required,gte=1,lte=10"`
	Columns  int    `url:"columns" validate:"required,gte=1,lte=10"`
	Username string `url:"username" validate:"required"`
	Period   Period `url:"period" validate:"required,validatePeriod"`
}

func get_collage(request *CollageRequest) image.Image {
	count := request.Rows * request.Columns
	albums := get_albums(request.Username, request.Period, count)

	err := downloadImagesForAlbums(albums)
	if err != nil {
		log.Println(err)
	}

	collage, _ := create_collage(albums, request.Rows, request.Columns)
	return collage
}

func collage(w http.ResponseWriter, r *http.Request) {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request CollageRequest
	err = qstring.Unmarshal(queryParams, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	validate.RegisterValidation("validatePeriod", validatePeriod)

	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := get_collage(&request)
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
