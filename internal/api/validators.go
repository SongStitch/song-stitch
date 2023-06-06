package api

import (
	"github.com/SongStitch/song-stitch/internal/session"
	"github.com/go-playground/validator/v10"
)

func validatePeriod(fl validator.FieldLevel) bool {
	period := session.Period(fl.Field().String())
	switch period {
	case session.OVERALL, session.SEVEN_DAYS, session.ONE_MONTH, session.THREE_MONTHS, session.SIX_MONTHS, session.TWELVE_MONTHS:
		return true
	default:
		return false
	}
}
