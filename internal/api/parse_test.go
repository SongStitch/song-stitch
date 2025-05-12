package api_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/SongStitch/song-stitch/internal/api"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
)

func TestParseQueryValues(t *testing.T) {
	defaultExpected := api.CollageRequest{
		Username:      "testuser",
		Method:        lastfm.MethodAlbum,
		TextLocation:  lastfm.LocationTopLeft,
		Period:        lastfm.PeriodSevenDays,
		Height:        0,
		Width:         0,
		Rows:          3,
		Columns:       3,
		FontSize:      12,
		DisplayAlbum:  false,
		DisplayArtist: false,
		DisplayTrack:  false,
		PlayCount:     false,
		BoldFont:      false,
		Grayscale:     false,
		Webp:          false,
	}

	tests := map[string]struct {
		query        url.Values
		expectedFunc func(c *api.CollageRequest)
		wantErr      bool
	}{
		"missing username": {
			query:   url.Values{},
			wantErr: true,
		},
		"default parameters": {
			query:        url.Values{"username": []string{"testuser"}},
			expectedFunc: func(c *api.CollageRequest) {},
		},
		"invalid method": {
			query:   url.Values{"username": []string{"test"}, "method": []string{"invalid"}},
			wantErr: true,
		},
		"valid method and text location": {
			query: url.Values{
				"username":     []string{"test"},
				"method":       []string{"track"},
				"textlocation": []string{"bottomright"},
			},
			expectedFunc: func(c *api.CollageRequest) {
				c.Username = "test"
				c.TextLocation = lastfm.LocationBottomRight
				c.Method = lastfm.MethodTrack
			},
		},
		"invalid height": {
			query:   url.Values{"username": []string{"test"}, "height": []string{"invalid"}},
			wantErr: true,
		},
		"height exceeds maximum": {
			query:   url.Values{"username": []string{"test"}, "height": []string{"3001"}},
			wantErr: true,
		},
		"valid custom dimensions": {
			query: url.Values{"username": []string{"test"}, "height": []string{"500"}, "width": []string{"800"}},
			expectedFunc: func(c *api.CollageRequest) {
				c.Username = "test"
				c.Height = 500
				c.Width = 800
			},
		},
		"custom rows and columns": {
			query: url.Values{"username": []string{"test"}, "rows": []string{"5"}, "columns": []string{"4"}},
			expectedFunc: func(c *api.CollageRequest) {
				c.Username = "test"
				c.Rows = 5
				c.Columns = 4
			},
		},
		"invalid font size": {
			query:   url.Values{"username": []string{"test"}, "fontsize": []string{"31"}},
			wantErr: true,
		},
		"boolean flags enabled": {
			query: url.Values{
				"username":  []string{"test"},
				"album":     []string{"true"},
				"artist":    []string{"1"},
				"track":     []string{"on"},
				"playcount": []string{"yes"},
				"boldfont":  []string{"true"},
				"grayscale": []string{"1"},
				"webp":      []string{"true"},
			},
			wantErr: true,
		},
		"valid boolean flags": {
			query: url.Values{
				"username":  []string{"test"},
				"album":     []string{"true"},
				"artist":    []string{"1"},
				"track":     []string{"t"},
				"playcount": []string{"1"},
				"boldfont":  []string{"true"},
				"grayscale": []string{"1"},
				"webp":      []string{"true"},
			},
			expectedFunc: func(c *api.CollageRequest) {
				c.Username = "test"
				c.DisplayAlbum = true
				c.DisplayArtist = true
				c.DisplayTrack = true
				c.PlayCount = true
				c.BoldFont = true
				c.Grayscale = true
				c.Webp = true
			},
		},
		"case insensitive parameters": {
			query: url.Values{"USERNAME": []string{"test"}, "METHOD": []string{"artist"}, "PERIOD": []string{"overall"}},
			expectedFunc: func(c *api.CollageRequest) {
				c.Username = "test"
				c.Method = lastfm.MethodArtist
				c.Period = lastfm.PeriodOverall
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := api.ParseQueryValues(tc.query)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			expected := defaultExpected
			tc.expectedFunc(&expected)

			if !reflect.DeepEqual(result, &expected) {
				t.Errorf("result mismatch:\nExpected: %+v\nGot:      %+v", expected, result)
			}
		})
	}
}
