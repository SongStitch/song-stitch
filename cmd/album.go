package main

import (
	"image"
	"log"
	"strconv"
)

type LastFMAlbum struct {
	Artist struct {
		URL        string `json:"url"`
		ArtistName string `json:"name"`
		Mbid       string `json:"mbid"`
	} `json:"artist"`
	Image     []LastFMImage `json:"image"`
	Mbid      string        `json:"mbid"`
	URL       string        `json:"url"`
	Playcount string        `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	AlbumName string `json:"name"`
}

type LastFMTopAlbums struct {
	TopAlbums struct {
		Album []LastFMAlbum `json:"album"`
		Attr  LastFMUser    `json:"@attr"`
	} `json:"topalbums"`
}

func (a *LastFMTopAlbums) Append(l LastFMResponse) {
	if albums, ok := l.(*LastFMTopAlbums); ok {
		a.TopAlbums.Album = append(a.TopAlbums.Album, albums.TopAlbums.Album...)
		return
	}
	log.Println("Error: LastFMResponse is not a LastFMTopAlbums")
}

func (a *LastFMTopAlbums) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopAlbums.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopAlbums) GetTotalFetched() int {
	return len(a.TopAlbums.Album)
}

func getAlbums(collageType CollageType, username string, period Period, count int, imageSize string) ([]*Album, error) {

	result, err := getLastFmResponse[*LastFMTopAlbums](collageType, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	r := *result

	albums := make([]*Album, len(r.TopAlbums.Album))
	for i, album := range r.TopAlbums.Album {
		newAlbum := &Album{
			Name:      album.AlbumName,
			Artist:    album.Artist.ArtistName,
			Playcount: album.Playcount,
		}

		for _, image := range album.Image {
			if image.Size == imageSize {
				newAlbum.ImageUrl = image.Link
			}
		}

		albums[i] = newAlbum
	}
	return albums, nil
}

type Album struct {
	Name      string
	Artist    string
	Playcount string
	ImageUrl  string
	Image     image.Image
}

func (a *Album) GetImageUrl() string {
	return a.ImageUrl
}

func (a *Album) SetImage(img *image.Image) {
	a.Image = *img
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

type Artist struct {
	Name      string
	Playcount string
	Image     image.Image
	ImageUrl  string
}

func (a *Artist) GetImageUrl() string {
	return a.ImageUrl
}

func (a *Artist) SetImage(img *image.Image) {
	a.Image = *img
}

func (a *Artist) GetImage() *image.Image {
	return &a.Image
}

func (a *Artist) GetParameters() map[string]string {
	return map[string]string{
		"artist":    a.Name,
		"playcount": a.Playcount,
	}
}
