package main

import (
	"image"

	"github.com/maitredede/puregolibs/drm"
)

type pixelsGrabber func(src []byte, rect image.Rectangle, dest *image.RGBA)

var pixelsGrabberMap map[drm.Fourcc]pixelsGrabber = map[drm.Fourcc]pixelsGrabber{
	drm.DRM_FORMAT_XRGB8888: grabPixelsXRGB8888,
}
