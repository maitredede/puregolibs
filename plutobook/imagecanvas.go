package plutobook

import (
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
)

type ImageCanvas struct {
	CanvasBase
}

var _ canvasInterface = (*ImageCanvas)(nil)

func CreateImageCanvas(width, height int, imageFormat ImageFormat) (*ImageCanvas, error) {
	libInit()
	ptr := libImageCanvasCreate(int32(width), int32(height), imageFormat)
	if ptr == 0 {
		msg := libGetErrorMessage()
		return nil, fmt.Errorf("image canvas create failed: %v", msg)
	}
	c := &ImageCanvas{
		CanvasBase: CanvasBase{
			ptr: ptr,
		},
	}
	return c, nil
}

func (c *ImageCanvas) GetFormat() (ImageFormat, error) {
	libInit()
	if c.ptr == 0 {
		return ImageFormat(0), ErrCanvasIsClosed
	}
	f := libImageCanvasGetFormat(c.ptr)
	return f, nil
}

func (c *ImageCanvas) GetWidth() int {
	libInit()
	if c.ptr == 0 {
		return 0
	}
	return int(libImageCanvasGetWidth(c.ptr))
}

func (c *ImageCanvas) GetHeight() int {
	libInit()
	if c.ptr == 0 {
		return 0
	}
	return int(libImageCanvasGetHeight(c.ptr))
}

func (c *ImageCanvas) GetStride() int {
	libInit()
	if c.ptr == 0 {
		return 0
	}
	return int(libImageCanvasGetStride(c.ptr))
}

func (c *ImageCanvas) WriteToPNG(file string) error {
	libInit()
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	cFile := strings.CString(file)
	ret := libImageCanvasWriteToPNG(c.ptr, uintptr(unsafe.Pointer(cFile)))
	if !ret {
		msg := libGetErrorMessage()
		return fmt.Errorf("image canvas png write failed: %v", msg)
	}
	return nil
}

func (c *ImageCanvas) WriteToPNGStream(output io.Writer) error {
	libInit()
	if c.ptr == 0 {
		return ErrCanvasIsClosed
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
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); status != ffi.OK {
		return fmt.Errorf("cif preparation failed: %v", status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(streamWriteCallback)

	stream := &streamWriteData{
		output: output,
		err:    nil,
	}
	userData := unsafe.Pointer(stream)
	// prepare the closure
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		return fmt.Errorf("closure preparation failed: %v", status)
	}

	isOk := libImageCanvasWriteToPNGStream(c.ptr, uintptr(callback), 0)

	if !isOk {
		if stream.err != nil {
			return stream.err
		}
		msg := libGetErrorMessage()
		return fmt.Errorf("image canvas write failed: %v", msg)
	}
	return nil
}
