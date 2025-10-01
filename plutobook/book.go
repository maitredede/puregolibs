package plutobook

import (
	"errors"
)

type Book struct {
	ptr uintptr
}

func NewBook(pageSize PageSize, margins PageMargins, mediaType MediaType) (*Book, error) {
	libInit()

	ptr := libCreate(pageSize, margins, mediaType)
	if ptr == 0 {
		msg := libGetErrorMessage()
		libClearErrorMessage()
		return nil, errors.New(msg)
	}
	b := Book{
		ptr: ptr,
	}
	return &b, nil
}

func (b *Book) Close() error {
	libInit()

	if b.ptr == 0 {
		return ErrBookIsClosed
	}
	libDestroy(uintptr(b.ptr))
	b.ptr = 0
	return nil
}

func (b *Book) GetPageSize() PageSize {
	libInit()
	if b.ptr == 0 {
		return PageSize{}
	}
	return libGetPageSize(b.ptr)
}

func (b *Book) GetPageSizeAt(index int) PageSize {
	libInit()
	if b.ptr == 0 {
		return PageSize{}
	}
	return libGetPageSizeAt(b.ptr, index)
}

func (b *Book) GetMediaType() MediaType {
	libInit()
	if b.ptr == 0 {
		return MediaTypePrint
	}
	ret := libGetMediaType(b.ptr)
	return MediaType(ret)
}

func (b *Book) GetPageMargins() PageMargins {
	libInit()
	if b.ptr == 0 {
		return PageMargins{}
	}
	ret := libGetPageMargins(b.ptr)
	return ret
}

func (b *Book) LoadURL(url string, userStyle string, userScript string) error {
	panic("TODO")
}

func (b *Book) LoadHTML(html string, userStyle string, userScript string, baseUrl string) error {
	panic("TODO")
}

func (b *Book) RenderPage(canvas *Canvas, pageIndex int) error {
	panic("TODO")
}

func (b *Book) WriteToPDF(file string) error {
	panic("TODO")
}

func (b *Book) WriteToPDFRange(file string, pageStart, pageEnd, pageStep int) error {
	panic("TODO")
}

func (b *Book) SetCustomResourceFetcher(fetcher CustomResourceFetcher) error {
	panic("TODO")
}
