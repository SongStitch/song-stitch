package spotify

type TracksResponse struct {
	SearchResult SearchResult[TrackItem] `json:"tracks"`
}

type AlbumResponse struct {
	SearchResult SearchResult[AlbumItem] `json:"albums"`
}

type SearchResult[T any] struct {
	Href     string `json:"href"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Items    []T    `json:"items"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
}

type TrackItem struct {
	Album            Album       `json:"album"`
	Name             string      `json:"name"`
	URI              string      `json:"uri"`
	Type             string      `json:"type"`
	ExternalIds      ExternalID  `json:"external_ids"`
	ExternalUrls     ExternalURL `json:"external_urls"`
	Href             string      `json:"href"`
	ID               string      `json:"id"`
	PreviewURL       string      `json:"preview_url"`
	Artists          []Artist    `json:"artists"`
	AvailableMarkets []string    `json:"available_markets"`
	DurationMs       int         `json:"duration_ms"`
	Popularity       int         `json:"popularity"`
	TrackNumber      int         `json:"track_number"`
	DiscNumber       int         `json:"disc_number"`
	IsLocal          bool        `json:"is_local"`
	Explicit         bool        `json:"explicit"`
}

type Album struct {
	ReleaseDatePrecision string      `json:"release_date_precision"`
	ExternalUrls         ExternalURL `json:"external_urls"`
	Href                 string      `json:"href"`
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	ReleaseDate          string      `json:"release_date"`
	AlbumType            string      `json:"album_type"`
	Type                 string      `json:"type"`
	URI                  string      `json:"uri"`
	Artists              []Artist    `json:"artists"`
	AvailableMarkets     []string    `json:"available_markets"`
	Images               []Image     `json:"images"`
	TotalTracks          int         `json:"total_tracks"`
}

type Artist struct {
	ExternalUrls ExternalURL `json:"external_urls"`
	Href         string      `json:"href"`
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	URI          string      `json:"uri"`
}

type ExternalID struct {
	ISRC string `json:"isrc"`
}

type ExternalURL struct {
	Spotify string `json:"spotify"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type AlbumItem struct {
	ReleaseDatePrecision string      `json:"release_date_precision"`
	ExternalURLs         ExternalURL `json:"external_urls"`
	Href                 string      `json:"href"`
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	ReleaseDate          string      `json:"release_date"`
	AlbumType            string      `json:"album_type"`
	Type                 string      `json:"type"`
	URI                  string      `json:"uri"`
	Artists              []Artist    `json:"artists"`
	Images               []Image     `json:"images"`
	TotalTracks          int         `json:"total_tracks"`
	IsPlayable           bool        `json:"is_playable"`
}
