package plutobook

type Canvas struct {
	ptr uintptr
}

func (c *Canvas) Close() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	// libDestroy(uintptr(c.ptr))
	// c.ptr = 0
	// return nil
	panic("TODO")
}

func NewImageCanvas(width, height int, imageFormat ImageFormat) (*Canvas, error) {
	libInit()
	panic("TODO")
}

func (c *Canvas) ClearSurface(r, g, b, a float32) error {
	libInit()
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	panic("TODO")
}

func (c *Canvas) WriteToPNG(file string) error {
	libInit()
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	panic("TODO")
}
