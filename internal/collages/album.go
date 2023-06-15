package collages

import (
	"context"
	"errors"
	"image"
	"strconv"
	"sync"

	"github.com/SongStitch/song-stitch/internal/cache"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/clients/spotify"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/SongStitch/song-stitch/internal/generator"
	"github.com/SongStitch/song-stitch/internal/models"
	"github.com/rs/zerolog"
)

type LastFMAlbum struct {
	Artist struct {
		URL        string `json:"url"`
		ArtistName string `json:"name"`
		Mbid       string `json:"mbid"`
	} `json:"artist"`
	Images    []lastfm.LastFMImage `json:"image"`
	Mbid      string               `json:"mbid"`
	URL       string               `json:"url"`
	Playcount string               `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	AlbumName string `json:"name"`
}

type LastFMTopAlbums struct {
	TopAlbums struct {
		Albums []LastFMAlbum     `json:"album"`
		Attr   lastfm.LastFMUser `json:"@attr"`
	} `json:"topalbums"`
}

func (a *LastFMTopAlbums) Append(l lastfm.LastFMResponse) error {
	if albums, ok := l.(*LastFMTopAlbums); ok {
		a.TopAlbums.Albums = append(a.TopAlbums.Albums, albums.TopAlbums.Albums...)
		return nil
	}
	return errors.New("type LastFMResponse is not a LastFMTopAlbums")
}

func (a *LastFMTopAlbums) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopAlbums.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopAlbums) GetTotalFetched() int {
	return len(a.TopAlbums.Albums)
}

func GenerateCollageForAlbum(ctx context.Context, username string, period constants.Period, count int, imageSize string, displayOptions generator.DisplayOptions) (*image.Image, error) {
	albums, err := getAlbums(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}

	generator.DownloadImages(ctx, albums)

	return generator.CreateCollage(ctx, albums, displayOptions)
}

func getAlbums(ctx context.Context, username string, period constants.Period, count int, imageSize string) ([]*Album, error) {
	result, err := lastfm.GetLastFmResponse[*LastFMTopAlbums](ctx, constants.ALBUM, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	r := *result
	cacheCount := 0

	logger := zerolog.Ctx(ctx)

	albums := make([]*Album, len(r.TopAlbums.Albums))
	var wg sync.WaitGroup
	wg.Add(len(r.TopAlbums.Albums))
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
			if cacheEntry, ok := imageCache.Get(newAlbum.GetIdentifier()); ok {
				newAlbum.ImageUrl = cacheEntry.Url
				cacheCount++
				return
			}
			albumInfo, err := getAlbumInfo(ctx, album, imageSize)
			if err != nil {
				logger.Error().Str("album", album.AlbumName).Str("artist", album.Artist.ArtistName).Err(err).Msg("Error getting album info")
				return
			}
			albums[i].ImageUrl = albumInfo.ImageUrl
		}(i, album)
	}
	wg.Wait()
	logger.Info().Int("cacheCount", cacheCount).Str("username", username).Int("totalCount", count).Msg("Albums fetched from cache")
	return albums, nil
}

func getAlbumInfo(ctx context.Context, album LastFMAlbum, imageSize string) (*models.AlbumInfo, error) {
	for _, image := range album.Images {
		if image.Size == imageSize && image.Link != "" {
			return &models.AlbumInfo{ImageUrl: image.Link}, nil
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
	ImageUrl  string
	Image     image.Image
	Mbid      string
	ImageSize string
}

func (a *Album) GetImageUrl() string {
	return a.ImageUrl
}

func (a *Album) SetImage(img *image.Image) {
	a.Image = *img
}

func (a *Album) GetIdentifier() string {
	if a.Mbid != "" {
		return a.Mbid + a.ImageSize
	}
	return a.Name + a.Artist + a.ImageSize
}

func (a *Album) GetCacheEntry() cache.CacheEntry {
	return cache.CacheEntry{Url: a.ImageUrl, Album: ""}
}

func (a *Album) GetImage() *image.Image {
	return &a.Image
}

func (a *Album) GetParameters() map[string]string {
	return map[string]string{
		"artist":    a.Artist,
		"album":     a.Name,
		"playcount": a.Playcount,
	}
}
