package libnfc

import (
	"errors"
	"fmt"
	"slices"
	"unsafe"

	"github.com/maitredede/puregolibs/tools"
)

type LibNfcError uintptr

var (
	LibNfcSuccess      LibNfcError = 0
	LibNfcEIO          LibNfcError = LibNfcSuccess - 1 // libNfcError(^uintptr(0))
	LibNfcEInvArg      LibNfcError = LibNfcSuccess - 2
	LibNfcEDEVNOTSUPP  LibNfcError = LibNfcSuccess - 3
	LibNfcENOTSUCHDEV  LibNfcError = LibNfcSuccess - 4
	LibNfcEOVFLOW      LibNfcError = LibNfcSuccess - 5
	LibNfcETIMEOUT     LibNfcError = LibNfcSuccess - 6
	LibNfcEOPABORTED   LibNfcError = LibNfcSuccess - 7
	LibNfcENOTIMPL     LibNfcError = LibNfcSuccess - 8
	LibNfcETGRELEASED  LibNfcError = LibNfcSuccess - 10
	LibNfcERFTRANS     LibNfcError = LibNfcSuccess - 20
	LibNfcEMFCAUTHFAIL LibNfcError = LibNfcSuccess - 30
	LibNfcESOFT        LibNfcError = LibNfcSuccess - 80
	LibNfcECHIP        LibNfcError = LibNfcSuccess - 90
)

func (e LibNfcError) Error() error {
	switch e {
	case LibNfcEIO:
		return errors.New("libNfcEIO")
	case LibNfcEInvArg:
		return errors.New("libNfcEInvArg")
	case LibNfcEDEVNOTSUPP:
		return errors.New("libNfcEDEVNOTSUPP")
	case LibNfcENOTSUCHDEV:
		return errors.New("libNfcENOTSUCHDEV")
	case LibNfcEOVFLOW:
		return errors.New("libNfcEOVFLOW")
	case LibNfcETIMEOUT:
		return errors.New("libNfcETIMEOUT")
	case LibNfcEOPABORTED:
		return errors.New("libNfcEOPABORTED")
	case LibNfcENOTIMPL:
		return errors.New("libNfcENOTIMPL")
	case LibNfcETGRELEASED:
		return errors.New("libNfcETGRELEASED")
	case LibNfcERFTRANS:
		return errors.New("libNfcERFTRANS")
	case LibNfcEMFCAUTHFAIL:
		return errors.New("libNfcEMFCAUTHFAIL")
	case LibNfcESOFT:
		return errors.New("libNfcESOFT")
	case LibNfcECHIP:
		return errors.New("libNfcECHIP")
	}
	return fmt.Errorf("unkown nfc error: %v", int64(e))
}

var _ tools.Timeout = (*LibNfcError)(nil)

func (e LibNfcError) Timeout() bool {
	return e == LibNfcETIMEOUT
}

var errList []LibNfcError = []LibNfcError{
	LibNfcEIO,
	LibNfcEInvArg,
	LibNfcEDEVNOTSUPP,
	LibNfcENOTSUCHDEV,
	LibNfcEOVFLOW,
	LibNfcETIMEOUT,
	LibNfcEOPABORTED,
	LibNfcENOTIMPL,
	LibNfcETGRELEASED,
	LibNfcERFTRANS,
	LibNfcEMFCAUTHFAIL,
	LibNfcESOFT,
	LibNfcECHIP,
}

func isLibErrorInt32(ret int32) bool {
	is := slices.Contains(errList, LibNfcError(ret))
	return is
}

func isLibErrorPtr(ptr unsafe.Pointer) bool {
	is := slices.Contains(errList, LibNfcError(ptr))
	return is
}

var ErrContextClosed *ErrorContextClosed = &ErrorContextClosed{m: "nfc context is closed"}

type ErrorContextClosed struct {
	m string
}

var _ error = (*ErrorContextClosed)(nil)

func (e ErrorContextClosed) Error() string {
	return e.m
}

var ErrDeviceClosed *ErrorDeviceClosed = &ErrorDeviceClosed{m: "nfc device is closed"}

type ErrorDeviceClosed struct {
	m string
}

var _ error = (*ErrorDeviceClosed)(nil)

func (e ErrorDeviceClosed) Error() string {
	return e.m
}
