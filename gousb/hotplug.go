package gousb

import (
	"errors"
	"fmt"
	"log/slog"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type HotplugEvent int32

const (
	HotplugEventDeviceArrived HotplugEvent = 1 << iota
	HotplugEventDeviceLeft
)

type HotplugFlags int32

const (
	HotplugFlagsNoFlags   HotplugFlags = 0
	HotplugFlagsEnumerate HotplugFlags = 1 << 0
)

const HotplugMatchAny int32 = -1

type HotplugCallback func(ctx *Context, device *Device, event HotplugEvent, userData any)

type hotplugCallbackHandle unsafe.Pointer

type HotplugHandle struct {
	ctx *Context
	h   hotplugCallbackHandle
	cb  HotplugCallback

	closure *ffi.Closure
}

func (c *Context) RegisterCallback(events HotplugEvent, flags HotplugFlags, vendorID int32, productID int32, devClass int32, cbfn HotplugCallback /*, userData any*/) (*HotplugHandle, error) {
	v := libusbGetVersion()
	if v.major < 1 || /*(v.major == 1 && v.minor < 0) ||*/ (v.major == 1 && v.minor == 0 && v.micro < 16) {
		return nil, fmt.Errorf("hotplug requires libusb >= 1.0.16")
	}

	// allocate the closure function
	var nativeCbFn unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &nativeCbFn)
	if closure == nil {
		return nil, errors.New("closure alloc failed")
	}
	// describe the closure's signature
	var cifCallback ffi.Cif
	// int hotplug_callback(struct libusb_context *ctx, struct libusb_device *dev, libusb_hotplug_event event, void *user_data)
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 4, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
		return nil, fmt.Errorf("ffi prepare cif status: %v", status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(nativeHotplugCallback)

	// prepare the closure
	userData := unsafe.Pointer(&cbfn)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, nativeCbFn); status != ffi.OK {
		return nil, fmt.Errorf("ffi prepare closur loc: %v", status)
	}

	goH := &HotplugHandle{
		ctx: c,
		// h:       h,
		cb:      cbfn,
		closure: closure,
	}

	var h hotplugCallbackHandle
	ret := libusbHotplugRegisterCallback(c.ptr, events, flags, vendorID, productID, devClass, nativeCbFn, unsafe.Pointer(goH), &h)
	err := errorFromRet(ret)
	if err != nil {
		ffi.ClosureFree(closure)
		return nil, err
	}
	goH.h = h

	return goH, nil
}

func nativeHotplugCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	// int hotplug_callback(struct libusb_context *ctx, struct libusb_device *dev, libusb_hotplug_event event, void *user_data)
	argArr := unsafe.Slice(args, cif.NArgs)
	ctxPtr := *(*libusbContext)(argArr[0])
	devPtr := *(*libusbDevice)(argArr[1])
	event := *(*HotplugEvent)(argArr[2])
	usbUserDataPtr := *(*unsafe.Pointer)(argArr[3])

	cbfn := *(*HotplugCallback)(userData)
	goH := *(**HotplugHandle)(usbUserDataPtr)

	ctx := func(p libusbContext) *Context {
		contextMapLck.Lock()
		defer contextMapLck.Unlock()
		val, ok := contextMap[p]
		if !ok {
			slog.Warn("libusb context not found in global context map")
			return nil
		}
		return val
	}(ctxPtr)

	desc, err := ctx.getDeviceDesc(devPtr)
	if err != nil {
		slog.Warn(fmt.Sprintf("error describing device: %v", err))
	}
	dev := &Device{
		ctx: ctx,
		//handle: devp,
		Desc: desc,
	}

	cbfn(ctx, dev, event, usbUserDataPtr)

	_ = goH

	return 0
}

func (h *HotplugHandle) Deregister() {
	libusbHotplugDeregisterCallback(h.ctx.ptr, h.h)
	if h.closure != nil {
		ffi.ClosureFree(h.closure)
		h.closure = nil
	}
}
