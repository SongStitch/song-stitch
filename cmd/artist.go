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
		// Since the order if albumResult is by playcount, we can assume if the artist is already in the map
		// that means it is pointing to the most played album
		key := album.Artist.Mbid
		if album.Artist.Mbid == "" {
			// If an artist name is already in the map, then we know if had a higher playcount
			if _, ok := artistImageMap[album.Artist.ArtistName]; ok {
				continue
			}
			log.Println("Artist", album.Artist.ArtistName, "has no mbid, using name instead")
			key = album.Artist.ArtistName
		}
		if _, ok := artistImageMap[key]; !ok {
			for _, image := range album.Images {
				if image.Size == imageSize {
					artistImageMap[key] = image.Link
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

		if artist.Mbid == "" {
			log.Println("Artist", artist.Name, "has no mbid, searching for it")
			artistsFromSearch, err := searchArtist(artist.Name)
			if err != nil {
				log.Println("Error searching for artist", artist.Name, ":", err)
			} else {
				for _, artistFromSearch := range artistsFromSearch.Results.ArtistMatches.Artists {
					log.Println("Found artist:", artistFromSearch.Name, "with mbid", artistFromSearch.Mbid)
					if artistFromSearch.Name == artist.Name {
						artist.Mbid = artistFromSearch.Mbid
						break
					}
				}
			}
		}
		val, ok := artistImageMap[artist.Mbid]
		if !ok {
			log.Println("No image found for artist", artist.Name, "with mbid", artist.Mbid, "trying with name")
			val, ok = artistImageMap[artist.Name]
			if !ok {
				log.Println("No image found for artist", artist.Name, "with mbid", artist.Mbid)
			}

		}
		newArtist.ImageUrl = val
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
