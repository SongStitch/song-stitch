package lastfm

import "errors"

var ErrTooManyImages = errors.New("too many images requested")
var ErrUserNotFound = errors.New("user not found")
var ErrNoImageFound = errors.New("no image found")
var ErrInvalidMethod = errors.New("invalid method")
var ErrInvalidLocation = errors.New("invalid text location")
var ErrInvalidPeriod = errors.New("invalid period")

type Period string

const (
	PeriodOverall      Period = "overall"
	PeriodSevenDays    Period = "7day"
	PeriodOneMonth     Period = "1month"
	PeriodThreeMonths  Period = "3month"
	PeriodSixMonths    Period = "6month"
	PeriodTwelveMonths Period = "12month"
)

func GetPeriodFromStr(s string) (Period, error) {
	switch s {
	case "overall":
		return PeriodOverall, nil
	case "7day":
		return PeriodSevenDays, nil
	case "1month":
		return PeriodOneMonth, nil
	case "3month":
		return PeriodThreeMonths, nil
	case "6month":
		return PeriodSixMonths, nil
	case "12month":
		return PeriodTwelveMonths, nil
	default:
		return PeriodOverall, ErrInvalidPeriod
	}
}

type Method string

const (
	MethodAlbum  Method = "album"
	MethodArtist Method = "artist"
	MethodTrack  Method = "track"
)

func GetMethodFromStr(s string) (Method, error) {
	switch s {
	case "album":
		return MethodAlbum, nil
	case "artist":
		return MethodArtist, nil
	case "track":
		return MethodTrack, nil
	default:
		return MethodAlbum, ErrInvalidMethod
	}
}

type TextLocation string

const (
	LocationTopLeft      TextLocation = "topleft"
	LocationTopCentre    TextLocation = "topcentre"
	LocationTopRight     TextLocation = "topright"
	LocationBottomLeft   TextLocation = "bottomleft"
	LocationBottomCentre TextLocation = "bottomcentre"
	LocationBottomRight  TextLocation = "bottomright"
)

func GetTextLocationFromStr(s string) (TextLocation, error) {
	switch s {
	case "topleft":
		return LocationTopLeft, nil
	case "topcentre":
		return LocationTopCentre, nil
	case "topright":
		return LocationTopRight, nil
	case "bottomleft":
		return LocationBottomLeft, nil
	case "bottomcentre":
		return LocationBottomCentre, nil
	case "bottomright":
		return LocationBottomRight, nil
	default:
		return LocationTopLeft, ErrInvalidLocation
	}
}

func (tl TextLocation) IsTop() bool {
	return tl == LocationTopLeft || tl == LocationTopCentre || tl == LocationTopRight
}
