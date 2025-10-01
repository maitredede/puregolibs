package libnfc

import (
	"errors"
	"fmt"
	"slices"
)

type libNfcError uintptr

var (
	libNfcSuccess      libNfcError = 0
	libNfcEIO          libNfcError = libNfcSuccess - 1 // libNfcError(^uintptr(0))
	libNfcEInvArg      libNfcError = libNfcSuccess - 2
	libNfcEDEVNOTSUPP  libNfcError = libNfcSuccess - 3
	libNfcENOTSUCHDEV  libNfcError = libNfcSuccess - 4
	libNfcEOVFLOW      libNfcError = libNfcSuccess - 5
	libNfcETIMEOUT     libNfcError = libNfcSuccess - 6
	libNfcEOPABORTED   libNfcError = libNfcSuccess - 7
	libNfcENOTIMPL     libNfcError = libNfcSuccess - 8
	libNfcETGRELEASED  libNfcError = libNfcSuccess - 10
	libNfcERFTRANS     libNfcError = libNfcSuccess - 20
	libNfcEMFCAUTHFAIL libNfcError = libNfcSuccess - 30
	libNfcESOFT        libNfcError = libNfcSuccess - 80
	libNfcECHIP        libNfcError = libNfcSuccess - 90
)

func (e libNfcError) Error() error {
	switch e {
	case libNfcEIO:
		return errors.New("libNfcEIO")
	case libNfcEInvArg:
		return errors.New("libNfcEInvArg")
	case libNfcEDEVNOTSUPP:
		return errors.New("libNfcEDEVNOTSUPP")
	case libNfcENOTSUCHDEV:
		return errors.New("libNfcENOTSUCHDEV")
	case libNfcEOVFLOW:
		return errors.New("libNfcEOVFLOW")
	case libNfcETIMEOUT:
		return errors.New("libNfcETIMEOUT")
	case libNfcEOPABORTED:
		return errors.New("libNfcEOPABORTED")
	case libNfcENOTIMPL:
		return errors.New("libNfcENOTIMPL")
	case libNfcETGRELEASED:
		return errors.New("libNfcETGRELEASED")
	case libNfcERFTRANS:
		return errors.New("libNfcERFTRANS")
	case libNfcEMFCAUTHFAIL:
		return errors.New("libNfcEMFCAUTHFAIL")
	case libNfcESOFT:
		return errors.New("libNfcESOFT")
	case libNfcECHIP:
		return errors.New("libNfcECHIP")
	}
	return fmt.Errorf("unkown nfc error: %v", int64(e))
}

func isLibError(ptr uintptr) bool {
	var errList []libNfcError = []libNfcError{
		libNfcEIO,
		libNfcEInvArg,
		libNfcEDEVNOTSUPP,
		libNfcENOTSUCHDEV,
		libNfcEOVFLOW,
		libNfcETIMEOUT,
		libNfcEOPABORTED,
		libNfcENOTIMPL,
		libNfcETGRELEASED,
		libNfcERFTRANS,
		libNfcEMFCAUTHFAIL,
		libNfcESOFT,
		libNfcECHIP,
	}
	is := slices.Contains(errList, libNfcError(ptr))
	return is
}

var (
	libStrError func() string
)

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
