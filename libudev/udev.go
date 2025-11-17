package libudev

import (
	"sync"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type UDev unsafe.Pointer

var (
	infoTrack map[UDev]track = make(map[UDev]track)
	infoLck   sync.Mutex
)

type track struct {
	refCount   int
	logClosure *ffi.Closure
}

func New() UDev {
	initLib()

	infoLck.Lock()
	defer infoLck.Unlock()

	ptr := libudevNew()
	infoTrack[ptr] = track{
		refCount: 1,
	}
	return ptr
}

func Ref(udev UDev) {
	initLib()

	infoLck.Lock()
	defer infoLck.Unlock()

	track, ok := infoTrack[udev]
	libudevRef(udev)
	if ok {
		track.refCount++
		infoTrack[udev] = track
	}
}

func Unref(udev UDev) {
	initLib()

	infoLck.Lock()
	defer infoLck.Unlock()

	track, ok := infoTrack[udev]
	libudevUnref(udev)
	if ok {
		track.refCount--
		if track.refCount <= 0 {
			if track.logClosure != nil {
				ffi.ClosureFree(track.logClosure)
			}
			delete(infoTrack, udev)
		}
	}
}
