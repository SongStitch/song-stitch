package lastfm

type LastFMResponse interface {
	Append(l LastFMResponse) error
	TotalPages() int
	TotalFetched() int
}
