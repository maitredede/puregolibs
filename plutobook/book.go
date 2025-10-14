package plutobook

import (
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
)

type bookPtr unsafe.Pointer

type Book struct {
	ptr bookPtr

	fetcherClosure *ffi.Closure
	fetcher        CustomResourceFetcher
}

func NewBook(pageSize PageSize, margins PageMargins, mediaType MediaType) (*Book, error) {
	libInit()

	ptr := libCreate(pageSize, margins, mediaType)
	if ptr == nil {
		msg := libGetErrorMessage()
		return nil, fmt.Errorf("book creation error: %v", msg)
	}
	b := Book{
		ptr: ptr,
	}
	return &b, nil
}

func (b *Book) Close() error {
	libInit()

	if b.ptr == nil {
		return ErrBookIsClosed
	}
	if b.fetcherClosure != nil {
		ffi.ClosureFree(b.fetcherClosure)
		b.fetcherClosure = nil
		b.fetcher = nil
	}
	libDestroy(b.ptr)
	b.ptr = nil
	return nil
}

func (b *Book) GetPageSize() PageSize {
	libInit()
	if b.ptr == nil {
		return PageSize{}
	}
	return libGetPageSize(b.ptr)
}

func (b *Book) GetPageSizeAt(index int) PageSize {
	libInit()
	if b.ptr == nil {
		return PageSize{}
	}
	return libGetPageSizeAt(b.ptr, index)
}

func (b *Book) GetMediaType() MediaType {
	libInit()
	if b.ptr == nil {
		return MediaTypePrint
	}
	ret := libGetMediaType(b.ptr)
	return MediaType(ret)
}

func (b *Book) GetPageMargins() PageMargins {
	libInit()
	if b.ptr == nil {
		return PageMargins{}
	}
	ret := libGetPageMargins(b.ptr)
	return ret
}

func (b *Book) LoadURL(url string, userStyle string, userScript string) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}

	cUrl := strings.CString(url)
	cUserStyle := strings.CString(userStyle)
	cUserScript := strings.CString(userScript)
	ok := libLoadUrl(b.ptr, stringPtr(cUrl), stringPtr(cUserStyle), stringPtr(cUserScript))
	if !ok {
		msg := libGetErrorMessage()
		return fmt.Errorf("error loading url %v: %v", url, msg)
	}
	return nil
}

func (b *Book) LoadHTML(html string, userStyle string, userScript string, baseUrl string) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}

	cHtml, lgHtml := strings.CStringL(html)
	cUserStyle := strings.CString(userStyle)
	cUserScript := strings.CString(userScript)
	cBaseUrl := strings.CString(baseUrl)
	ok := libLoadHtml(b.ptr, binPtr(cHtml), int32(lgHtml), stringPtr(cUserStyle), stringPtr(cUserScript), stringPtr(cBaseUrl))
	if !ok {
		msg := libGetErrorMessage()
		return fmt.Errorf("error loading html: %v", msg)
	}
	return nil
}

func (b *Book) RenderPage(canvas *ImageCanvas, pageIndex int) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}
	if canvas.ptr == nil {
		return ErrCanvasIsClosed
	}
	libRenderPage(b.ptr, canvas.ptr, uint32(pageIndex))
	return nil
}

func (b *Book) WriteToPDF(file string) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}
	ret := libWriteToPDF(b.ptr, file)
	if !ret {
		err := errors.New(libGetErrorMessage())
		return err
	}
	return nil
}

func (b *Book) WriteToPDFRange(file string, pageStart, pageEnd, pageStep int) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}
	ret := libWriteToPDFRange(b.ptr, file, uint32(pageStart), uint32(pageEnd), int32(pageStep))
	if !ret {
		err := errors.New(libGetErrorMessage())
		return err
	}
	return nil
}

func (b *Book) WriteToPDFStream(output io.Writer) error {
	libInit()
	if b.ptr == nil {
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

	isOk := libWriteToPDFStream(b.ptr, callback, unsafe.Pointer(stream))

	if !isOk {
		if stream.err != nil {
			return stream.err
		}
		msg := libGetErrorMessage()
		return fmt.Errorf("pdf write failed: %v", msg)
	}
	return nil
}

func (b *Book) WriteToPDFStreamRange(output io.Writer, pageStart, pageEnd, pageStep int) error {
	libInit()
	if b.ptr == nil {
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

	isOk := libWriteToPDFStreamRange(b.ptr, callback, unsafe.Pointer(stream), uint32(pageStart), uint32(pageEnd), int32(pageStep))
	if !isOk {
		if stream.err != nil {
			return stream.err
		}
		msg := libGetErrorMessage()
		return fmt.Errorf("pdf range write failed: %v", msg)
	}
	return nil
}

func (b *Book) SetCustomResourceFetcher(fetcher CustomResourceFetcher) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}
	if fetcher == nil {
		if b.fetcherClosure != nil {
			ffi.ClosureFree(b.fetcherClosure)
			b.fetcherClosure = nil
			b.fetcher = nil
		}

		libSetCustomResourceFetcher(b.ptr, nil, nil)
		return nil
	}

	// allocate the closure function
	var callback unsafe.Pointer
	b.fetcherClosure = ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if b.fetcherClosure == nil {
		return errors.New("closure not allocated")
	}

	// describe the closure's signature
	// plutobook_resource_data_t* (*plutobook_resource_fetch_callback_t)(void* closure, const char* url);
	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		return fmt.Errorf("cif preparation failed: %v", status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(customResourceFetcherCallback)

	// prepare the closure
	if status := ffi.PrepClosureLoc(b.fetcherClosure, &cifCallback, fn, nil, callback); status != ffi.OK {
		return fmt.Errorf("closure preparation failed: %v", status)
	}
	b.fetcher = fetcher

	// pbClosure := unsafe.Pointer(&fetcher)
	pbClosure := unsafe.Pointer(b)
	libSetCustomResourceFetcher(b.ptr, callback, pbClosure)
	return nil
}

func (b *Book) GetDocumentWidth() float32 {
	libInit()
	if b.ptr == nil {
		return 0
	}
	return libGetDocumentWidth(b.ptr)
}

func (b *Book) GetDocumentHeight() float32 {
	libInit()
	if b.ptr == nil {
		return 0
	}
	return libGetDocumentHeight(b.ptr)
}

func (b *Book) WriteToPNG(file string, width, height int) error {
	libInit()
	if b.ptr == nil {
		return ErrBookIsClosed
	}

	cFile := strings.CString(file)
	ret := libWriteToPNG(b.ptr, stringPtr(cFile), int32(width), int32(height))
	if !ret {
		msg := libGetErrorMessage()
		return fmt.Errorf("file write error: %v", msg)
	}
	return nil
}

func (b *Book) WriteToPNGStream(output io.Writer, width, height int) error {
	libInit()
	if b.ptr == nil {
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

	isOk := libWriteToPNGStream(b.ptr, callback, unsafe.Pointer(stream), int32(width), int32(height))

	if !isOk {
		if stream.err != nil {
			return stream.err
		}
		msg := libGetErrorMessage()
		return fmt.Errorf("png write failed: %v", msg)
	}
	return nil
}
