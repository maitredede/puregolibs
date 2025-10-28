package evdi

import (
	"sync/atomic"
	"unsafe"

	"github.com/maitredede/puregolibs/libc"
)

var (
	bufferNumerator atomic.Int32
)

type Buffer struct {
	evdiHandle evdiHandle
	buffer     *evdiBuffer
}

func NewBuffer(mode EvdiMode, evdiHandle evdiHandle) *Buffer {
	id := bufferNumerator.Add(1)

	b := &Buffer{
		evdiHandle: evdiHandle,
		buffer:     &evdiBuffer{},
	}
	stride := mode.Width
	pitchMask := int32(63)

	stride += pitchMask
	stride &= ^pitchMask
	stride *= 4

	b.buffer.id = id
	b.buffer.width = mode.Width
	b.buffer.height = mode.Height
	b.buffer.stride = stride
	b.buffer.rectCount = 16
	b.buffer.rects = (*evdiRect)(libc.CAlloc(b.buffer.rectCount, int32(unsafe.Sizeof(evdiRect{}))))
	bytesPerPixel := mode.BitsPerPixel / 8
	bufferSize := mode.Width * mode.Height * bytesPerPixel
	b.buffer.buffer = (*byte)(libc.CAlloc(1, bufferSize))

	libevdiRegisterBuffer(b.evdiHandle, b.buffer)

	return b
}

func (b *Buffer) Close() {
	libevdiUnregisterBuffer(b.evdiHandle, b.buffer.id)
	libc.Free(unsafe.Pointer(b.buffer.buffer))
	libc.Free(unsafe.Pointer(b.buffer.rects))
}
