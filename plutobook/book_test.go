package plutobook

import (
	"bytes"
	_ "embed"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed basic.html
	basicHtml string
)

func TestBookCreate(t *testing.T) {
	size := PageSize{
		Width:  100,
		Height: 100,
	}
	margins := PageMargins{
		Top:    1,
		Right:  1,
		Bottom: 1,
		Left:   1,
	}
	b, err := NewBook(size, margins, MediaTypeScreen)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBookSizes(t *testing.T) {
	size := PageSize{
		Width:  101,
		Height: 102,
	}
	margins := PageMargins{
		Top:    12,
		Right:  34,
		Bottom: 56,
		Left:   78,
	}
	media := MediaTypeScreen
	b, err := NewBook(size, margins, media)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	savedPageSize := b.GetPageSize()
	t.Logf("book pageSize: ctor=%+v, actual=%+v", size, savedPageSize)
	assert.EqualExportedValues(t, size, savedPageSize)
	savedMargins := b.GetPageMargins()
	t.Logf("book margins: ctor=%+v, actual=%+v", margins, savedMargins)
	assert.EqualExportedValues(t, margins, savedMargins)
	savedMediaType := b.GetMediaType()
	t.Logf("book media: ctor=%+v, actual=%+v", media, savedMediaType)
	assert.EqualValues(t, media, savedMediaType)
}

func TestBookHtmlSampleToPNG(t *testing.T) {
	size := PageSize{
		Width:  1900,
		Height: 1200,
	}
	margins := PageMargins{
		Top:    0,
		Right:  0,
		Bottom: 0,
		Left:   0,
	}
	media := MediaTypeScreen
	b, err := NewBook(size, margins, media)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	if err := b.LoadHTML(basicHtml, "", "", ""); err != nil {
		t.Fatal(err)
	}

	w := b.GetDocumentWidth()
	h := b.GetDocumentHeight()
	t.Logf("docSize: w=%v h=%v", w, h)

	tmp := t.TempDir()
	file := filepath.Join(tmp, "output.png")
	if err := b.WriteToPNG(file, int(w), int(h)); err != nil {
		t.Fatal(err)
	}
	input, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()

	img, err := png.Decode(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("img size: %v %v", img.Bounds().Dx(), img.Bounds().Dy())
}

func TestBookHtmlSampleToPNGStream(t *testing.T) {
	size := PageSize{
		Width:  1900,
		Height: 1200,
	}
	margins := PageMargins{
		Top:    0,
		Right:  0,
		Bottom: 0,
		Left:   0,
	}
	media := MediaTypeScreen
	b, err := NewBook(size, margins, media)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	if err := b.SetCustomResourceFetcher(DefaultHttpLoader); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := b.LoadHTML(basicHtml, "", "", ""); err != nil {
		t.Fatal(err)
	}

	w := b.GetDocumentWidth()
	h := b.GetDocumentHeight()
	t.Logf("docSize: w=%v h=%v", w, h)

	output := &bytes.Buffer{}
	if err := b.WriteToPNGStream(output, 0, 0); err != nil {
		t.Fatal(err)
	}
	t.Logf("png size: %d", output.Len())

	bin := output.Bytes()
	input := bytes.NewReader(bin)
	img, err := png.Decode(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("img size: %v %v", img.Bounds().Dx(), img.Bounds().Dy())
}

func TestBookUrlToPngStream(t *testing.T) {
	size := PageSize{
		Width:  1280,
		Height: 1080,
	}
	margins := PageMargins{
		Top:    0,
		Right:  0,
		Bottom: 0,
		Left:   0,
	}
	media := MediaTypeScreen
	b, err := NewBook(size, margins, media)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	// if err := b.SetCustomResourceFetcher(DefaultHttpLoader); err != nil {
	// 	t.Log(err)
	// 	t.Fail()
	// }

	if err := b.LoadURL("https://github.com", "", ""); err != nil {
		t.Fatal(err)
	}

	w := int(b.GetDocumentWidth())
	h := int(b.GetDocumentHeight())

	output := &bytes.Buffer{}
	if err := b.WriteToPNGStream(output, w, h); err != nil {
		t.Fatal(err)
	}
	t.Logf("png size: %d", output.Len())
}
