package sane

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
)

// Device represents a scanning device.
type Device struct {
	Name, Vendor, Model, Type string
}

type internalSANE_Device struct {
	Name   unsafe.Pointer // char*
	Vendor unsafe.Pointer // char*
	Model  unsafe.Pointer // char*
	Type   unsafe.Pointer // char*
}

// var (
// 	ffiDeviceType = ffi.NewType(
// 		&ffi.TypeUint8,
// 		&ffi.TypeUint8,
// 		&ffi.TypeUint8,
// 		&ffi.TypeUint8,
// 	)
// )

var (
	lckGetDevice sync.Mutex
)

func GetDevices(localOnly bool) ([]Device, error) {

	lckGetDevice.Lock()
	defer lckGetDevice.Unlock()

	sym := mustGetSymbol("sane_get_devices")
	argTypes := []*ffi.Type{
		&ffi.TypePointer,
		&ffi.TypeSint32,
	}
	retType := ffi.TypeSint16

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, uint32(len(argTypes)), &retType, argTypes...); ok != ffi.OK {
		panic("sane_get_devices cif prep is not OK")
	}

	var ret SANE_Status
	var argDeviceList unsafe.Pointer

	argDeviceListPtr := unsafe.Pointer(&argDeviceList)

	var argLocalOnly int32 = 0
	if localOnly {
		argLocalOnly = 1
	}

	args := []unsafe.Pointer{
		unsafe.Pointer(&argDeviceListPtr),
		unsafe.Pointer(&argLocalOnly),
	}

	ffi.Call(&cif, sym, unsafe.Pointer(&ret), args...)

	if ret != StatusGood {
		return nil, fmt.Errorf("device list error: %v", ret)
	}

	if uintptr(argDeviceList) == 0 {
		// no devices found
		return nil, nil
	}
	devicePtrs := (*[1 << 30]*internalSANE_Device)(argDeviceList)
	devices := make([]Device, 0)
	for i := 0; devicePtrs[i] != nil; i++ {
		nd := *devicePtrs[i]
		d := Device{
			Name:   strings.GoString((*byte)(nd.Name)),
			Vendor: strings.GoString((*byte)(nd.Vendor)),
			Model:  strings.GoString((*byte)(nd.Model)),
			Type:   strings.GoString((*byte)(nd.Type)),
		}
		devices = append(devices, d)
	}
	return devices, nil
}
