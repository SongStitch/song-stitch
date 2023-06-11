package collages

import (
	"context"
	"errors"
	"image"
	"strconv"
	"sync"

	"github.com/rs/zerolog"

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

func (a *LastFMTopTracks) Append(l lastfm.LastFMResponse) error {
	if tracks, ok := l.(*LastFMTopTracks); ok {
		a.TopTracks.Tracks = append(a.TopTracks.Tracks, tracks.TopTracks.Tracks...)
		return nil
	}
	return errors.New("type LastFMResponse is not a LastFMTopAlbums")
}

func (a *LastFMTopTracks) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopTracks.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopTracks) GetTotalFetched() int {
	return len(a.TopTracks.Tracks)
}

func GenerateCollageForTrack(ctx context.Context, username string, period constants.Period, count int, imageSize string, displayOptions generator.DisplayOptions) (image.Image, error) {
	if count > 25 {
		return nil, constants.ErrTooManyImages
	}
	tracks, err := getTracks(ctx, username, period, count, imageSize)
	if err != nil {
		return nil, err
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
	logger := zerolog.Ctx(ctx)

	tracks := make([]*Track, len(r.TopTracks.Tracks))
	spotifyClient, err := spotify.NewSpotifyClient()
	if err != nil {
		logger.Warn().Err(err).Msg("Error creating spotify client")
	}

	var wg sync.WaitGroup
	wg.Add(len(r.TopTracks.Tracks))
	for i, track := range r.TopTracks.Tracks {
		newTrack := &Track{
			Name:      track.Name,
			Artist:    track.Artist.Name,
			Playcount: track.Playcount,
		}

		go func(trackName string, artistName string) {
			defer wg.Done()
			trackInfo, err := getTrackInfoFromLastFm(trackName, artistName, imageSize)
			if err != nil {
				logger.Warn().Err(err).Msg("Error getting track info from lastfm")
				trackInfo, err = getTrackInfoFromSpotify(ctx, spotifyClient, trackName, artistName)
				if err != nil {
					logger.Warn().Err(err).Msg("Error getting track info from spotify")
					return
				}
			}

			newTrack.ImageUrl = trackInfo.ImageUrl
			newTrack.Album = trackInfo.AlbumName
		}(track.Name, track.Artist.Name)

		tracks[i] = newTrack
	}
	wg.Wait()
	return tracks, nil
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

func getTrackInfoFromSpotify(ctx context.Context, client *spotify.SpotifyClient, trackName string, artistName string) (*models.TrackInfo, error) {
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
}

func (a *Track) GetImageUrl() string {
	return a.ImageUrl
}

func (a *Track) SetImage(img *image.Image) {
	a.Image = *img
}

func (a *Track) GetImage() *image.Image {
	return &a.Image
}

func (a *Track) GetParameters() map[string]string {
	return map[string]string{
		"artist":    a.Artist,
		"track":     a.Name,
		"playcount": a.Playcount,
		"album":     a.Album,
	}
}
