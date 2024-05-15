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
	"github.com/SongStitch/song-stitch/internal/clients"
	"github.com/SongStitch/song-stitch/internal/clients/lastfm"
	"github.com/SongStitch/song-stitch/internal/clients/spotify"
	"github.com/SongStitch/song-stitch/internal/config"
	"github.com/SongStitch/song-stitch/internal/constants"
)

type LastfmTrack struct {
	Artist struct {
		URL  string `json:"url"`
		Name string `json:"name"`
		Mbid string `json:"mbid"`
	} `json:"artist"`
	Mbid     string `json:"mbid"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Duration string `json:"duration"`
	Attr     struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Playcount string               `json:"playcount"`
	Images    []lastfm.LastfmImage `json:"image"`
}

type LastfmTopTracks struct {
	TopTracks struct {
		Attr   lastfm.LastfmUser `json:"@attr"`
		Tracks []LastfmTrack     `json:"track"`
	} `json:"toptracks"`
}

func GetElementsForTrack(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
	displayOptions DisplayOptions,
) ([]CollageElement, error) {
	config := config.GetConfig()
	if count > config.MaxImages.Tracks {
		return nil, constants.ErrTooManyImages
	}
	tracks, err := getTracks(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func getLastfmTracks(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
) ([]LastfmTrack, error) {
	tracks := []LastfmTrack{}
	totalPages := 0

	handler := func(data []byte) (int, int, error) {
		var lastfmTopTracks LastfmTopTracks
		err := json.Unmarshal(data, &lastfmTopTracks)
		if err != nil {
			return 0, 0, err
		}
		tracks = append(tracks, lastfmTopTracks.TopTracks.Tracks...)
		if totalPages == 0 {
			total, err := strconv.Atoi(lastfmTopTracks.TopTracks.Attr.TotalPages)
			if err != nil {
				return 0, 0, err
			}
			totalPages = total
		}
		return len(tracks), totalPages, nil
	}
	err := lastfm.GetLastFmResponse(ctx, constants.TRACK, username, period, count, handler)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}
func getTracks(
	ctx context.Context,
	username string,
	period constants.Period,
	count int,
	imageSize string,
) ([]CollageElement, error) {
	tracks, err := getLastfmTracks(ctx, username, period, count)
	if err != nil {
		return nil, err
	}
	cacheCount := 0
	logger := zerolog.Ctx(ctx)

	elements := make([]CollageElement, len(tracks))

	var wg sync.WaitGroup
	wg.Add(len(tracks))
	start := time.Now()
	for i, track := range tracks {
		go func(i int, lastfmTrack LastfmTrack) {
			defer wg.Done()
			track := parseLastfmTrack(ctx, lastfmTrack, imageSize, &cacheCount)
			track.Image, err = DownloadImageWithRetry(ctx, track.ImageUrl)
			if err != nil {
				logger.Error().
					Err(err).
					Str("imageUrl", track.ImageUrl).
					Msg("Error downloading image")
			}
			elements[i] = CollageElement{
				Parameters: track.Parameters(),
				Image:      track.Image,
			}
		}(i, track)

	}
	wg.Wait()
	logger.Info().
		Int("cacheCount", cacheCount).
		Str("username", username).
		Int("totalCount", count).
		Dur("duration", time.Since(start)).
		Str("method", "track").
		Msg("Image URLs fetched")

	return elements, nil
}

func parseLastfmTrack(
	ctx context.Context,
	track LastfmTrack,
	imageSize string,
	cacheCount *int,
) Track {
	logger := zerolog.Ctx(ctx)
	newTrack := Track{
		Name:      track.Name,
		Artist:    track.Artist.Name,
		Playcount: track.Playcount,
		Mbid:      track.Mbid,
		ImageSize: imageSize,
	}

	imageCache := cache.GetImageUrlCache()
	if cacheEntry, ok := imageCache.Get(newTrack.Identifier()); ok {
		newTrack.ImageUrl = cacheEntry.Url
		newTrack.Album = cacheEntry.Album
		(*cacheCount)++
		return newTrack
	}

	trackInfo, err := getTrackInfo(ctx, newTrack.Name, newTrack.Artist, imageSize)
	if err != nil {
		logger.Error().
			Str("track", newTrack.Name).
			Str("artist", newTrack.Artist).
			Err(err).
			Msg("Error getting track info")
		return newTrack
	}
	newTrack.ImageUrl = trackInfo.ImageUrl
	newTrack.Album = trackInfo.AlbumName
	return newTrack
}

func getTrackInfo(
	ctx context.Context,
	trackName string,
	artistName string,
	imageSize string,
) (clients.TrackInfo, error) {
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
	return trackInfo, constants.ErrNoImageFound
}

func getTrackInfoFromLastFm(
	trackName string,
	artistName string,
	imageSize string,
) (clients.TrackInfo, error) {
	result, err := lastfm.GetTrackInfo(trackName, artistName, imageSize)
	if err != nil {
		return result, err
	}
	if result.ImageUrl == "" {
		return result, constants.ErrNoImageFound
	}
	return result, nil
}

func getTrackInfoFromSpotify(
	ctx context.Context,
	trackName string,
	artistName string,
) (clients.TrackInfo, error) {
	client, err := spotify.GetSpotifyClient()
	if err != nil {
		return clients.TrackInfo{}, err
	}
	result, err := client.GetTrackInfo(ctx, trackName, artistName)
	if err != nil {
		return clients.TrackInfo{}, err
	}
	if result.ImageUrl == "" {
		return clients.TrackInfo{}, constants.ErrNoImageFound
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

func (t *Track) Identifier() string {
	if t.Mbid != "" {
		return t.Mbid + t.ImageSize
	}
	return t.Name + t.Artist + t.ImageSize
}

func (t *Track) CacheEntry() cache.CacheEntry {
	return cache.CacheEntry{Url: t.ImageUrl, Album: t.Album}
}

func (t *Track) Parameters() map[string]string {
	return map[string]string{
		"artist":    t.Artist,
		"track":     t.Name,
		"playcount": t.Playcount,
		"album":     t.Album,
	}
}
