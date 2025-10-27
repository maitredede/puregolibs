package evdi

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

var (
	libLck sync.Mutex
	libPtr uintptr
	libErr error
)

func getSymbol(sym string) (uintptr, error) {
	return purego.Dlsym(libPtr, sym)
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func initLib() {
	libLck.Lock()
	defer libLck.Unlock()

	if libErr != nil {
		panic(libErr)
	}
	if libPtr != 0 {
		return
	}

	libPtr, libErr = purego.Dlopen("libevdi.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if libErr != nil {
		panic(libErr)
	}

	purego.RegisterLibFunc(&libEvdiCheckDevice, libPtr, "evdi_check_device")
	purego.RegisterLibFunc(&libEvdiAddDevice, libPtr, "evdi_add_device")
	purego.RegisterLibFunc(&libEvdiOpen, libPtr, "evdi_open")
	purego.RegisterLibFunc(&libEvdiOpenAttachedToFixed, libPtr, "evdi_open_attached_to_fixed")

	purego.RegisterLibFunc(&libEvdiClose, libPtr, "evdi_close")
	purego.RegisterLibFunc(&libEvdiConnect, libPtr, "evdi_connect")
	purego.RegisterLibFunc(&libEvdiConnect2, libPtr, "evdi_connect2")
	purego.RegisterLibFunc(&libEvdiDisconnect, libPtr, "evdi_disconnect")
	purego.RegisterLibFunc(&libEvdiEnableCursorEvents, libPtr, "evdi_enable_cursor_events")

	purego.RegisterLibFunc(&libevdiGrabPixels, libPtr, "evdi_grab_pixels")
	//purego.RegisterLibFunc(&libevdiRegisterBuffer, libPtr, "evdi_register_buffer")
	{
		typeEvdiBuffer := ffi.NewType(
			&ffi.TypeSint32,  // id
			&ffi.TypePointer, // buffr
			&ffi.TypeSint32,  // width
			&ffi.TypeSint32,  // height
			&ffi.TypeSint32,  // stride
			&ffi.TypePointer, // rects
			&ffi.TypeSint32,  // rect_count
		)
		libevdiRegisterBufferSym := mustGetSymbol("evdi_register_buffer")
		argTypes := []*ffi.Type{
			&ffi.TypePointer, // handle
			&typeEvdiBuffer,  // buffer
		}
		var libevdiRegisterBufferCif ffi.Cif
		if status := ffi.PrepCif(&libevdiRegisterBufferCif, ffi.DefaultAbi, uint32(len(argTypes)), &ffi.TypeVoid, argTypes...); status != ffi.OK {
			panic(status)
		}

		libevdiRegisterBuffer = func(handle evdiHandle, buffer *evdiBuffer) {
			args := []unsafe.Pointer{
				unsafe.Pointer(&handle),
				unsafe.Pointer(buffer),
			}
			ffi.Call(&libevdiRegisterBufferCif, libevdiRegisterBufferSym, nil, args...)
		}
	}
	purego.RegisterLibFunc(&libevdiUnregisterBuffer, libPtr, "evdi_unregister_buffer")
	purego.RegisterLibFunc(&libevdiRequestUpdate, libPtr, "evdi_request_update")
	purego.RegisterLibFunc(&libevdiDdcciResponse, libPtr, "evdi_ddcci_response")

	purego.RegisterLibFunc(&libEvdiHandleEvents, libPtr, "evdi_handle_events")
	purego.RegisterLibFunc(&libEvdiGetEventReady, libPtr, "evdi_get_event_ready")
	purego.RegisterLibFunc(&libEvdiGetLibVersion, libPtr, "evdi_get_lib_version")
	// purego.RegisterLibFunc(&libEvdiSetLogging, libPtr, "evdi_set_logging")

	purego.RegisterLibFunc(&libXorgRunning, libPtr, "Xorg_running")
}

var (
	libEvdiCheckDevice         func(device int32) DeviceStatus
	libEvdiAddDevice           func() int32
	libEvdiOpen                func(device int32) evdiHandle
	libEvdiOpenAttachedToFixed func(sysfsParentDevice unsafe.Pointer, length uint32) evdiHandle

	libEvdiClose              func(handle evdiHandle)
	libEvdiConnect            func(handle evdiHandle, edid *byte, edidLength uint32, skuAreaLimit uint32)
	libEvdiConnect2           func(handle evdiHandle, edid *byte, edidLength uint32, pixelArealimit uint32, pixelPerSecondLimit uint32)
	libEvdiDisconnect         func(handle evdiHandle)
	libEvdiEnableCursorEvents func(handle evdiHandle, enable bool)

	libevdiGrabPixels       func(handle evdiHandle, rects **evdiRect, numRects *int32)
	libevdiRegisterBuffer   func(handle evdiHandle, buffer *evdiBuffer)
	libevdiUnregisterBuffer func(handle evdiHandle, bufferId int32)
	libevdiRequestUpdate    func(handle evdiHandle, bufferId int32) bool
	libevdiDdcciResponse    func(handle evdiHandle, buffer unsafe.Pointer, bufferLength uint32, result bool)

	libEvdiHandleEvents  func(handle evdiHandle, evtctx *evdiEventContext)
	libEvdiGetEventReady func(handle evdiHandle) evdiSelectable
	libEvdiGetLibVersion func(version unsafe.Pointer)
	// libEvdiSetLogging       func(logging) bool

	libXorgRunning func() bool
)

type evdiSelectable int32
