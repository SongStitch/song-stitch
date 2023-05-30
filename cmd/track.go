package main

import (
	"image"
	"log"
	"strconv"
)

type LastFMTrack struct {
	Mbid   string        `json:"mbid"`
	Name   string        `json:"name"`
	Images []LastFMImage `json:"image"`
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
		Tracks []LastFMTrack `json:"track"`
		Attr   LastFMUser    `json:"@attr"`
	} `json:"toptracks"`
}

func (a *LastFMTopTracks) Append(l LastFMResponse) {
	if tracks, ok := l.(*LastFMTopTracks); ok {
		a.TopTracks.Tracks = append(a.TopTracks.Tracks, tracks.TopTracks.Tracks...)
		return
	}
	log.Println("Error: LastFMResponse is not a LastFMTopAlbums")
}

func (a *LastFMTopTracks) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopTracks.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopTracks) GetTotalFetched() int {
	return len(a.TopTracks.Tracks)
}

func getTracks(username string, period Period, count int, imageSize string) ([]*Track, error) {
	result, err := getLastFmResponse[*LastFMTopTracks](TRACK, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	r := *result

	tracks := make([]*Track, len(r.TopTracks.Tracks))
	for i, track := range r.TopTracks.Tracks {
		newTrack := &Track{
			Name:      track.Name,
			Artist:    track.Artist.Name,
			Playcount: track.Playcount,
		}

		for _, image := range track.Images {
			if image.Size == imageSize {
				newTrack.ImageUrl = image.Link
			}
		}

		tracks[i] = newTrack
	}
	return tracks, nil
}

type Track struct {
	Name      string
	Artist    string
	Playcount string
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
	}
}
