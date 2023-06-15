package generator

import (
	"image"

	"github.com/SongStitch/song-stitch/internal/cache"
)

type Drawable interface {
	GetImage() *image.Image
	GetParameters() map[string]string
}
type Downloadable interface {
	GetImageUrl() string
	SetImage(*image.Image)
	GetIdentifier() string
	GetCacheEntry() cache.CacheEntry
}
