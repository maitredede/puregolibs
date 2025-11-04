package main

import (
	"image"
	"image/color"
)

func grabPixelsXRGB8888(data []byte, rect image.Rectangle, dest *image.RGBA) {
	// XRGB8888 = 4 bytes par pixel (X, R, G, B)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			offset := (y*dest.Bounds().Dx() + x) * 4

			// XRGB8888: byte 0 = X (ignoré), 1 = R, 2 = G, 3 = B
			r := data[offset+1]
			g := data[offset+2]
			b := data[offset+3]

			dest.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
}
