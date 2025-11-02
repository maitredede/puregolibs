//go:build linux

package libevdi

type Mode struct {
	Width        int32
	Height       int32
	RefreshRate  int32
	BitsPerPixel int32
	PixelFormat  int32
}

func toEvdiMode(event drmEvdiModeChanged) Mode {
	mode := Mode{
		Width:        event.hdisplay,
		Height:       event.vdisplay,
		RefreshRate:  event.vrefresh,
		BitsPerPixel: event.bitsPerPixel,
		PixelFormat:  event.pixelFormat,
	}
	return mode
}
