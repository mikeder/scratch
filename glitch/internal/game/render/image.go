package render

import (
	"image"
	"image/color"
)

func NewImage(w, h int, c color.RGBA) *image.NRGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{w, h}

	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			switch {
			case x < w/2 && y < h/2: // upper left quadrant
				img.Set(x, y, c)
			case x >= w/2 && y >= h/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}
	return img
}
