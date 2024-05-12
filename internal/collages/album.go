package collages

import (
	"context"
	"encoding/json"
	"image"
	"strconv"
	"sync"
	"time"

	"github.com/SongStitch/song-stitch/internal/cache"
	"github.com/SongStitch/song-stitch/internal/clients"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/clients/spotify"
	"github.com/SongStitch/song-stitch/internal/config"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/rs/zerolog"
)

type LastfmAlbum struct {
	Artist struct {
		URL        string `json:"url"`
		ArtistName string `json:"name"`
		Mbid       string `json:"mbid"`
	} `json:"artist"`
	Mbid      string `json:"mbid"`
	URL       string `json:"url"`
	Playcount string `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	AlbumName string               `json:"name"`
	Images    []lastfm.LastfmImage `json:"image"`
}

type LastfmTopAlbums struct {
	TopAlbums struct {
		Attr   lastfm.LastfmUser `json:"@attr"`
		Albums []LastfmAlbum     `json:"album"`
	} `json:"topalbums"`
}

func GetElementsForAlbum(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
	displayOptions DisplayOptions,
) ([]CollageElement, error) {
	config := config.GetConfig()
	if count > config.MaxImages.Albums {
		return nil, constants.ErrTooManyImages
	}
	albums, err := getAlbums(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func getLastfmAlbums(ctx context.Context, username string, period constants.Period, count int) ([]LastfmAlbum, error) {
	albums := []LastfmAlbum{}
	totalPages := 0

	handler := func(data []byte) (int, int, error) {
		var lastfmTopAlbums LastfmTopAlbums
		err := json.Unmarshal(data, &lastfmTopAlbums)
		if err != nil {
			return 0, 0, err
		}
		albums = append(albums, lastfmTopAlbums.TopAlbums.Albums...)
		if totalPages == 0 {
			total, err := strconv.Atoi(lastfmTopAlbums.TopAlbums.Attr.TotalPages)
			if err != nil {
				return 0, 0, err
			}
			totalPages = total
		}
		return len(albums), totalPages, nil
	}
	err := lastfm.GetLastFmResponse(ctx, constants.ALBUM, username, period, count, handler)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func getAlbums(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
) ([]CollageElement, error) {
	albums, err := getLastfmAlbums(ctx, username, period, count)
	if err != nil {
		return nil, err
	}
	cacheCount := 0

	logger := zerolog.Ctx(ctx)

	elements := make([]CollageElement, len(albums))
	var wg sync.WaitGroup
	wg.Add(len(albums))
	start := time.Now()
	for i, lastfmAlbum := range albums {
		go func(i int, lastfmAlbum LastfmAlbum) {
			defer wg.Done()
			album := parseLastfmAlbum(ctx, lastfmAlbum, imageSize, &cacheCount)

			album.Image, err = DownloadImageWithRetry(ctx, album.ImageUrl)
			if err != nil {
				logger.Error().
					Err(err).
					Str("imageUrl", album.ImageUrl).
					Msg("Error downloading image")
			}
			elements[i] = CollageElement{
				Parameters: album.Parameters(),
				Image:      album.Image,
			}
		}(i, lastfmAlbum)
	}
	wg.Wait()
	logger.Info().
		Int("cacheCount", cacheCount).
		Str("username", username).
		Int("totalCount", count).
		Dur("duration", time.Since(start)).
		Str("method", "album").
		Msg("Image URLs fetched")
	return elements, nil
}

func parseLastfmAlbum(ctx context.Context, album LastfmAlbum, imageSize string, cacheCount *int) Album {
	logger := zerolog.Ctx(ctx)
	newAlbum := Album{
		Artist:    album.Artist.ArtistName,
		Name:      album.AlbumName,
		Playcount: album.Playcount,
		Mbid:      album.Mbid,
		ImageSize: imageSize,
	}

	imageCache := cache.GetImageUrlCache()
	if cacheEntry, ok := imageCache.Get(newAlbum.Identifier()); ok {
		newAlbum.ImageUrl = cacheEntry.Url
		(*cacheCount)++
		return newAlbum
	}
	albumInfo, err := getAlbumInfo(ctx, album, imageSize)
	if err != nil {
		logger.Error().
			Str("album", album.AlbumName).
			Str("artist", album.Artist.ArtistName).
			Err(err).
			Msg("Error getting album info")
		return newAlbum
	}
	newAlbum.ImageUrl = albumInfo.ImageUrl
	return newAlbum
}

func getAlbumInfo(
	ctx context.Context,
	album LastfmAlbum,
	imageSize string,
) (clients.AlbumInfo, error) {
	for _, image := range album.Images {
		if image.Size == imageSize && image.Link != "" {
			return clients.AlbumInfo{ImageUrl: image.Link}, nil
		}
	}
	client, err := spotify.GetSpotifyClient()
	if err != nil {
		return clients.AlbumInfo{}, err
	}
	albumInfo, err := client.GetAlbumInfo(ctx, album.AlbumName, album.Artist.ArtistName)
	if err != nil {
		return clients.AlbumInfo{}, err
	}
	return albumInfo, nil
}

type Album struct {
	Image     image.Image
	ImageUrl  string
	Artist    string
	Name      string
	Playcound string
	ImageSize string
	Mbid      string
	Playcount string
}

func (a *Album) Identifier() string {
	if a.Mbid != "" {
		return a.Mbid + a.ImageSize
	}
	return a.Name + a.Artist + a.ImageSize
}

func (a *Album) CacheEntry() cache.CacheEntry {
	return cache.CacheEntry{Url: a.ImageUrl, Album: ""}
}

func (a *Album) Parameters() map[string]string {
	return map[string]string{
		"artist":    a.Artist,
		"album":     a.Name,
		"playcount": a.Playcount,
	}
}
