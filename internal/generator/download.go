package generator

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

const (
	jpgFileType = ".jpg"
	gifFileType = ".gif"
)

func DownloadImage[T Downloadable](a T) error {
	url := a.GetImageUrl()
	if len(url) == 0 {
		// Skip album art if it doesn't exist
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
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
		a.SetImage(&img)
		return err
	} else if strings.ToLower(extension) == gifFileType {
		img, err := gif.Decode(ioBody)
		if err != nil {
			return err
		}
		a.SetImage(&img)
		return err
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		img, _, err := image.Decode(bytes.NewReader(body))
		a.SetImage(&img)
		return err

	}
}

func DownloadImages[T Downloadable](ctx context.Context, entities []T) error {

	var wg sync.WaitGroup
	wg.Add(len(entities))

	for i := range entities {
		entity := &entities[i]
		// download each image in a separate goroutine
		go func(entity *T) {
			defer wg.Done()
			err := DownloadImage(*entity)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Str("imageUrl", (*entity).GetImageUrl()).Msg("Error downloading image")
			}
		}(entity)
	}

	// wait for all downloads to finish
	wg.Wait()

	return nil
}
