package api

import (
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/go-playground/validator/v10"
)

func validatePeriod(fl validator.FieldLevel) bool {
	period := constants.Period(fl.Field().String())
	switch period {
	case constants.OVERALL, constants.SEVEN_DAYS, constants.ONE_MONTH, constants.THREE_MONTHS, constants.SIX_MONTHS, constants.TWELVE_MONTHS:
		return true
	default:
		return false
	}
}
