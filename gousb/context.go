package gousb

import (
	"unsafe"
)

type Context struct {
	ptr      unsafe.Pointer
	ptrValid bool
}

func Init() (*Context, error) {
	libInit()

	var ptr unsafe.Pointer
	ret := libusbInit(&ptr)
	err := errorFromRet(ret)
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		ptr:      ptr,
		ptrValid: true,
	}
	return ctx, nil
}

func (c *Context) Close() error {
	libInit()

	if !c.ptrValid {
		return ErrInvalidContext
	}
	libusbExit(c.ptr)
	c.ptrValid = false
	return nil
}

func (c *Context) SetDebug(level LogLevel) {
	libInit()

	if !c.ptrValid {
		return
	}
	libusbSetDebug(c.ptr, level)
}

func (c *Context) SetLocale(locale string) error {
	libInit()

	if !c.ptrValid {
		return ErrInvalidContext
	}
	// cLocale := append([]byte(locale), 0)
	// ret := libusbSetLocale(c.ptr, unsafe.Pointer(&cLocale[0]))
	ret := libusbSetLocale(c.ptr, locale)
	err := errorFromRet(ret)
	return err
}
