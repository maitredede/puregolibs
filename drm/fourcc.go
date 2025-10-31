//go:build linux

package drm

func FourccCode(a, b, c, d uint8) uint32 {
	return uint32(a) |
		uint32(b)<<8 |
		uint32(c)<<16 |
		uint32(d)<<24
}

// TODO place here used fourcc codes
var (
	DRM_FORMAT_ARGB8888 = FourccCode('A', 'R', '2', '4')
)
