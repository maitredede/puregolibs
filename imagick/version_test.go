package imagick

import "testing"

func TestVersion(t *testing.T) {
	cv := GetMagickVersion()
	t.Logf("imagemagick version: %v", cv)
	wv := GetMagickVersionWand()
	t.Logf("imagemagick version (from wand): %v", wv)
}
