package collages

import (
	"bytes"
	"context"
	"errors"
	"image"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/cache"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/config"
	"github.com/SongStitch/song-stitch/internal/constants"
)

type LastFMArtist struct {
	Mbid      string `json:"mbid"`
	URL       string `json:"url"`
	Playcount string `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Name   string               `json:"name"`
	Images []lastfm.LastFMImage `json:"image"`
}

type LastFMTopArtists struct {
	TopArtists struct {
		Attr    lastfm.LastFMUser `json:"@attr"`
		Artists []LastFMArtist    `json:"artist"`
	} `json:"topartists"`
}

func (a *LastFMTopArtists) Append(l lastfm.LastFMResponse) error {
	if artists, ok := l.(*LastFMTopArtists); ok {
		a.TopArtists.Artists = append(a.TopArtists.Artists, artists.TopArtists.Artists...)
		return nil
	}
	return errors.New("type LastFMResponse is not a LastFMTopArtists")
}

func (a *LastFMTopArtists) TotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopArtists.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopArtists) TotalFetched() int {
	return len(a.TopArtists.Artists)
}

func GenerateCollageForArtist(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
	displayOptions DisplayOptions,
) (image.Image, *bytes.Buffer, error) {
	config := config.GetConfig()
	if count > config.MaxImages.Artists {
		return nil, nil, constants.ErrTooManyImages
	}
	artists, err := getArtists(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, nil, err
	}

	return CreateCollage(ctx, artists, displayOptions)
}

func getArtists(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
) ([]CollageElement, error) {
	result, err := lastfm.GetLastFmResponse[*LastFMTopArtists](
		ctx,
		constants.ARTIST,
		username,
		period,
		count,
	)
	if err != nil {
		return nil, err
	}
	r := *result
	cacheCount := 0
	logger := zerolog.Ctx(ctx)

	elements := make([]CollageElement, len(r.TopArtists.Artists))
	var wg sync.WaitGroup
	wg.Add(len(elements))

	start := time.Now()
	for i, lastfmArtist := range r.TopArtists.Artists {
		go func(i int, lastfmArtist LastFMArtist) {
			defer wg.Done()
			artist := parseLastfmArtist(ctx, lastfmArtist, imageSize, &cacheCount)
			artist.Image, err = DownloadImageWithRetry(ctx, artist.ImageUrl)
			if err != nil {
				logger.Error().
					Err(err).
					Str("imageUrl", artist.ImageUrl).
					Msg("Error downloading image")
			}

			elements[i] = CollageElement{Image: artist.Image, Parameters: artist.Parameters()}
		}(i, lastfmArtist)
	}
	wg.Wait()
	logger.Info().
		Int("cacheCount", cacheCount).
		Str("username", username).
		Int("totalCount", count).
		Dur("duration", time.Since(start)).
		Str("method", "artist").
		Msg("Image URLs fetched")
	return elements, nil
}

func parseLastfmArtist(ctx context.Context, artist LastFMArtist, imageSize string, cacheCount *int) Artist {
	logger := zerolog.Ctx(ctx)
	newArtist := Artist{
		Name:      artist.Name,
		Playcount: artist.Playcount,
		Mbid:      artist.Mbid,
		ImageSize: imageSize,
	}

	imageCache := cache.GetImageUrlCache()
	if cacheEntry, ok := imageCache.Get(newArtist.Identifier()); ok {
		newArtist.ImageUrl = cacheEntry.Url
		(*cacheCount)++
		return newArtist
	}

	// last.fm api doesn't return images for artists, so we can fetch the images from the website directly
	id, err := lastfm.GetImageIdForArtist(ctx, artist.URL)
	if err != nil {
		logger.Error().
			Err(err).
			Str("artist", artist.Name).
			Str("artistUrl", artist.URL).
			Msg("Error getting image url for artist")
		return newArtist
	}
	newArtist.ImageUrl = "https://lastfm.freetls.fastly.net/i/u/300x300/" + id
	return newArtist
}

type Artist struct {
	Name      string
	Playcount string
	Image     image.Image
	ImageUrl  string
	Mbid      string
	ImageSize string
	Url       string
}

func (a *Artist) Identifier() string {
	if a.Mbid != "" {
		return a.Mbid + a.ImageSize
	}
	return a.Url + a.ImageSize
}

func (a *Artist) CacheEntry() cache.CacheEntry {
	return cache.CacheEntry{Url: a.ImageUrl, Album: ""}
}

func (a *Artist) Parameters() map[string]string {
	return map[string]string{
		"artist":    a.Name,
		"playcount": a.Playcount,
	}
}
