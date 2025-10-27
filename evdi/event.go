package evdi

import "unsafe"

type evdiEventContext struct {
	dpmsHandler        unsafe.Pointer
	modeChangedHandler unsafe.Pointer
	updateReadyHandler unsafe.Pointer
	crtcStateHandler   unsafe.Pointer
	cursorSetHandler   unsafe.Pointer
	cursorMoveHandler  unsafe.Pointer
	ddcciDataHandler   unsafe.Pointer

	userData unsafe.Pointer
}

type Events struct {
	Dpms        func(dpmsMode int32, userData unsafe.Pointer)
	ModeChanged func(mode EvdiMode, userdata unsafe.Pointer)
	UpdateReady func(bufferToUpdate int32, userData unsafe.Pointer)
	CrtcState   func(state int32, userData unsafe.Pointer)
	CursorSet   func(cursorSet CursorSet, userData unsafe.Pointer)
	CursorMove  func(cursorMove CursorMove, userData unsafe.Pointer)
	DdcCiData   func(data DdcciData, userData unsafe.Pointer)
}

type EvdiMode struct {
	Width        int32
	Height       int32
	RefreshRate  int32
	BitsPerPixel int32
	PixelFormat  uint32
}
