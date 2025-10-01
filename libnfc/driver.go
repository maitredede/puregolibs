package libnfc

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
)

type nativeDriver struct {
	// const char *name
	name uintptr
	// const scan_type_enum scan_type
	scanType ScanMode
	// size_t (*scan)(const nfc_context *context, nfc_connstring connstrings[], const size_t connstrings_len)
	scan uintptr
	// struct nfc_device *(*open)(const nfc_context *context, const nfc_connstring connstring)
	open uintptr
	// void (*close)(struct nfc_device *pnd)
	close uintptr
	// const char *(*strerror)(const struct nfc_device *pnd)
	strerror uintptr

	initiatorInit                 uintptr
	initiatorInitSecureElement    uintptr
	initiatorSelectPassiveTarget  uintptr
	initiatorPollTarget           uintptr
	initiatorSelectDepTarget      uintptr
	initiatorDeselectTarget       uintptr
	initiatorTransceiveBytes      uintptr
	initiatorTransceiveBits       uintptr
	initiatorTransceiveBytesTimed uintptr
	initiatorTransceiveBitsTimed  uintptr
	initiatorTargetIsPresent      uintptr

	targetInit         uintptr
	targetSendBytes    uintptr
	targetReceiveBytes uintptr
	targetSendBits     uintptr
	targetReceiveBits  uintptr

	deviceSetPropertyBool     uintptr
	deviceSetPropertyInt      uintptr
	getSupportedModulation    uintptr
	getSupportedBaudRate      uintptr
	deviceGetInformationAbout uintptr

	abortCommand uintptr
	idle         uintptr
	powerdown    uintptr
}

var (
	// int nfc_register_driver(const nfc_driver *driver)
	libNfcRegisterDriver func(driver uintptr) int16
)

type Driver struct {
	Name string
}

type driverInfo struct {
	driver *Driver

	nameBin     *byte
	nd          *nativeDriver
	disposables []func()
}

func (di *driverInfo) Dispose() {
	for _, d := range di.disposables {
		d()
	}
}

type ScanMode int16

const (
	ScanModeNotIntrusive ScanMode = iota
	ScanModeIntrusive
	ScanModeNotAvailable
)

var (
	drivers = make([]driverInfo, 0)
)

func RegisterDriver(driver *Driver) error {
	libInit()

	di := driverInfo{
		driver: driver,
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	di.disposables = append(di.disposables, func() {
		ffi.ClosureFree(closure)
	})
	// describe the closure's signature
	var cifCallback ffi.Cif
	args := []*ffi.Type{
		&ffi.TypePointer,
		&ffi.TypePointer, //TODO
		&ffi.TypeSint32,
	}
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeSint32, args...); status != ffi.OK {
		panic(status)
	}
	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		fmt.Println("Hello, World Scan NFC Go!")
		return 0
	})
	// prepare the closure
	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
	}

	di.nameBin = strings.CString(driver.Name)

	namePtr := uintptr(unsafe.Pointer(di.nameBin))
	nd := nativeDriver{
		name: namePtr,
		// scanType: ScanModeNotAvailable,
		scanType: ScanModeNotIntrusive,
		scan:     uintptr(callback),
	}
	di.nd = &nd

	driverPtr := uintptr(unsafe.Pointer(&di.nd))

	ret := libNfcRegisterDriver(driverPtr)
	if ret == 0 {
		drivers = append(drivers, di)
		return nil
	}
	defer di.Dispose()

	return fmt.Errorf("register driver ret=%v", ret)
}
