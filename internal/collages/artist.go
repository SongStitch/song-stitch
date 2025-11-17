package collages

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/cache"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/config"
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
	period lastfm.Period,
	count int,
	imageSize string,
	displayOptions DisplayOptions,
	jobChan chan<- CollageElement,
) error {
	config := config.GetConfig()
	if count > config.MaxImages.Artists {
		return lastfm.ErrTooManyImages
	}
	return getArtists(ctx, username, period, count, imageSize, jobChan)
}

func getLastfmArtists(
	ctx context.Context,
	username string,
	period lastfm.Period,
	count int,
) ([]LastfmArtist, error) {
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
	err := lastfm.GetLastFmResponse(ctx, lastfm.MethodArtist, username, period, count, handler)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func getArtists(
	ctx context.Context,
	username string,
	period lastfm.Period,
	count int,
	imageSize string,
	jobChan chan<- CollageElement,
) error {
	artists, err := getLastfmArtists(ctx, username, period, count)
	if err != nil {
		return err
	}

	var cacheCount int64
	logger := zerolog.Ctx(ctx)

	var wg sync.WaitGroup

	start := time.Now()
	for i, lastfmArtist := range artists {
		wg.Add(1)

		go func(i int, lastfmArtist LastfmArtist) {
			defer wg.Done()

			artist := parseLastfmArtist(ctx, lastfmArtist, imageSize, &cacheCount)

			img, ext, imgErr := DownloadImageWithRetry(ctx, artist.ImageUrl)
			if imgErr != nil {
				logger.Error().
					Err(imgErr).
					Str("imageUrl", artist.ImageUrl).
					Msg("Error downloading image")
			}

			jobChan <- CollageElement{
				Index:      i,
				ImageBytes: img,
				ImageExt:   ext,
				Parameters: artist.Parameters(),
			}
		}(i, lastfmArtist)
	}

	wg.Wait()

	logger.Info().
		Int64("cacheCount", atomic.LoadInt64(&cacheCount)).
		Str("username", username).
		Int("totalCount", count).
		Dur("duration", time.Since(start)).
		Str("method", "artist").
		Msg("Image URLs fetched")

	return nil
}

func parseLastfmArtist(
	ctx context.Context,
	artist LastfmArtist,
	imageSize string,
	cacheCount *int64,
) Artist {
	logger := zerolog.Ctx(ctx)

	newArtist := Artist{
		Name:      artist.Name,
		Playcount: artist.Playcount,
		Mbid:      artist.Mbid,
		Url:       artist.URL,
		ImageSize: imageSize,
	}

	key := newArtist.Identifier()
	if key != "" {
		imageCache := cache.GetImageUrlCache()
		if cacheEntry, ok := imageCache.Get(key); ok {
			newArtist.ImageUrl = cacheEntry.Url
			atomic.AddInt64(cacheCount, 1)
			logger.Info().Msg("Image URL found in cache")
			return newArtist
		}
	}

	idOrURL, err := lastfm.GetImageIdForArtist(ctx, artist.Name, artist.Mbid)
	if err != nil {
		logger.Error().
			Err(err).
			Str("artist", artist.Name).
			Str("artistUrl", artist.URL).
			Msg("Error getting image url for artist")
		return newArtist
	}

	imageURL := lastfm.BuildArtistImageURL(idOrURL)
	if imageURL == "" {
		logger.Warn().Msg("No image URL found for artist")
		return newArtist
	}

	newArtist.ImageUrl = imageURL

	if key != "" {
		imageCache := cache.GetImageUrlCache()
		imageCache.Set(key, newArtist.CacheEntry())
	}

	return newArtist
}

type Artist struct {
	Name      string
	Playcount string
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
