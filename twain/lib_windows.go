package twain

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

func libDSMEntry(origin *Identity, dest *Identity, dataGroup DataGroup, dataArgument DataArgumentType, message Message, data uintptr) uint16 {
	var libName string
	switch runtime.GOARCH {
	case "386":
		libName = "twain_32.dll"
	case "amd64":
		libName = "twain_64.dll"
	default:
		panic(fmt.Errorf("GOARCH '%s' not supported", runtime.GOARCH))
	}
	dll, err := syscall.LoadDLL(libName)
	if err != nil {
		panic(err)
	}
	defer dll.Release()
	dsmEntryFunc, err := dll.FindProc("DSM_Entry")
	if err != nil {
		panic(err)
	}

	r1, r2, err := dsmEntryFunc.Call(
		uintptr(unsafe.Pointer(&origin)),
		uintptr(unsafe.Pointer(&dest)),
		uintptr(unsafe.Pointer(&dataGroup)),
		uintptr(unsafe.Pointer(&dataArgument)),
		uintptr(unsafe.Pointer(&message)),
		uintptr(unsafe.Pointer(&data)),
	)
	_ = r1
	_ = r2
	_ = err
	panic("TODO")
}
