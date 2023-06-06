package lastfm

type LastFMResponse interface {
	Append(l LastFMResponse) error
	GetTotalPages() int
	GetTotalFetched() int
}
