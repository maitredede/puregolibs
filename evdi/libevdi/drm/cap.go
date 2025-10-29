package drm

import (
	"os"
	"unsafe"

	"github.com/maitredede/puregolibs/evdi/libevdi/drm/ioctl"
)

type (
	drmGetCap struct {
		id  uint64
		val uint64
	}
)

const (
	CapDumbBuffer uint64 = iota + 1
	CapVBlankHighCRTC
	CapDumbPreferredDepth
	CapDumbPreferShadow
	CapPrime
	CapTimestampMonotonic
	CapAsyncPageFlip
	CapCursorWidth
	CapCursorHeight
	CapAddFB2Modifiers
	CapPageFlipTarget
	CapCrtcInVblankEvent
	CapSyncObj
	CapSyncObjTimeline
	CapAtomicAsyncPageFlip
)

func HasDumbBuffer(file *os.File) bool {
	cap, err := GetCap(file, CapDumbBuffer)
	if err != nil {
		return false
	}
	return cap != 0
}

func GetCap(file *os.File, capid uint64) (uint64, error) {
	cap := &drmGetCap{}
	cap.id = capid
	err := ioctl.Do(uintptr(file.Fd()), uintptr(IOCTLGetCap), uintptr(unsafe.Pointer(cap)))
	if err != nil {
		return 0, err
	}
	return cap.val, nil
}
