package lastfm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/SongStitch/song-stitch/internal/clients"
	"github.com/SongStitch/song-stitch/internal/config"
)

type LastfmImage struct {
	Size string `json:"size"`
	Link string `json:"#text"`
}

type LastfmUser struct {
	User       string `json:"user"`
	TotalPages string `json:"totalPages"`
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	Total      string `json:"total"`
}

func getMethodForCollageType(collageType Method) string {
	switch collageType {
	case MethodAlbum:
		return "user.gettopalbums"
	case MethodArtist:
		return "user.gettopartists"
	case MethodTrack:
		return "user.gettoptracks"
	default:
		return ""
	}
}

type CleanError struct {
	errStr string
}

func (e CleanError) Error() string {
	return e.errStr
}

// strip sensitive information from error message
func cleanError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()
	modified := apiKeyRedactionRegex.ReplaceAllString(errStr, "$1")
	return CleanError{errStr: modified}
}

var (
	defaultHTTPClient = &http.Client{
		Timeout: 60 * time.Second,
	}

	apiKeyRedactionRegex = regexp.MustCompile(`([&?])api_key=[^&]+(&|\b)`)

	defaultUserAgent = "songstitch/1.0 (+https://songstitch.art)"
)

func GetLastFmResponse(
	ctx context.Context,
	collageType Method,
	username string,
	period Period,
	count int,
	handler func(data []byte) (fetched int, totalPages int, err error),
) error {
	cfg := config.GetConfig()
	endpoint := cfg.Lastfm.Endpoint
	apiKey := cfg.Lastfm.APIKey

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Fetching Last.fm data")

	method := getMethodForCollageType(collageType)
	if method == "" {
		return fmt.Errorf("unsupported collage type: %v", collageType)
	}

	const maxPerPage = 500

	totalFetched := 0
	page := 1

	for count > totalFetched {
		logger.Info().
			Int("page", page).
			Int("totalFetched", totalFetched).
			Int("count", count).
			Msg("Fetching Last.fm page")

		limit := min(count-totalFetched, maxPerPage)

		u, err := url.Parse(endpoint)
		if err != nil {
			logger.Error().Err(err).Msg("invalid Last.fm endpoint")
			return fmt.Errorf("invalid lastfm endpoint: %w", err)
		}

		q := u.Query()
		q.Set("user", username)
		q.Set("method", method)
		q.Set("period", string(period))
		q.Set("limit", strconv.Itoa(limit))
		q.Set("page", strconv.Itoa(page))
		q.Set("api_key", apiKey)
		q.Set("format", "json")
		u.RawQuery = q.Encode()

		body, err := func() ([]byte, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				return nil, err
			}

			start := time.Now()
			res, err := defaultHTTPClient.Do(req)
			if res != nil {
				defer res.Body.Close()
			}

			logger.Info().
				Dur("duration", time.Since(start)).
				Str("method", method).
				Int("status", statusCodeOrZero(res)).
				Msg("Last.fm request completed")

			if err != nil {
				return nil, cleanError(err)
			}

			if res.StatusCode == http.StatusNotFound {
				return nil, ErrUserNotFound
			}

			if res.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("lastfm unexpected status code: %d", res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			return body, nil
		}()
		if err != nil {
			return err
		}

		fetched, totalPages, err := handler(body)
		if err != nil {
			return err
		}

		totalFetched = fetched

		if totalPages == page || totalPages == 0 {
			break
		}

		page++
	}

	return nil
}

func statusCodeOrZero(res *http.Response) int {
	if res == nil {
		return 0
	}
	return res.StatusCode
}

type GetTrackInfoResponse struct {
	Track struct {
		Album struct {
			AlbumName string        `json:"title"`
			Images    []LastfmImage `json:"image"`
		} `json:"album"`
	} `json:"track"`
}

func GetTrackInfo(
	trackName string,
	artistName string,
	imageSize string,
) (clients.TrackInfo, error) {
	cfg := config.GetConfig()
	endpoint := cfg.Lastfm.Endpoint
	apiKey := cfg.Lastfm.APIKey

	u, err := url.Parse(endpoint)
	if err != nil {
		return clients.TrackInfo{}, fmt.Errorf("invalid lastfm endpoint: %w", err)
	}

	q := u.Query()
	q.Set("track", trackName)
	q.Set("artist", artistName)
	q.Set("method", "track.getInfo")
	q.Set("api_key", apiKey)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return clients.TrackInfo{}, err
	}
	res, err := defaultHTTPClient.Do(req)
	if err != nil {
		return clients.TrackInfo{}, cleanError(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return clients.TrackInfo{}, errors.New("track not found")
	}

	if res.StatusCode != http.StatusOK {
		return clients.TrackInfo{}, fmt.Errorf(
			"lastfm track.getInfo unexpected status code: %d",
			res.StatusCode,
		)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return clients.TrackInfo{}, err
	}

	var response GetTrackInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return clients.TrackInfo{}, err
	}

	for _, image := range response.Track.Album.Images {
		if image.Size == imageSize {
			return clients.TrackInfo{
				AlbumName: response.Track.Album.AlbumName,
				ImageUrl:  image.Link,
			}, nil
		}
	}

	return clients.TrackInfo{}, errors.New("no image found for requested size")
}

// fanart.tv image info.
type fanartImage struct {
	URL   string `json:"url"`
	Likes string `json:"likes"`
}

// fanart.tv artist response.
type fanartArtistResponse struct {
	Name             string        `json:"name"`
	MBID             string        `json:"mbid_id"`
	ArtistThumb      []fanartImage `json:"artistthumb"`
	ArtistBackground []fanartImage `json:"artistbackground"`
	HDMusicLogo      []fanartImage `json:"hdmusiclogo"`
	MusicBanner      []fanartImage `json:"musicbanner"`
	MusicLogo        []fanartImage `json:"musiclogo"`
}

// Wikipedia pageimages response.
type wikipediaQueryResponse struct {
	Query struct {
		Pages map[string]struct {
			Title     string `json:"title"`
			Thumbnail struct {
				Source string `json:"source"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnail"`
		} `json:"pages"`
	} `json:"query"`
}

// Deezer artist search response.
type deezerArtistSearchResponse struct {
	Data []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXL     string `json:"picture_xl"`
	} `json:"data"`
	Total int `json:"total"`
}

func fetchArtistArtworkFromFanart(
	ctx context.Context,
	mbid, apiKey string,
) (*fanartArtistResponse, error) {
	if mbid == "" {
		return nil, nil
	}
	if apiKey == "" {
		return nil, fmt.Errorf("fanart.tv API key is empty")
	}

	endpoint := fmt.Sprintf("https://webservice.fanart.tv/v3/music/%s?api_key=%s", mbid, apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fanart.tv status: %s", resp.Status)
	}

	var payload fanartArtistResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// bestArtistThumbURL picks the "best" artist thumb URL from a fanart.tv response.
// In case of multiple candidates, the one with the highest "likes" count is chosen
func bestArtistThumbURL(f *fanartArtistResponse) string {
	if f == nil {
		return ""
	}

	candidates := [][]fanartImage{
		f.ArtistThumb,
		f.ArtistBackground,
		f.HDMusicLogo,
		f.MusicLogo,
		f.MusicBanner,
	}

	for _, group := range candidates {
		if len(group) == 0 {
			continue
		}

		best := bestByLikes(group)
		if best.URL != "" {
			return best.URL
		}
	}

	return ""
}

func bestByLikes(images []fanartImage) fanartImage {
	if len(images) == 0 {
		return fanartImage{}
	}

	best := images[0]
	bestLikes := parseLikes(best.Likes)

	for _, img := range images[1:] {
		if l := parseLikes(img.Likes); l > bestLikes {
			bestLikes = l
			best = img
		}
	}

	return best
}

func parseLikes(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

// Hack to normalise quote characters for Wikipedia titles
var wikiTitleReplacer = strings.NewReplacer(
	"’", "'",
	"‘", "'",
	"“", `"`,
	"”", `"`,
)

// normaliseForWikipediaTitle normalises quote characters to match typical Wikipedia page titles
func normaliseForWikipediaTitle(s string) string {
	return wikiTitleReplacer.Replace(s)
}

const wikipediaThumbSize = "600"

func fetchArtistImageFromWikipedia(ctx context.Context, artistName string) (string, error) {
	artistName = strings.TrimSpace(artistName)
	if artistName == "" {
		return "", nil
	}

	artistName = normaliseForWikipediaTitle(artistName)

	q := url.Values{}
	q.Set("action", "query")
	q.Set("format", "json")
	q.Set("prop", "pageimages")
	q.Set("piprop", "thumbnail")
	q.Set("pithumbsize", wikipediaThumbSize)
	q.Set("redirects", "1")
	q.Set("titles", artistName)

	endpoint := "https://en.wikipedia.org/w/api.php?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wikipedia status: %s", resp.Status)
	}

	var payload wikipediaQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	for _, page := range payload.Query.Pages {
		if page.Thumbnail.Source != "" {
			return page.Thumbnail.Source, nil
		}
	}

	return "", nil
}

// fetchArtistImageFromDeezer fetches an artist image from the Deezer search API as a fallback.
func fetchArtistImageFromDeezer(ctx context.Context, artistName string) (string, error) {
	artistName = strings.TrimSpace(artistName)
	if artistName == "" {
		return "", nil
	}

	q := url.Values{}
	q.Set("q", artistName)

	endpoint := "https://api.deezer.com/search/artist?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deezer status: %s", resp.Status)
	}

	var payload deezerArtistSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if len(payload.Data) == 0 {
		return "", errors.New("empty payload")
	}

	getImageUrl := func(payload deezerArtistSearchResponse) string {
		a := payload.Data[0]
		switch {
		case a.PictureXL != "":
			return a.PictureXL
		case a.PictureBig != "":
			return a.PictureBig
		case a.PictureMedium != "":
			return a.PictureMedium
		case a.PictureSmall != "":
			return a.PictureSmall
		case a.Picture != "":
			return a.Picture
		default:
			return ""
		}
	}

	url := getImageUrl(payload)
	isValidUrl := func(url string) bool {
		// format is https://cdn-images.dzcdn.net/images/artist/<id>/1000x1000-000000-80-0-0.jpg
		// if <id> is missing, we return false

		s := strings.Split(url, "artist/")
		if len(s) != 2 {
			return false
		}

		// <id>/1000x1000-000000-80-0-0.jpg
		s2 := strings.Split(s[1], "/")
		if len(s2) != 2 {
			return false
		}

		return s2[0] != ""
	}

	valid := isValidUrl(url)
	if valid {
		url = strings.Replace(url, "1000x1000", "300x300", 1)
		return url, nil
	}
	return "", nil
}

func GetImageIdForArtist(ctx context.Context, artistName string, mbid string) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str("artistName", artistName).Str("mbid", mbid).Logger()
	cfg := config.GetConfig()

	lookupFanart := func(ctx context.Context, mbid string) string {
		fa, err := fetchArtistArtworkFromFanart(ctx, mbid, cfg.Fanart.APIKey)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("fanart.tv artwork lookup failed")
			return ""
		}

		artURL := bestArtistThumbURL(fa)
		if artURL != "" {
			logger.Info().
				Str("artwork_url", artURL).
				Msg("Successfully resolved artist artwork URL from fanart.tv")
			return artURL
		}
		return ""
	}

	if mbid != "" {
		url := lookupFanart(ctx, mbid)
		if url != "" {
			return url, nil
		}
	}

	deezerURL, err := fetchArtistImageFromDeezer(ctx, artistName)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Deezer artwork lookup failed")
	} else if deezerURL != "" {
		logger.Info().
			Str("artwork_url", deezerURL).
			Msg("Successfully resolved artist artwork URL from Deezer (no MBID)")
		return deezerURL, nil
	}

	wikiURL, err := fetchArtistImageFromWikipedia(ctx, artistName)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Wikipedia artwork lookup failed")
	} else if wikiURL != "" {
		logger.Info().
			Str("artwork_url", wikiURL).
			Msg("Successfully resolved artist artwork URL from Wikipedia (no MBID)")
		return wikiURL, nil
	}

	return "", fmt.Errorf("no image found")
}

// BuildArtistImageURL normalises either a raw URL or a legacy Last.fm image ID
// into a full HTTP URL.
func BuildArtistImageURL(idOrURL string) string {
	idOrURL = strings.TrimSpace(idOrURL)
	if idOrURL == "" {
		return ""
	}

	if strings.HasPrefix(idOrURL, "http://") || strings.HasPrefix(idOrURL, "https://") {
		return idOrURL
	}

	return "https://lastfm.freetls.fastly.net/i/u/300x300/" + idOrURL
}
