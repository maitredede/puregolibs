package libevdi

import "unsafe"

type Handle unsafe.Pointer

func Open(device int) Handle {
	libInit()
	return evdiOpen(device)
}

func OpenAttachedTo(sysfsParentDevice string) Handle {
	libInit()
	panic("TODO")
}

func OpenAttachedToNone() Handle {
	libInit()
	return evdiOpenAttachedToFixed(nil, 0)
}

func Close(handle Handle) {
	libInit()
	evdiClose(handle)
}
