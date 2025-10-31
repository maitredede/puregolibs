package libevdi

import (
	"syscall"
)

type CursorSet struct {
	HotX    int32
	HotY    int32
	Width   uint32
	Height  uint32
	Enabled bool
	// BufferLength uint32
	// Buffer       *uint32
	BufferData  []byte
	PixelFormat uint32
	Stride      uint32
}

type CursorMove struct {
	X int32
	Y int32
}

func toBool(v uint8) bool {
	return v != 0
}

func (h *Handle) toEvdiCursorSet(event drmEvdiEventCursorSet) CursorSet {
	cursorSet := CursorSet{
		HotX:    event.hotX,
		HotY:    event.hotY,
		Width:   event.width,
		Height:  event.height,
		Enabled: toBool(event.enabled),
		// BufferLength: 0,
		// Buffer:       nil,
		BufferData:  nil,
		PixelFormat: event.pixelFormat,
		Stride:      event.stride,
	}

	if event.enabled != 0 {
		size := event.bufferLength
		offset := uint64(0)

		if err := h.getDumbOffset(event.bufferHandle, &offset); err != nil {
			evdiLogInfo("error: DRM_IOCTL_MODE_MAP_DUMB failed with error: %v", err)
			return cursorSet
		}

		bin, err := syscall.Mmap(int(h.fd.Fd()), int64(offset), int(size), syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			evdiLogInfo("error: mmap failed with error: %v", err)
			return cursorSet
		}
		cursorSet.BufferData = bin[:]
		if err := syscall.Munmap(bin); err != nil {
			evdiLogDebug("error: mmap failed with error: %v", err)
		}
	}
	return cursorSet
}

func toEvdiCursorMove(event drmEvdiEventCursorMove) CursorMove {
	c := CursorMove{
		X: event.x,
		Y: event.y,
	}
	return c
}
