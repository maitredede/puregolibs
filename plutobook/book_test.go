package plutobook

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
