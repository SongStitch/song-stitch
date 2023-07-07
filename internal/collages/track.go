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
	"github.com/SongStitch/song-stitch/internal/clients/spotify"
	"github.com/SongStitch/song-stitch/internal/constants"
	"github.com/SongStitch/song-stitch/internal/generator"
	"github.com/SongStitch/song-stitch/internal/models"
)

type LastFMTrack struct {
	Mbid   string               `json:"mbid"`
	Name   string               `json:"name"`
	Images []lastfm.LastFMImage `json:"image"`
	Artist struct {
		URL  string `json:"url"`
		Name string `json:"name"`
		Mbid string `json:"mbid"`
	} `json:"artist"`
	URL      string `json:"url"`
	Duration string `json:"duration"`
	Attr     struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Playcount string `json:"playcount"`
}

type LastFMTopTracks struct {
	TopTracks struct {
		Tracks []LastFMTrack     `json:"track"`
		Attr   lastfm.LastFMUser `json:"@attr"`
	} `json:"toptracks"`
}

func (t *LastFMTopTracks) Append(l lastfm.LastFMResponse) error {
	if tracks, ok := l.(*LastFMTopTracks); ok {
		t.TopTracks.Tracks = append(t.TopTracks.Tracks, tracks.TopTracks.Tracks...)
		return nil
	}
	return errors.New("type LastFMResponse is not a LastFMTopAlbums")
}

func (t *LastFMTopTracks) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(t.TopTracks.Attr.TotalPages)
	return totalPages
}

func (t *LastFMTopTracks) GetTotalFetched() int {
	return len(t.TopTracks.Tracks)
}

func GenerateCollageForTrack(ctx context.Context, username string, period constants.Period, count int, imageSize string, displayOptions generator.DisplayOptions) (*image.Image, *bytes.Buffer, error) {
	if count > 25 {
		return nil, nil, constants.ErrTooManyImages
	}
	tracks, err := getTracks(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, nil, err
	}

	generator.DownloadImages(ctx, tracks)

	return generator.CreateCollage(ctx, tracks, displayOptions)
}

func getTracks(ctx context.Context, username string, period constants.Period, count int, imageSize string) ([]*Track, error) {
	result, err := lastfm.GetLastFmResponse[*LastFMTopTracks](ctx, constants.TRACK, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	r := *result
	cacheCount := 0
	logger := zerolog.Ctx(ctx)

	tracks := make([]*Track, len(r.TopTracks.Tracks))

	var wg sync.WaitGroup
	wg.Add(len(r.TopTracks.Tracks))
	start := time.Now()
	for i, track := range r.TopTracks.Tracks {
		newTrack := &Track{
			Name:      track.Name,
			Artist:    track.Artist.Name,
			Playcount: track.Playcount,
			Mbid:      track.Mbid,
			ImageSize: imageSize,
		}
		tracks[i] = newTrack

		imageCache := cache.GetImageUrlCache()
		if cacheEntry, ok := imageCache.Get(newTrack.GetIdentifier()); ok {
			newTrack.ImageUrl = cacheEntry.Url
			newTrack.Album = cacheEntry.Album
			cacheCount++
			wg.Done()
			continue
		}

		go func(trackName string, artistName string) {
			defer wg.Done()

			trackInfo, err := getTrackInfo(ctx, trackName, artistName, imageSize)
			if err != nil {
				logger.Error().Str("track", trackName).Str("artist", artistName).Err(err).Msg("Error getting track info")
				return
			}
			newTrack.ImageUrl = trackInfo.ImageUrl
			newTrack.Album = trackInfo.AlbumName
		}(track.Name, track.Artist.Name)

	}
	wg.Wait()
	logger.Info().Int("cacheCount", cacheCount).Str("username", username).Int("totalCount", count).Dur("duration", time.Since(start)).Str("method", "track").Msg("Image URLs fetched")

	return tracks, nil
}

func getTrackInfo(ctx context.Context, trackName string, artistName string, imageSize string) (*models.TrackInfo, error) {
	logger := zerolog.Ctx(ctx)

	trackInfo, err := getTrackInfoFromLastFm(trackName, artistName, imageSize)
	if err == nil {
		return trackInfo, nil
	}
	logger.Warn().Err(err).Msg("Error getting track info from lastfm")
	trackInfo, err = getTrackInfoFromSpotify(ctx, trackName, artistName)
	if err == nil {
		return trackInfo, nil
	}
	logger.Warn().Err(err).Msg("Error getting track info from spotify")
	return nil, constants.ErrNoImageFound
}

func getTrackInfoFromLastFm(trackName string, artistName string, imageSize string) (*models.TrackInfo, error) {
	result, err := lastfm.GetTrackInfo(trackName, artistName, imageSize)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, constants.ErrNoImageFound
	}
	if result.ImageUrl == "" {
		return nil, constants.ErrNoImageFound
	}
	return result, nil
}

func getTrackInfoFromSpotify(ctx context.Context, trackName string, artistName string) (*models.TrackInfo, error) {
	client, err := spotify.GetSpotifyClient()
	if err != nil {
		return nil, err
	}
	result, err := client.GetTrackInfo(ctx, trackName, artistName)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, constants.ErrNoImageFound
	}
	if result.ImageUrl == "" {
		return nil, constants.ErrNoImageFound
	}
	return result, nil
}

type Track struct {
	Name      string
	Artist    string
	Playcount string
	Album     string
	ImageUrl  string
	Image     image.Image
	Mbid      string
	ImageSize string
}

func (t *Track) GetImageUrl() string {
	return t.ImageUrl
}

func (t *Track) SetImage(img *image.Image) {
	t.Image = *img
}

func (t *Track) GetImage() *image.Image {
	return &t.Image
}

func (t *Track) GetIdentifier() string {
	if t.Mbid != "" {
		return t.Mbid + t.ImageSize
	}
	return t.Name + t.Artist + t.ImageSize
}
func (t *Track) GetCacheEntry() cache.CacheEntry {
	return cache.CacheEntry{Url: t.ImageUrl, Album: t.Album}
}

func (t *Track) GetParameters() map[string]string {
	return map[string]string{
		"artist":    t.Artist,
		"track":     t.Name,
		"playcount": t.Playcount,
		"album":     t.Album,
	}
}
