package collages

import (
	"context"
	"errors"
	"image"
	"strconv"
	"sync"

	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/SongStitch/song-stitch/internal/generator"
)

type LastFMArtist struct {
	Images    []lastfm.LastFMImage `json:"image"`
	Mbid      string               `json:"mbid"`
	URL       string               `json:"url"`
	Playcount string               `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Name string `json:"name"`
}

type LastFMTopArtists struct {
	TopArtists struct {
		Artists []LastFMArtist    `json:"artist"`
		Attr    lastfm.LastFMUser `json:"@attr"`
	} `json:"topartists"`
}

func (a *LastFMTopArtists) Append(l lastfm.LastFMResponse) error {

	if artists, ok := l.(*LastFMTopArtists); ok {
		a.TopArtists.Artists = append(a.TopArtists.Artists, artists.TopArtists.Artists...)
		return nil
	}
	return errors.New("type LastFMResponse is not a LastFMTopArtists")
}
func (a *LastFMTopArtists) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopArtists.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopArtists) GetTotalFetched() int {
	return len(a.TopArtists.Artists)
}

func GenerateCollageForArtist(ctx context.Context, username string, period constants.Period, count int, imageSize string, displayOptions generator.DisplayOptions) (*image.Image, error) {
	if count > 100 {
		return nil, constants.ErrTooManyImages
	}
	artists, err := getArtists(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}

	generator.DownloadImages(ctx, artists)

	return generator.CreateCollage(ctx, artists, displayOptions)
}

func getArtists(ctx context.Context, username string, period constants.Period, count int, imageSize string) ([]*Artist, error) {

	result, err := lastfm.GetLastFmResponse[*LastFMTopArtists](ctx, constants.ARTIST, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	r := *result

	artists := make([]*Artist, len(r.TopArtists.Artists))
	var wg sync.WaitGroup
	wg.Add(len(artists))
	for i, artist := range r.TopArtists.Artists {
		newArtist := &Artist{
			Name:      artist.Name,
			Playcount: artist.Playcount,
		}

		// last.fm api doesn't return images for artists, so we can fetch the images from the website directly
		go func(url string) {
			defer wg.Done()
			id, err := lastfm.GetImageIdForArtist(ctx, url)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Str("artist", artist.Name).Str("artistUrl", url).Msg("Error getting image url for artist")
				return
			}
			newArtist.ImageUrl = "https://lastfm.freetls.fastly.net/i/u/300x300/" + id
		}(artist.URL)
		artists[i] = newArtist
	}
	wg.Wait()
	return artists, nil
}

type Artist struct {
	Name      string
	Playcount string
	Image     image.Image
	ImageUrl  string
}

func (a *Artist) GetImageUrl() string {
	return a.ImageUrl
}

func (a *Artist) SetImage(img *image.Image) {
	a.Image = *img
}

func (a *Artist) GetImage() *image.Image {
	return &a.Image
}

func (a *Artist) GetParameters() map[string]string {
	return map[string]string{
		"artist":    a.Name,
		"playcount": a.Playcount,
	}
}
