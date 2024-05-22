package collages

import (
	"context"
	"encoding/json"
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

type LastfmArtist struct {
	Mbid      string `json:"mbid"`
	URL       string `json:"url"`
	Playcount string `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Name   string               `json:"name"`
	Images []lastfm.LastfmImage `json:"image"`
}

type LastfmTopArtists struct {
	TopArtists struct {
		Attr    lastfm.LastfmUser `json:"@attr"`
		Artists []LastfmArtist    `json:"artist"`
	} `json:"topartists"`
}

func GetElementsForArtist(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
	displayOptions DisplayOptions,
) ([]CollageElement, error) {
	config := config.GetConfig()
	if count > config.MaxImages.Artists {
		return nil, constants.ErrTooManyImages
	}
	artists, err := getArtists(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func getLastfmArtists(ctx context.Context, username string, period constants.Period, count int) ([]LastfmArtist, error) {
	artists := []LastfmArtist{}
	totalPages := 0

	handler := func(data []byte) (int, int, error) {
		var lastfmTopArtists LastfmTopArtists
		err := json.Unmarshal(data, &lastfmTopArtists)
		if err != nil {
			return 0, 0, err
		}
		artists = append(artists, lastfmTopArtists.TopArtists.Artists...)
		if totalPages == 0 {
			total, err := strconv.Atoi(lastfmTopArtists.TopArtists.Attr.TotalPages)
			if err != nil {
				return 0, 0, err
			}
			totalPages = total
		}
		return len(artists), totalPages, nil
	}
	err := lastfm.GetLastFmResponse(ctx, constants.ARTIST, username, period, count, handler)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func getArtists(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
) ([]CollageElement, error) {
	artists, err := getLastfmArtists(ctx, username, period, count)
	if err != nil {
		return nil, err
	}
	cacheCount := 0
	logger := zerolog.Ctx(ctx)

	elements := make([]CollageElement, len(artists))
	var wg sync.WaitGroup
	wg.Add(len(elements))

	start := time.Now()
	for i, lastfmArtist := range artists {
		go func(i int, lastfmArtist LastfmArtist) {
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

func parseLastfmArtist(ctx context.Context, artist LastfmArtist, imageSize string, cacheCount *int) Artist {
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
	imageCache.Set(newArtist.Identifier(), newArtist.CacheEntry())
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
