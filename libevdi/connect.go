package libevdi

import "unsafe"

func Connect(h Handle, edid []byte, skuAreaLimit int) {
	libInit()

	var edidPtr unsafe.Pointer
	edidLen := len(edid)
	if len(edid) > 0 {
		edidPtr = unsafe.Pointer(&edid[0])
	}
	evdiConnect(h, edidPtr, uint(edidLen), uint32(skuAreaLimit))
}

func Connect2(h Handle, edid []byte, pixelAreaLimit, pixelPerSecondLimit int) {
	libInit()

	var edidPtr unsafe.Pointer
	edidLen := len(edid)
	if len(edid) > 0 {
		edidPtr = unsafe.Pointer(&edid[0])
	}
	evdiConnect2(h, edidPtr, uint(edidLen), uint32(pixelAreaLimit), uint32(pixelPerSecondLimit))
}

func Disconnect(h Handle) {
	libInit()

	evdiDisconnect(h)
}
