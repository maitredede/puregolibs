//go:build linux

package drm

import "fmt"

type Fourcc uint32

func FourccValue(a, b, c, d uint8) Fourcc {
	return Fourcc(uint32(a) |
		uint32(b)<<8 |
		uint32(c)<<16 |
		uint32(d)<<24)
}

var fourccStrings = make(map[Fourcc]string)

func buildFourcc(a, b, c, d uint8, name string) Fourcc {
	val := FourccValue(a, b, c, d)

	str := fmt.Sprintf("%s (%s)", name, string([]byte{a, b, c, d}))
	fourccStrings[val] = str
	return val
}

func (v Fourcc) String() string {
	if s, ok := fourccStrings[v]; ok {
		return s
	}
	return fmt.Sprintf("0x%08x", uint32(v))
}

// TODO place here used fourcc codes
var (
	DRM_FORMAT_XRGB8888 = buildFourcc('X', 'R', '2', '4', "XRGB8888")
	DRM_FORMAT_XBGR8888 = buildFourcc('X', 'B', '2', '4', "XBGR8888")
	DRM_FORMAT_RGBX8888 = buildFourcc('R', 'X', '2', '4', "RGBX8888")
	DRM_FORMAT_BGRX8888 = buildFourcc('B', 'X', '2', '4', "BGRX8888")

	DRM_FORMAT_ARGB8888 = buildFourcc('A', 'R', '2', '4', "ARGB8888")
	DRM_FORMAT_ABGR8888 = buildFourcc('A', 'B', '2', '4', "ABGR8888")
	DRM_FORMAT_RGBA8888 = buildFourcc('R', 'A', '2', '4', "RGBA8888")
	DRM_FORMAT_BGRA8888 = buildFourcc('B', 'A', '2', '4', "BGRA8888")
)
