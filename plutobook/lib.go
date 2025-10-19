package plutobook

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libplutobook.dylib"
	case "linux":
		return "libplutobook.so"
	case "windows":
		return "libplutobook-0.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func libInitFuncs() {
	purego.RegisterLibFunc(&libVersion, initPtr, "plutobook_version")
	purego.RegisterLibFunc(&libVersionString, initPtr, "plutobook_version_string")
	purego.RegisterLibFunc(&libBuildInfo, initPtr, "plutobook_build_info")

	purego.RegisterLibFunc(&libGetErrorMessage, initPtr, "plutobook_get_error_message")
	purego.RegisterLibFunc(&libClearErrorMessage, initPtr, "plutobook_clear_error_message")

	// libCreateSym = mustGetSymbol("plutobook_create")
	// purego.RegisterLibFunc(&libCreate, initPtr, "plutobook_create")
	registerFFICreate()

	purego.RegisterLibFunc(&libDestroy, initPtr, "plutobook_destroy")

	//purego.RegisterLibFunc(&libGetPageSize, initPtr, "plutobook_get_page_size")
	registerFFIGetPageSize()
	registerFFIGetPageSizeAt()
	registerFFIGetPageMargins()
	purego.RegisterLibFunc(&libGetMediaType, initPtr, "plutobook_get_media_type")
	purego.RegisterLibFunc(&libGetDocumentWidth, initPtr, "plutobook_get_document_width")
	purego.RegisterLibFunc(&libGetDocumentHeight, initPtr, "plutobook_get_document_height")
	purego.RegisterLibFunc(&libLoadHtml, initPtr, "plutobook_load_html")
	purego.RegisterLibFunc(&libLoadUrl, initPtr, "plutobook_load_url")
	purego.RegisterLibFunc(&libRenderPage, initPtr, "plutobook_render_page")
	purego.RegisterLibFunc(&libWriteToPDF, initPtr, "plutobook_write_to_pdf")
	purego.RegisterLibFunc(&libWriteToPDFRange, initPtr, "plutobook_write_to_pdf_range")
	purego.RegisterLibFunc(&libWriteToPDFStream, initPtr, "plutobook_write_to_pdf_stream")
	purego.RegisterLibFunc(&libWriteToPDFStreamRange, initPtr, "plutobook_write_to_pdf_stream_range")
	purego.RegisterLibFunc(&libWriteToPNG, initPtr, "plutobook_write_to_png")
	purego.RegisterLibFunc(&libWriteToPNGStream, initPtr, "plutobook_write_to_png_stream")

	purego.RegisterLibFunc(&libImageCanvasCreate, initPtr, "plutobook_image_canvas_create")
	// purego.RegisterLibFunc(&libPDFCanvasCreate, initPtr, "plutobook_pdf_canvas_create")
	purego.RegisterLibFunc(&libPDFCanvasSetMetadata, initPtr, "plutobook_pdf_canvas_set_metadata")
	// purego.RegisterLibFunc(&libPDFCanvasSetSize, initPtr, "plutobook_pdf_canvas_set_size")
	purego.RegisterLibFunc(&libPDFCanvasShowPage, initPtr, "plutobook_pdf_canvas_show_page")
	purego.RegisterLibFunc(&libCanvasDestroy, initPtr, "plutobook_canvas_destroy")
	purego.RegisterLibFunc(&libCanvasFlush, initPtr, "plutobook_canvas_flush")
	purego.RegisterLibFunc(&libCanvasFinish, initPtr, "plutobook_canvas_finish")
	purego.RegisterLibFunc(&libCanvasTranslate, initPtr, "plutobook_canvas_translate")
	purego.RegisterLibFunc(&libCanvasScale, initPtr, "plutobook_canvas_scale")
	purego.RegisterLibFunc(&libCanvasRotate, initPtr, "plutobook_canvas_rotate")
	purego.RegisterLibFunc(&libCanvasTransform, initPtr, "plutobook_canvas_transform")
	purego.RegisterLibFunc(&libCanvasSetMatrix, initPtr, "plutobook_canvas_set_matrix")
	purego.RegisterLibFunc(&libCanvasResetMatrix, initPtr, "plutobook_canvas_reset_matrix")
	purego.RegisterLibFunc(&libCanvasClipRect, initPtr, "plutobook_canvas_clip_rect")
	purego.RegisterLibFunc(&libCanvasClearSurface, initPtr, "plutobook_canvas_clear_surface")
	purego.RegisterLibFunc(&libCanvasSaveState, initPtr, "plutobook_canvas_save_state")
	purego.RegisterLibFunc(&libCanvasRestoreState, initPtr, "plutobook_canvas_restore_state")
	purego.RegisterLibFunc(&libImageCanvasGetFormat, initPtr, "plutobook_image_canvas_get_format")
	purego.RegisterLibFunc(&libImageCanvasGetWidth, initPtr, "plutobook_image_canvas_get_width")
	purego.RegisterLibFunc(&libImageCanvasGetHeight, initPtr, "plutobook_image_canvas_get_height")
	purego.RegisterLibFunc(&libImageCanvasGetStride, initPtr, "plutobook_image_canvas_get_stride")
	purego.RegisterLibFunc(&libImageCanvasWriteToPNG, initPtr, "plutobook_image_canvas_write_to_png")
	purego.RegisterLibFunc(&libImageCanvasWriteToPNGStream, initPtr, "plutobook_image_canvas_write_to_png_stream")

	purego.RegisterLibFunc(&libSetSSLCaInfo, initPtr, "plutobook_set_ssl_cainfo")
	purego.RegisterLibFunc(&libSetSSLCaPath, initPtr, "plutobook_set_ssl_capath")
	purego.RegisterLibFunc(&libSetSSLVerifyPeer, initPtr, "plutobook_set_ssl_verify_peer")
	purego.RegisterLibFunc(&libSetSSLVerifyHost, initPtr, "plutobook_set_ssl_verify_host")
	purego.RegisterLibFunc(&libSetHttpFollowRedirects, initPtr, "plutobook_set_http_follow_redirects")
	purego.RegisterLibFunc(&libSetHttpMaxRedirects, initPtr, "plutobook_set_http_max_redirects")
	purego.RegisterLibFunc(&libSetHttpTimeout, initPtr, "plutobook_set_http_timeout")
	purego.RegisterLibFunc(&libSetCustomResourceFetcher, initPtr, "plutobook_set_custom_resource_fetcher")
	purego.RegisterLibFunc(&libResourceDataCreate, initPtr, "plutobook_resource_data_create")
	purego.RegisterLibFunc(&libResourceDataDestroy, initPtr, "plutobook_resource_data_destroy")
	purego.RegisterLibFunc(&libResourceDataGetReferenceCount, initPtr, "plutobook_resource_data_get_reference_count")
}

type stringPtr unsafe.Pointer
type binPtr unsafe.Pointer

var (
	libSetSSLCaInfo           func(path string)
	libSetSSLCaPath           func(path string)
	libSetSSLVerifyPeer       func(verify bool)
	libSetSSLVerifyHost       func(verify bool)
	libSetHttpFollowRedirects func(follow bool)
	libSetHttpMaxRedirects    func(amount int32)
	libSetHttpTimeout         func(amount int32)

	libCreate                   func(pageSize PageSize, margins PageMargins, mediaType MediaType) bookPtr
	libDestroy                  func(book bookPtr)
	libGetPageSize              func(book bookPtr) PageSize
	libGetPageSizeAt            func(book bookPtr, index int) PageSize
	libGetDocumentWidth         func(book bookPtr) float32
	libGetDocumentHeight        func(book bookPtr) float32
	libLoadHtml                 func(book bookPtr, data binPtr, length int32, userStyle stringPtr, userScript stringPtr, baseUrl stringPtr) bool
	libLoadUrl                  func(book bookPtr, url stringPtr, userStyle stringPtr, userScript stringPtr) bool
	libSetCustomResourceFetcher func(book bookPtr, callback unsafe.Pointer, closure unsafe.Pointer)
	libGetMediaType             func(book bookPtr) int32
	libGetPageMargins           func(book bookPtr) PageMargins
	libRenderPage               func(book bookPtr, canvas canvasPtr, pageInted uint32)
	libWriteToPDF               func(book bookPtr, file string) bool
	libWriteToPDFRange          func(book bookPtr, file string, pageStart, pageEnd uint32, pageStep int32) bool
	libWriteToPDFStream         func(book bookPtr, callback unsafe.Pointer, closure unsafe.Pointer) bool
	libWriteToPDFStreamRange    func(book bookPtr, callback unsafe.Pointer, closure unsafe.Pointer, pageStart, pageEnd uint32, pageStep int32) bool
	libWriteToPNG               func(book bookPtr, file stringPtr, width, height int32) bool
	libWriteToPNGStream         func(book bookPtr, callback unsafe.Pointer, closure unsafe.Pointer, width, height int32) bool

	libCanvasDestroy      func(canvas canvasPtr)
	libCanvasFlush        func(canvas canvasPtr)
	libCanvasFinish       func(canvas canvasPtr)
	libCanvasTranslate    func(canvas canvasPtr, tx, ty float32)
	libCanvasScale        func(canvas canvasPtr, sx, sy float32)
	libCanvasRotate       func(canvas canvasPtr, angle float32)
	libCanvasTransform    func(canvas canvasPtr, a, b, c, d, e, f float32)
	libCanvasSetMatrix    func(canvas canvasPtr, a, b, c, d, e, f float32)
	libCanvasResetMatrix  func(canvas canvasPtr)
	libCanvasClipRect     func(canvas canvasPtr, x, y, width, height float32)
	libCanvasClearSurface func(canvas canvasPtr, r, g, b, a float32)
	libCanvasSaveState    func(canvas canvasPtr)
	libCanvasRestoreState func(canvas canvasPtr)

	libImageCanvasCreate           func(width, height int32, format ImageFormat) canvasPtr
	libImageCanvasGetFormat        func(canvas canvasPtr) ImageFormat
	libImageCanvasGetWidth         func(canvas canvasPtr) int32
	libImageCanvasGetHeight        func(canvas canvasPtr) int32
	libImageCanvasGetStride        func(canvas canvasPtr) int32
	libImageCanvasWriteToPNG       func(canvas canvasPtr, file stringPtr) bool
	libImageCanvasWriteToPNGStream func(canvas canvasPtr, callback unsafe.Pointer, closure unsafe.Pointer) bool

	// libPDFCanvasCreate      func(filename string, size PageSize) canvasPtr
	libPDFCanvasSetMetadata func(canvas canvasPtr, metadata PdfMetadata, value string)
	// libPDFCanvasSetSize     func(canvas canvasPtr, size PageSize)
	libPDFCanvasShowPage func(canvas canvasPtr)

	libResourceDataCreate            func(content binPtr, length uint32, mimeType stringPtr, textEncoding stringPtr) resourceDataPtr
	libResourceDataDestroy           func(resource resourceDataPtr)
	libResourceDataGetReferenceCount func(resource resourceDataPtr) uint32
)

func registerFFIGetPageSize() {
	sym := mustGetSymbol("plutobook_get_page_size")

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffiPageSizeType, &ffi.TypePointer); ok != ffi.OK {
		panic("plutobook_get_page_size cif prep is not OK")
	}

	libGetPageSize = func(book bookPtr) PageSize {
		var ret PageSize
		ffi.Call(&cif, sym, unsafe.Pointer(&ret), unsafe.Pointer(&book))
		return ret
	}
}

func registerFFIGetPageSizeAt() {
	sym := mustGetSymbol("plutobook_get_page_size_at")

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffiPageSizeType, &ffi.TypePointer, &ffi.TypeUint32); ok != ffi.OK {
		panic("plutobook_get_page_size_at cif prep is not OK")
	}

	libGetPageSizeAt = func(book bookPtr, index int) PageSize {
		var ret PageSize
		cIndex := uint32(index)
		ffi.Call(&cif, sym, unsafe.Pointer(&ret), unsafe.Pointer(&book), unsafe.Pointer(&cIndex))
		return ret
	}
}

func registerFFIGetPageMargins() {
	sym := mustGetSymbol("plutobook_get_page_margins")

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffiPageMarginsType, &ffi.TypePointer); ok != ffi.OK {
		panic("plutobook_get_page_margins cif prep is not OK")
	}

	libGetPageMargins = func(book bookPtr) PageMargins {
		var ret PageMargins
		ffi.Call(&cif, sym, unsafe.Pointer(&ret), unsafe.Pointer(&book))
		return ret
	}
}
