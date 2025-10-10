package plutobook

import (
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
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
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}

	cUrl := strings.CString(url)
	cUserStyle := strings.CString(userStyle)
	cUserScript := strings.CString(userScript)
	ok := libLoadUrl(b.ptr, uintptr(unsafe.Pointer(cUrl)), uintptr(unsafe.Pointer(cUserStyle)), uintptr(unsafe.Pointer(cUserScript)))
	if !ok {
		return fmt.Errorf("error loading url %v", url)
	}
	return nil
}

func (b *Book) LoadHTML(html string, userStyle string, userScript string, baseUrl string) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}

	cHtml, lgHtml := strings.CStringL(html)
	cUserStyle := strings.CString(userStyle)
	cUserScript := strings.CString(userScript)
	cBaseUrl := strings.CString(baseUrl)
	ok := libLoadHtml(b.ptr, uintptr(unsafe.Pointer(cHtml)), int32(lgHtml), uintptr(unsafe.Pointer(cUserStyle)), uintptr(unsafe.Pointer(cUserScript)), uintptr(unsafe.Pointer(cBaseUrl)))
	if !ok {
		return errors.New("error loading html")
	}
	return nil
}

func (b *Book) RenderPage(canvas *ImageCanvas, pageIndex int) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}
	panic("TODO: Book.RenderPage")
}

func (b *Book) WriteToPDF(file string) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}
	panic("TODO: Book.WriteToPDF")
}

func (b *Book) WriteToPDFRange(file string, pageStart, pageEnd, pageStep int) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}
	panic("TODO: Book.WriteToPDFRange")
}

func (b *Book) SetCustomResourceFetcher(fetcher CustomResourceFetcher) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}
	panic("TODO: SetCustomResourceFetcher")
}

func (b *Book) GetDocumentWidth() float32 {
	libInit()
	if b.ptr == 0 {
		return 0
	}
	return libGetDocumentWidth(b.ptr)
}

func (b *Book) GetDocumentHeight() float32 {
	libInit()
	if b.ptr == 0 {
		return 0
	}
	return libGetDocumentHeight(b.ptr)
}

func (b *Book) WriteToPNG(file string, width, height int) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}

	cFile := strings.CString(file)
	ret := libWriteToPNG(b.ptr, uintptr(unsafe.Pointer(cFile)), int32(width), int32(height))
	if !ret {
		return errors.New("file write error")
	}
	return nil
}

func (b *Book) WriteToPNGStream(output io.Writer, width, height int) error {
	libInit()
	if b.ptr == 0 {
		return ErrBookIsClosed
	}
	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		return errors.New("closure not allocated")
	}
	defer ffi.ClosureFree(closure)

	// describe the closure's signature
	// plutobook_stream_status_t (*plutobook_stream_write_callback_t)(void* closure, const char* data, unsigned int length);
	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); status != ffi.OK {
		return fmt.Errorf("cif preparation failed: %v", status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(streamWriteCallback)

	stream := &streamWriteData{
		output: output,
		err:    nil,
	}
	// prepare the closure
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
		return fmt.Errorf("closure preparation failed: %v", status)
	}

	isOk := libWriteToPNGStream(b.ptr, uintptr(callback), uintptr(unsafe.Pointer(stream)), int32(width), int32(height))

	if !isOk {
		if stream.err != nil {
			return stream.err
		}
		return errors.New("png write failed")
	}
	return nil
}
