package generator

import (
	"image"

	"github.com/SongStitch/song-stitch/internal/cache"
)

type Drawable interface {
	Image() image.Image
	Parameters() map[string]string
	ClearImage()
}

type Downloadable interface {
	ImageUrl() string
	SetImage(image.Image)
	Identifier() string
	CacheEntry() cache.CacheEntry
}
