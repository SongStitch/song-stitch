package spotify

type TracksResponse struct {
	SearchResult SearchResult[TrackItem] `json:"tracks"`
}

type AlbumResponse struct {
	SearchResult SearchResult[AlbumItem] `json:"albums"`
}

type SearchResult[T any] struct {
	Href     string `json:"href"`
	Items    []T    `json:"items"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

type TrackItem struct {
	Album            Album       `json:"album"`
	Artists          []Artist    `json:"artists"`
	AvailableMarkets []string    `json:"available_markets"`
	DiscNumber       int         `json:"disc_number"`
	DurationMs       int         `json:"duration_ms"`
	Explicit         bool        `json:"explicit"`
	ExternalIds      ExternalID  `json:"external_ids"`
	ExternalUrls     ExternalURL `json:"external_urls"`
	Href             string      `json:"href"`
	ID               string      `json:"id"`
	IsLocal          bool        `json:"is_local"`
	Name             string      `json:"name"`
	Popularity       int         `json:"popularity"`
	PreviewURL       string      `json:"preview_url"`
	TrackNumber      int         `json:"track_number"`
	Type             string      `json:"type"`
	URI              string      `json:"uri"`
}

type Album struct {
	AlbumType            string      `json:"album_type"`
	Artists              []Artist    `json:"artists"`
	AvailableMarkets     []string    `json:"available_markets"`
	ExternalUrls         ExternalURL `json:"external_urls"`
	Href                 string      `json:"href"`
	ID                   string      `json:"id"`
	Images               []Image     `json:"images"`
	Name                 string      `json:"name"`
	ReleaseDate          string      `json:"release_date"`
	ReleaseDatePrecision string      `json:"release_date_precision"`
	TotalTracks          int         `json:"total_tracks"`
	Type                 string      `json:"type"`
	URI                  string      `json:"uri"`
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
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type AlbumItem struct {
	AlbumType            string      `json:"album_type"`
	Artists              []Artist    `json:"artists"`
	ExternalURLs         ExternalURL `json:"external_urls"`
	Href                 string      `json:"href"`
	ID                   string      `json:"id"`
	Images               []Image     `json:"images"`
	IsPlayable           bool        `json:"is_playable"`
	Name                 string      `json:"name"`
	ReleaseDate          string      `json:"release_date"`
	ReleaseDatePrecision string      `json:"release_date_precision"`
	TotalTracks          int         `json:"total_tracks"`
	Type                 string      `json:"type"`
	URI                  string      `json:"uri"`
}
