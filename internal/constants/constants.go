package constants

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

type TextLocation string

const (
	TOP_LEFT      TextLocation = "topleft"
	TOP_CENTRE    TextLocation = "topcentre"
	TOP_RIGHT     TextLocation = "topright"
	BOTTOM_LEFT   TextLocation = "bottomleft"
	BOTTOM_CENTRE TextLocation = "bottomcentre"
	BOTTOM_RIGHT  TextLocation = "bottomright"
)

func GetTextLocationFromStr(s string) TextLocation {
	switch s {
	case "topleft":
		return TOP_LEFT
	case "topcentre":
		return TOP_CENTRE
	case "topright":
		return TOP_RIGHT
	case "bottomleft":
		return BOTTOM_LEFT
	case "bottomcentre":
		return BOTTOM_CENTRE
	case "bottomright":
		return BOTTOM_RIGHT
	default:
		return TOP_LEFT
	}
}

func (tl TextLocation) IsTop() bool {
	return tl == TOP_LEFT || tl == TOP_CENTRE || tl == TOP_RIGHT
}
