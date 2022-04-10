package db

import "image"

type BitmapLoader interface {
	GetTileBitmap(string) (image.Image, error)
}
