package plutobook

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
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

	purego.RegisterLibFunc(&libWriteToPNGStream, initPtr, "plutobook_write_to_png_stream")
	purego.RegisterLibFunc(&libWriteToPNG, initPtr, "plutobook_write_to_png")

	purego.RegisterLibFunc(&libSetSSLCaInfo, initPtr, "plutobook_set_ssl_cainfo")
	purego.RegisterLibFunc(&libSetSSLCaPath, initPtr, "plutobook_set_ssl_capath")
	purego.RegisterLibFunc(&libSetSSLVerifyPeer, initPtr, "plutobook_set_ssl_verify_peer")
	purego.RegisterLibFunc(&libSetSSLVerifyHost, initPtr, "plutobook_set_ssl_verify_host")
	purego.RegisterLibFunc(&libSetHttpFollowRedirects, initPtr, "plutobook_set_http_follow_redirects")
	purego.RegisterLibFunc(&libSetHttpMaxRedirects, initPtr, "plutobook_set_http_max_redirects")
	purego.RegisterLibFunc(&libSetHttpTimeout, initPtr, "plutobook_set_http_timeout")
	purego.RegisterLibFunc(&libSetCustomResourceFetcher, initPtr, "plutobook_set_custom_resource_fetcher")
	purego.RegisterLibFunc(&libResourceDataCreate, initPtr, "plutobook_resource_data_create")
}

var (
	libSetSSLCaInfo           func(path string)
	libSetSSLCaPath           func(path string)
	libSetSSLVerifyPeer       func(verify bool)
	libSetSSLVerifyHost       func(verify bool)
	libSetHttpFollowRedirects func(follow bool)
	libSetHttpMaxRedirects    func(amount int32)
	libSetHttpTimeout         func(amount int32)

	libGetDocumentWidth         func(book uintptr) float32
	libGetDocumentHeight        func(book uintptr) float32
	libLoadHtml                 func(book uintptr, data uintptr, length int32, userStyle uintptr, userScript uintptr, baseUrl uintptr) bool
	libLoadUrl                  func(book uintptr, url uintptr, userStyle uintptr, userScript uintptr) bool
	libWriteToPNG               func(book uintptr, file uintptr, width, height int32) bool
	libWriteToPNGStream         func(book uintptr, callback uintptr, closure uintptr, width, height int32) bool
	libSetCustomResourceFetcher func(book uintptr, callback uintptr, closure uintptr)

	libCanvasDestroy      func(canvas uintptr)
	libCanvasFlush        func(canvas uintptr)
	libCanvasFinish       func(canvas uintptr)
	libCanvasTranslate    func(canvas uintptr, tx, ty float32)
	libCanvasScale        func(canvas uintptr, sx, sy float32)
	libCanvasRotate       func(canvas uintptr, angle float32)
	libCanvasTransform    func(canvas uintptr, a, b, c, d, e, f float32)
	libCanvasSetMatrix    func(canvas uintptr, a, b, c, d, e, f float32)
	libCanvasResetMatrix  func(canvas uintptr)
	libCanvasClipRect     func(canvas uintptr, x, y, width, height float32)
	libCanvasClearSurface func(canvas uintptr, r, g, b, a float32)
	libCanvasSaveState    func(canvas uintptr)
	libCanvasRestoreState func(canvas uintptr)

	libImageCanvasCreate           func(width, height int32, format ImageFormat) uintptr
	libImageCanvasGetFormat        func(canvas uintptr) ImageFormat
	libImageCanvasGetWidth         func(canvas uintptr) int32
	libImageCanvasGetHeight        func(canvas uintptr) int32
	libImageCanvasGetStride        func(canvas uintptr) int32
	libImageCanvasWriteToPNG       func(canvas uintptr, file uintptr) bool
	libImageCanvasWriteToPNGStream func(canvas uintptr, callback uintptr, closure uintptr) bool

	// libPDFCanvasCreate      func(filename string, size PageSize) uintptr
	libPDFCanvasSetMetadata func(canvas uintptr, metadata PdfMetadata, value string)
	// libPDFCanvasSetSize     func(canvas uintptr, size PageSize)
	libPDFCanvasShowPage func(canvas uintptr)

	libResourceDataCreate func(content uintptr, length uint32, mimeType uintptr, textEncoding uintptr) uintptr
)
