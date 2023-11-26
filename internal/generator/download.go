package generator

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SongStitch/song-stitch/internal/cache"
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

func DownloadImageWithRetry(ctx context.Context, entity Downloadable) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = DownloadImage(ctx, entity)
		if err == nil {
			return nil
		}
		zerolog.Ctx(ctx).
			Warn().
			Err(err).
			Str("imageUrl", entity.GetImageUrl()).
			Msg("Error downloading image")
		delay := backoffSchedule[i]
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			continue
		}
	}
	return fmt.Errorf("failed to download image after %d retries: %w", maxRetries, err)
}

func DownloadImage(ctx context.Context, entity Downloadable) error {
	url := entity.GetImageUrl()
	if len(url) == 0 {
		// Skip album art if it doesn't exist
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	ioBody := resp.Body
	extension, err := getExtension(url)
	if err != nil {
		return err
	}

	if strings.ToLower(extension) == jpgFileType {
		img, err := jpeg.Decode(ioBody)
		if err != nil {
			return err
		}
		entity.SetImage(img)
		return err
	} else if strings.ToLower(extension) == gifFileType {
		img, err := gif.Decode(ioBody)
		if err != nil {
			return err
		}
		entity.SetImage(img)
		return err
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		img, _, err := image.Decode(bytes.NewReader(body))
		entity.SetImage(img)
		return err

	}
}

func DownloadImages[T Downloadable](ctx context.Context, entities []T) error {

	logger := zerolog.Ctx(ctx)
	var wg sync.WaitGroup
	wg.Add(len(entities))

	start := time.Now()
	for i := range entities {
		entity := &entities[i]
		// download each image in a separate goroutine
		go func(entity *T) {
			defer wg.Done()
			err := DownloadImageWithRetry(ctx, *entity)
			if err != nil {
				logger.Error().
					Err(err).
					Str("imageUrl", (*entity).GetImageUrl()).
					Msg("Error downloading image")
			}
			cache := cache.GetImageUrlCache()
			cache.Set((*entity).GetIdentifier(), (*entity).GetCacheEntry())
		}(entity)
	}

	// wait for all downloads to finish
	wg.Wait()
	logger.Info().Dur("duration", time.Since(start)).Msg("Downloaded images")

	return nil
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
