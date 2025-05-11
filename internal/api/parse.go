package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
)

type CollageRequest struct {
	Method        lastfm.Method
	TextLocation  lastfm.TextLocation
	Username      string
	Period        lastfm.Period
	Height        uint
	Width         uint
	Rows          int
	Columns       int
	FontSize      int
	DisplayAlbum  bool
	DisplayArtist bool
	DisplayTrack  bool
	PlayCount     bool
	BoldFont      bool
	Grayscale     bool
	Webp          bool
}

var ErrInvalidValue = errors.New("invalid value")

func parseIntWithDefault(value string, d int) (int, error) {
	if value == "" {
		return d, nil
	} else {
		value, err := strconv.Atoi(value)
		if err != nil {
			return -1, err
		}
		return value, nil
	}
}

func parseIntWithDefaultAndRange(value string, d, min, max int) (int, error) {
	if value == "" {
		return d, nil
	} else {
		value, err := strconv.Atoi(value)
		if err != nil {
			return -1, err
		}
		if value < min || value > max {
			return -1, ErrInvalidValue
		}
		return value, nil
	}
}

func parseBoolWithDefault(value string, d bool) (bool, error) {
	if value == "" {
		return d, nil
	} else {
		value, err := strconv.ParseBool(value)
		if err != nil {
			return false, err
		}
		return value, nil
	}
}

func ParseRequest(r *http.Request) (*CollageRequest, error) {
	params := &CollageRequest{}
	q := r.URL.Query()

	{
		method := q.Get("method")
		if method == "" {
			params.Method = lastfm.MethodAlbum
		} else {
			method, err := lastfm.GetMethodFromStr(method)
			if err != nil {
				return nil, err
			}
			params.Method = method
		}
	}

	{
		textLocation := q.Get("textlocation")
		if textLocation == "" {
			params.TextLocation = lastfm.LocationTopLeft
		} else {
			location, err := lastfm.GetTextLocationFromStr(textLocation)
			if err != nil {
				return nil, err
			}
			params.TextLocation = location
		}
	}

	{
		username := q.Get("username")
		if username == "" {
			return nil, errors.New("username is required")
		}
		params.Username = username
	}

	{
		period := q.Get("period")
		if period == "" {
			params.Period = lastfm.PeriodSevenDays
		} else {
			period, err := lastfm.GetPeriodFromStr(period)
			if err != nil {
				return nil, err
			}
			params.Period = period
		}
	}

	{
		height := q.Get("height")
		value, err := parseIntWithDefaultAndRange(height, 0, 0, 3000)
		if err != nil {
			return nil, fmt.Errorf("invalid height: %w", err)
		}
		params.Height = uint(value)
	}

	{
		width := q.Get("width")
		value, err := parseIntWithDefaultAndRange(width, 0, 0, 3000)
		if err != nil {
			return nil, fmt.Errorf("invalid width: %w", err)
		}
		params.Width = uint(value)
	}

	{
		rows := q.Get("rows")
		value, err := parseIntWithDefault(rows, 3)
		if err != nil {
			return nil, fmt.Errorf("invalid rows: %w", err)
		}
		params.Rows = value
	}

	{
		columns := q.Get("columns")
		value, err := parseIntWithDefault(columns, 3)
		if err != nil {
			return nil, fmt.Errorf("invalid columns: %w", err)
		}
		params.Columns = value
	}

	{
		size := q.Get("fontsize")
		value, err := parseIntWithDefaultAndRange(size, 12, 8, 30)
		if err != nil {
			return nil, fmt.Errorf("invalid font size: %w", err)
		}
		params.FontSize = value
	}

	{
		album := q.Get("album")
		value, err := parseBoolWithDefault(album, false)
		if err != nil {
			return nil, fmt.Errorf("invalid album: %w", err)
		}
		params.DisplayAlbum = value
	}

	{
		artist := q.Get("artist")
		value, err := parseBoolWithDefault(artist, false)
		if err != nil {
			return nil, fmt.Errorf("invalid artist: %w", err)
		}
		params.DisplayArtist = value
	}

	{
		track := q.Get("track")
		value, err := parseBoolWithDefault(track, false)
		if err != nil {
			return nil, fmt.Errorf("invalid track: %w", err)
		}
		params.DisplayTrack = value
	}

	{
		playcount := q.Get("playcount")
		value, err := parseBoolWithDefault(playcount, false)
		if err != nil {
			return nil, fmt.Errorf("invalid playcount: %w", err)
		}
		params.PlayCount = value
	}

	{
		boldfont := q.Get("boldfont")
		value, err := parseBoolWithDefault(boldfont, false)
		if err != nil {
			return nil, fmt.Errorf("invalid boldfont: %w", err)
		}
		params.BoldFont = value
	}

	{
		grayscale := q.Get("grayscale")
		value, err := parseBoolWithDefault(grayscale, false)
		if err != nil {
			return nil, fmt.Errorf("invalid grayscale: %w", err)
		}
		params.Grayscale = value
	}

	{
		Webp := q.Get("Webp")
		value, err := parseBoolWithDefault(Webp, false)
		if err != nil {
			return nil, fmt.Errorf("invalid Webp: %w", err)
		}
		params.Webp = value
	}

	return params, nil
}
