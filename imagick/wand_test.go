package imagick

import (
	"bytes"
	_ "embed"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed sample.png
	samplePNG []byte
)

func TestWandCreate(t *testing.T) {
	magickEnv := MagickWandGenesis()
	t.Cleanup(magickEnv.Terminus)

	w, newWandErr := magickEnv.NewMagickWand()
	assert.NoError(t, newWandErr, "new wand should not return error")

	closeErr := w.Close()
	assert.NoError(t, closeErr, "first close should not raise error")

	expectedErr := w.Close()
	assert.Error(t, expectedErr, "second close should raise invalid wand error")
}

func TestMakeProgressivePNG(t *testing.T) {
	magickEnv := MagickWandGenesis()
	t.Cleanup(magickEnv.Terminus)

	w, err := magickEnv.NewMagickWand()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	t.Logf("source blob size: %d", len(samplePNG))
	if err := w.ReadImageBlob(samplePNG); err != nil {
		t.Fatalf("read image error: %v", err)
	}
	if err := w.SetImageFormat("PNG"); err != nil {
		t.Fatalf("image format error: %v", err)
	}
	if err := w.SetImageInterlaceScheme(PNGInterlace); err != nil {
		t.Fatalf("image interlace scheme error: %v", err)
	}
	blob, err := w.GetImageBlob()
	if err != nil {
		t.Fatalf("get image blob error: %v", err)
	}
	t.Logf("new blob size: %d", len(blob))

	//try to decode with golang png
	reader := bytes.NewReader(blob)
	img, err := png.Decode(reader)
	if err != nil {
		t.Fatalf("go png decode failed: %v", err)
	}
	t.Logf("go img bounds: %v", img.Bounds())
}
