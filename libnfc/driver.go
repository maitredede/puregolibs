package libnfc

import (
	"fmt"
	"log/slog"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type nfcDriverPtr unsafe.Pointer

type nativeDriverList struct {
	next   *nativeDriverList
	driver *NativeDriver
}

type NativeDriver struct {
	// const char *name
	name unsafe.Pointer
	// const scan_type_enum scan_type
	scanType ScanMode
	// size_t (*scan)(const nfc_context *context, nfc_connstring connstrings[], const size_t connstrings_len)
	scan unsafe.Pointer
	// struct nfc_device *(*open)(const nfc_context *context, const nfc_connstring connstring)
	open unsafe.Pointer
	// void (*close)(struct nfc_device *pnd)
	close unsafe.Pointer
	// const char *(*strerror)(const struct nfc_device *pnd)
	strerror unsafe.Pointer

	initiatorInit                 unsafe.Pointer
	initiatorInitSecureElement    unsafe.Pointer
	initiatorSelectPassiveTarget  unsafe.Pointer
	initiatorPollTarget           unsafe.Pointer
	initiatorSelectDepTarget      unsafe.Pointer
	initiatorDeselectTarget       unsafe.Pointer
	initiatorTransceiveBytes      unsafe.Pointer
	initiatorTransceiveBits       unsafe.Pointer
	initiatorTransceiveBytesTimed unsafe.Pointer
	initiatorTransceiveBitsTimed  unsafe.Pointer
	initiatorTargetIsPresent      unsafe.Pointer

	targetInit         unsafe.Pointer
	targetSendBytes    unsafe.Pointer
	targetReceiveBytes unsafe.Pointer
	targetSendBits     unsafe.Pointer
	targetReceiveBits  unsafe.Pointer

	deviceSetPropertyBool     unsafe.Pointer
	deviceSetPropertyInt      unsafe.Pointer
	getSupportedModulation    unsafe.Pointer
	getSupportedBaudRate      unsafe.Pointer
	deviceGetInformationAbout unsafe.Pointer

	abortCommand unsafe.Pointer
	idle         unsafe.Pointer
	powerdown    unsafe.Pointer
}

type Driver struct {
	Name     string
	ScanMode ScanMode
	Scan     func() ([]string, LibNfcError)
}

type driverInfo struct {
	driver *Driver

	nameBin     *byte
	nd          *NativeDriver
	ndPtr       nfcDriverPtr
	disposables []func()
}

func (di *driverInfo) Dispose() {
	for _, d := range di.disposables {
		d()
	}
}

var (
	goDrivers = make([]*driverInfo, 0)
)

func GetNativeDrivers() ([]*NativeDriver, error) {
	sym, err := getSymbol("nfc_drivers")
	if err != nil {
		return nil, err
	}
	ptr := unsafe.Pointer(sym)
	nfcDrivers := (*nativeDriverList)(ptr)
	cur := nfcDrivers
	res := make([]*NativeDriver, 0, 10)
	for {
		if cur == nil {
			break
		}
		if cur.driver == nil {
			cur = cur.next
			continue
		}

		// name := strings.GoStringN(uintptr(cur.driver.name), 256)
		// t.Logf("driver: %s", name)

		res = append(res, cur.driver)
		cur = cur.next
	}

	return res, nil
}

func RegisterDriver(driver *Driver) error {
	libInit()

	di := &driverInfo{
		driver: driver,
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}
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
		slog.Debug("Hello, World Scan NFC Go!")

		argArr := unsafe.Slice(args, cif.NArgs)
		ctx := *(*nfcContextPtr)(argArr[0])
		conStringsPtr := *(**nfcConnString)(unsafe.Pointer(argArr[1]))
		connstringsLen := *(*int32)(argArr[2])
		connStrings := (unsafe.Slice(conStringsPtr, connstringsLen))
		di := (*driverInfo)(userData)

		for i := 0; i < int(connstringsLen); i++ {
			connStrings[i].Set("")
		}

		result, err := di.driver.Scan()
		if err != LibNfcSuccess {
			*(*int32)(ret) = int32(err)
			return 0
		}

		maxNum := min(len(result), int(connstringsLen))

		for i := 0; i < maxNum; i++ {
			cs := result[i]
			// bin := append([]byte(cs), 0)

			// ptr := (*byte)(unsafe.Pointer(&connStrings[i][0]))
			// dstSlice := unsafe.Slice(ptr, len(bin))
			// copy(dstSlice, bin)
			connStrings[i].Set(cs)
		}

		_ = ctx

		*(*int32)(ret) = int32(maxNum)
		return 0
	})
	// prepare the closure
	closureUserData := unsafe.Pointer(di)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, closureUserData, callback); status != ffi.OK {
		panic(status)
	}

	cName := append([]byte(driver.Name), 0)
	if len(cName) > 255 {
		cName = append(cName[:255], 0)
	}
	di.nameBin = &cName[0]
	// di.nameBin = strings.CString(driver.Name)
	// namePtr := unsafe.Pointer(di.nameBin)
	nd := NativeDriver{
		name: unsafe.Pointer(di.nameBin),
		// scanType: ScanModeNotAvailable,
		// scanType: ScanModeNotIntrusive,
		scanType: driver.ScanMode,
		scan:     callback,
	}
	di.nd = &nd
	di.ndPtr = nfcDriverPtr(di.nd)

	ret := libNfcRegisterDriver(di.ndPtr)
	if ret == 0 {
		goDrivers = append(goDrivers, di)
		return nil
	}
	// defer di.Dispose()

	return fmt.Errorf("register driver ret=%v", ret)
}
