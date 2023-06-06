package clients

type LastFMResponse interface {
	Append(l LastFMResponse) error
	GetTotalPages() int
	GetTotalFetched() int
}
