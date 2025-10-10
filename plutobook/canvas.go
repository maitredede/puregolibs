package plutobook

type canvasInterface interface {
	canvasPtr() uintptr

	Close() error
	Flush() error
	Finish() error
	Translate(tx, ty float32) error
	Scale(sx, sy float32) error
	Rotate(angle float32) error
	Transform(a, b, c, d, e, f float32) error
	SetMatrix(a, b, c, d, e, f float32) error
	ResetMatrix() error
	ClipRect(x, y, width, height float32) error
	ClearSurface(r, g, b, a float32) error
	SaveState() error
	RestoreState() error
	// GetSurface() (uintptr, error)
	// GetContext() (uintptr, error)
}

type CanvasBase struct {
	ptr uintptr
}

func (c *CanvasBase) canvasPtr() uintptr {
	return c.ptr
}

func (c *CanvasBase) Close() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasDestroy(c.ptr)
	c.ptr = 0
	return nil
}

func (c *CanvasBase) Flush() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasFlush(c.ptr)
	return nil
}

func (c *CanvasBase) Finish() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasFinish(c.ptr)
	return nil
}

func (c *CanvasBase) Translate(tx, ty float32) error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasTranslate(c.ptr, tx, ty)
	return nil
}

func (c *CanvasBase) Scale(sx, sy float32) error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasScale(c.ptr, sx, sy)
	return nil
}

func (c *CanvasBase) Rotate(angle float32) error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasRotate(c.ptr, angle)
	return nil
}

func (k *CanvasBase) Transform(a, b, c, d, e, f float32) error {
	if k.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasTransform(k.ptr, a, b, c, d, e, f)
	return nil
}

func (k *CanvasBase) SetMatrix(a, b, c, d, e, f float32) error {
	if k.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasSetMatrix(k.ptr, a, b, c, d, e, f)
	return nil
}

func (c *CanvasBase) ResetMatrix() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasResetMatrix(c.ptr)
	return nil
}

func (c *CanvasBase) ClipRect(x, y, width, height float32) error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasClipRect(c.ptr, x, y, width, height)
	return nil
}

func (c *CanvasBase) ClearSurface(r, g, b, a float32) error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasClearSurface(c.ptr, r, g, b, a)
	return nil
}

func (c *CanvasBase) SaveState() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasSaveState(c.ptr)
	return nil
}

func (c *CanvasBase) RestoreState() error {
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libCanvasRestoreState(c.ptr)
	return nil
}
