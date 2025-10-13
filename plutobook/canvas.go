package plutobook

import "unsafe"

type canvasPtr unsafe.Pointer

type canvasInterface interface {
	canvasPtr() canvasPtr

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
	ptr canvasPtr
}

func (c *CanvasBase) canvasPtr() canvasPtr {
	return c.ptr
}

func (c *CanvasBase) Close() error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasDestroy(c.ptr)
	c.ptr = nil
	return nil
}

func (c *CanvasBase) Flush() error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasFlush(c.ptr)
	return nil
}

func (c *CanvasBase) Finish() error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasFinish(c.ptr)
	return nil
}

func (c *CanvasBase) Translate(tx, ty float32) error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasTranslate(c.ptr, tx, ty)
	return nil
}

func (c *CanvasBase) Scale(sx, sy float32) error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasScale(c.ptr, sx, sy)
	return nil
}

func (c *CanvasBase) Rotate(angle float32) error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasRotate(c.ptr, angle)
	return nil
}

func (k *CanvasBase) Transform(a, b, c, d, e, f float32) error {
	if k.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasTransform(k.ptr, a, b, c, d, e, f)
	return nil
}

func (k *CanvasBase) SetMatrix(a, b, c, d, e, f float32) error {
	if k.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasSetMatrix(k.ptr, a, b, c, d, e, f)
	return nil
}

func (c *CanvasBase) ResetMatrix() error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasResetMatrix(c.ptr)
	return nil
}

func (c *CanvasBase) ClipRect(x, y, width, height float32) error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasClipRect(c.ptr, x, y, width, height)
	return nil
}

func (c *CanvasBase) ClearSurface(r, g, b, a float32) error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasClearSurface(c.ptr, r, g, b, a)
	return nil
}

func (c *CanvasBase) SaveState() error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasSaveState(c.ptr)
	return nil
}

func (c *CanvasBase) RestoreState() error {
	if c.ptr == nil {
		return ErrCanvasIsClosed
	}
	libCanvasRestoreState(c.ptr)
	return nil
}
