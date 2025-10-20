package plutobook

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
)

type PDFCanvas struct {
	CanvasBase
}

var _ canvasInterface = (*PDFCanvas)(nil)

func CreatePDFCanvas(filename string, size PageSize) (*PDFCanvas, error) {
	libInit()

	sym := mustGetSymbol("plutobook_pdf_canvas_create")

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffi.TypePointer, &ffi.TypePointer, &ffiPageSizeType); ok != ffi.OK {
		panic("plutobook_pdf_canvas_create cif prep is not OK")
	}

	cFilename := strings.CString(filename)
	var ptr canvasPtr
	ffi.Call(&cif, sym, unsafe.Pointer(&ptr), unsafe.Pointer(cFilename), unsafe.Pointer(&size))

	if ptr == nil {
		msg := libGetErrorMessage()
		return nil, fmt.Errorf("pdf canvas create failed: %v", msg)
	}
	c := &PDFCanvas{
		CanvasBase: CanvasBase{
			ptr: ptr,
		},
	}
	return c, nil
}

func (c *PDFCanvas) SetMetadata(metadata PdfMetadata, value string) error {
	libInit()
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libPDFCanvasSetMetadata(c.ptr, metadata, value)
	return nil
}

func (c *PDFCanvas) SetSize(size PageSize) error {
	libInit()
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	// libPDFCanvasSetSize(c.ptr, size)
	// return nil
	sym := mustGetSymbol("plutobook_pdf_canvas_set_size")

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypePointer, &ffiPageSizeType); ok != ffi.OK {
		panic("plutobook_pdf_canvas_set_size cif prep is not OK")
	}
	ffi.Call(&cif, sym, nil, unsafe.Pointer(&c.ptr), unsafe.Pointer(&size))
	return nil
}

func (c *PDFCanvas) ShowPage() error {
	libInit()
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libPDFCanvasShowPage(c.ptr)
	return nil
}
