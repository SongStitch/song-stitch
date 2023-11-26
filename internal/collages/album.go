package collages

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/SongStitch/song-stitch/internal/generator"
	"github.com/rs/zerolog"
)

type LastFMAlbum struct {
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
	Images    []lastfm.LastFMImage `json:"image"`
}

type LastFMTopAlbums struct {
	TopAlbums struct {
		Attr   lastfm.LastFMUser `json:"@attr"`
		Albums []LastFMAlbum     `json:"album"`
	} `json:"topalbums"`
}

func (a *LastFMTopAlbums) Append(l lastfm.LastFMResponse) error {
	if albums, ok := l.(*LastFMTopAlbums); ok {
		a.TopAlbums.Albums = append(a.TopAlbums.Albums, albums.TopAlbums.Albums...)
		return nil
	}
	return errors.New("type LastFMResponse is not a LastFMTopAlbums")
}

func (a *LastFMTopAlbums) TotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopAlbums.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopAlbums) TotalFetched() int {
	return len(a.TopAlbums.Albums)
}

func GenerateCollageForAlbum(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
	displayOptions generator.DisplayOptions,
) (image.Image, *bytes.Buffer, error) {
	config := config.GetConfig()
	if count > config.MaxImages.Albums {
		return nil, nil, constants.ErrTooManyImages
	}
	albums, err := getAlbums(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, nil, err
	}

	generator.DownloadImages(ctx, albums)

	return generator.CreateCollage(ctx, albums, displayOptions)
}

func getAlbums(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
) ([]*Album, error) {
	result, err := lastfm.GetLastFmResponse[*LastFMTopAlbums](
		ctx,
		constants.ALBUM,
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

	albums := make([]*Album, len(r.TopAlbums.Albums))
	var wg sync.WaitGroup
	wg.Add(len(r.TopAlbums.Albums))
	start := time.Now()
	for i, album := range r.TopAlbums.Albums {
		go func(i int, album LastFMAlbum) {
			defer wg.Done()
			newAlbum := &Album{
				Name:      album.AlbumName,
				Artist:    album.Artist.ArtistName,
				Playcount: album.Playcount,
				Mbid:      album.Mbid,
				ImageSize: imageSize,
			}
			albums[i] = newAlbum

			imageCache := cache.GetImageUrlCache()
			if cacheEntry, ok := imageCache.Get(newAlbum.Identifier()); ok {
				newAlbum.imageUrl = cacheEntry.Url
				cacheCount++
				return
			}
			albumInfo, err := getAlbumInfo(ctx, album, imageSize)
			if err != nil {
				logger.Error().
					Str("album", album.AlbumName).
					Str("artist", album.Artist.ArtistName).
					Err(err).
					Msg("Error getting album info")
				return
			}
			albums[i].imageUrl = albumInfo.ImageUrl
		}(i, album)
	}
	wg.Wait()
	logger.Info().
		Int("cacheCount", cacheCount).
		Str("username", username).
		Int("totalCount", count).
		Dur("duration", time.Since(start)).
		Str("method", "album").
		Msg("Image URLs fetched")
	return albums, nil
}

func getAlbumInfo(
	ctx context.Context,
	album LastFMAlbum,
	imageSize string,
) (*clients.AlbumInfo, error) {
	for _, image := range album.Images {
		if image.Size == imageSize && image.Link != "" {
			return &clients.AlbumInfo{ImageUrl: image.Link}, nil
		}
	}
	client, err := spotify.GetSpotifyClient()
	if err != nil {
		return nil, err
	}
	albumInfo, err := client.GetAlbumInfo(ctx, album.AlbumName, album.Artist.ArtistName)
	if err != nil {
		return nil, err
	}
	return albumInfo, nil
}

type Album struct {
	Name      string
	Artist    string
	Playcount string
	imageUrl  string
	image     image.Image
	Mbid      string
	ImageSize string
}

func (a *Album) ImageUrl() string {
	return a.imageUrl
}

func (a *Album) SetImage(img image.Image) {
	a.image = img
}

func (a *Album) Identifier() string {
	if a.Mbid != "" {
		return a.Mbid + a.ImageSize
	}
	return a.Name + a.Artist + a.ImageSize
}

func (a *Album) CacheEntry() cache.CacheEntry {
	return cache.CacheEntry{Url: a.imageUrl, Album: ""}
}

func (a *Album) Image() image.Image {
	return a.image
}

func (a *Album) Parameters() map[string]string {
	return map[string]string{
		"artist":    a.Artist,
		"album":     a.Name,
		"playcount": a.Playcount,
	}
}

func (a *Album) ClearImage() {
	a.image = nil
}
