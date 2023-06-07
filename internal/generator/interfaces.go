package generator

import "image"

type Drawable interface {
	GetImage() *image.Image
	GetParameters() map[string]string
}
type Downloadable interface {
	GetImageUrl() string
	SetImage(*image.Image)
}
