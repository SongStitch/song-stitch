package session

type Period string

const (
	OVERALL       Period = "overall"
	SEVEN_DAYS    Period = "7day"
	ONE_MONTH     Period = "1month"
	THREE_MONTHS  Period = "3month"
	SIX_MONTHS    Period = "6month"
	TWELVE_MONTHS Period = "12month"
)

func GetPeriodFromStr(s string) Period {
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

type CollageType string

const (
	ALBUM  CollageType = "album"
	ARTIST CollageType = "artist"
	TRACK  CollageType = "track"
)

func GetCollageTypeFromStr(s string) CollageType {
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
