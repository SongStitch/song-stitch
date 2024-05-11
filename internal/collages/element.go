package collages

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	jpgFileType = ".jpg"
	gifFileType = ".gif"
)
const maxRetries = 3

var (
	backoffSchedule = []time.Duration{
		200 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
)

func DownloadImageWithRetry(ctx context.Context, url string) (image.Image, error) {
	var err error
	for i := 0; i < maxRetries; i++ {
		img, err := DownloadImage(ctx, url)
		if err == nil {
			return img, nil
		}
		zerolog.Ctx(ctx).
			Warn().
			Err(err).
			Str("imageUrl", url).
			Msg("Error downloading image")
		delay := backoffSchedule[i]
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			continue
		}
	}
	return nil, fmt.Errorf("failed to download image after %d retries: %w", maxRetries, err)
}

func DownloadImage(ctx context.Context, url string) (image.Image, error) {
	if url == "" {
		return nil, errors.New("empty url")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	ioBody := resp.Body
	extension, err := getExtension(url)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(extension) {
	case jpgFileType:
		img, err := jpeg.Decode(ioBody)
		if err != nil {
			return nil, err
		}
		return img, nil
	case gifFileType:
		img, err := gif.Decode(ioBody)
		if err != nil {
			return nil, err
		}
		return img, nil
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		img, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		return img, nil
	}
}

func getExtension(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	// Split the path component of the URL into a slice of path elements
	pathElements := strings.Split(parsedURL.Path, "/")

	// The last element of the path should be the filename
	fileName := pathElements[len(pathElements)-1]

	// Extract the file extension from the filename
	ext := filepath.Ext(fileName)
	return ext, nil
}
