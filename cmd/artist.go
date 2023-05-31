package main

import (
	"image"
	"log"
	"strconv"
)

type LastFMArtist struct {
	Images    []LastFMImage `json:"image"`
	Mbid      string        `json:"mbid"`
	URL       string        `json:"url"`
	Playcount string        `json:"playcount"`
	Attr      struct {
		Rank string `json:"rank"`
	} `json:"@attr"`
	Name string `json:"name"`
}

type LastFMTopArtists struct {
	TopArtists struct {
		Artists []LastFMArtist `json:"artist"`
		Attr    LastFMUser     `json:"@attr"`
	} `json:"topartists"`
}

func (a *LastFMTopArtists) Append(l LastFMResponse) {

	if artists, ok := l.(*LastFMTopArtists); ok {
		a.TopArtists.Artists = append(a.TopArtists.Artists, artists.TopArtists.Artists...)
		return
	}
	log.Println("Error: LastFMResponse is not a LastFMTopArtists")
}
func (a *LastFMTopArtists) GetTotalPages() int {
	totalPages, _ := strconv.Atoi(a.TopArtists.Attr.TotalPages)
	return totalPages
}

func (a *LastFMTopArtists) GetTotalFetched() int {
	return len(a.TopArtists.Artists)
}

func getArtists(username string, period Period, count int, imageSize string) ([]*Artist, error) {

	// We get the artist images from the top albums, since the top artists endpoint doesn't return images
	// We instead display the image for the most played album for the artist
	albumResult, err := getLastFmResponse[*LastFMTopAlbums](ALBUM, username, period, 400, imageSize)
	if err != nil {
		return nil, err
	}
	albums := *albumResult
	artistImageMap := make(map[string]string)
	for _, album := range albums.TopAlbums.Albums {
		if album.Artist.Mbid == "" {
			continue
		}
		// Since the order if albumResult is by playcount, we can assume if the artist is already in the map
		// that means it is pointing to the most played album
		if _, ok := artistImageMap[album.Artist.Mbid]; !ok {
			for _, image := range album.Images {
				if image.Size == imageSize {
					artistImageMap[album.Artist.Mbid] = image.Link
				}
			}
		}
	}

	result, err := getLastFmResponse[*LastFMTopArtists](ARTIST, username, period, count, imageSize)
	if err != nil {
		return nil, err
	}
	r := *result

	artists := make([]*Artist, len(r.TopArtists.Artists))
	for i, artist := range r.TopArtists.Artists {
		newArtist := &Artist{
			Name:      artist.Name,
			Playcount: artist.Playcount,
		}

		val, ok := artistImageMap[artist.Mbid]
		newArtist.ImageUrl = val
		if !ok {
			log.Println("No image found for artist:", artist.Name, "with mbid", artist.Mbid)
		}

		artists[i] = newArtist
	}
	return artists, nil
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
