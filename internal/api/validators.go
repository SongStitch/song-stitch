package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
)

func validatePeriod(fl validator.FieldLevel) bool {
	period := lastfm.Period(fl.Field().String())
	switch period {
	case lastfm.OVERALL,
		lastfm.SEVEN_DAYS,
		lastfm.ONE_MONTH,
		lastfm.THREE_MONTHS,
		lastfm.SIX_MONTHS,
		lastfm.TWELVE_MONTHS:
		return true
	default:
		return false
	}
}

func validateTextLocation(fl validator.FieldLevel) bool {
	textLocation := lastfm.TextLocation(fl.Field().String())
	switch textLocation {
	case lastfm.TOP_LEFT,
		lastfm.TOP_CENTRE,
		lastfm.TOP_RIGHT,
		lastfm.BOTTOM_LEFT,
		lastfm.BOTTOM_CENTRE,
		lastfm.BOTTOM_RIGHT:
		return true
	default:
		return false
	}
}
