package libnfc

import (
	"errors"
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type nfcDevicePtr unsafe.Pointer

type NfcDevice struct {
	ctx *NfcContext

	ptr nfcDevicePtr
}

func (d *NfcDevice) Close() error {
	if d.ptr == nil {
		return ErrDeviceClosed
	}
	libNfcClose(d.ptr)
	d.ptr = nil
	return nil
}

func (d *NfcDevice) Name() string {
	if d.ptr == nil {
		return ""
	}

	return libNfcDeviceGetName(d.ptr)
}

func (d *NfcDevice) ConnString() (string, error) {
	if d.ptr == nil {
		return "", ErrDeviceClosed
	}

	name := libNfcDeviceGetConnString(d.ptr)
	return name, nil
}

func (d *NfcDevice) lastError() LibNfcError {
	if d.ptr == nil {
		panic(ErrDeviceClosed)
	}
	e := libNfcDeviceGetLastError(d.ptr)
	return LibNfcError(e)
}

func (d *NfcDevice) GetInformationAbout() (string, error) {
	if d.ptr == nil {
		return "", ErrDeviceClosed
	}

	var strinfo unsafe.Pointer
	ret := libNfcDeviceGetInformationAbout(d.ptr, &strinfo)
	if isLibErrorInt32(ret) {
		return "", LibNfcError(ret).Error()
	}
	defer libNfcFree(strinfo)

	info := strings.GoString((*byte)(strinfo))
	return info, nil
}

func (d *NfcDevice) Ptr() unsafe.Pointer {
	return unsafe.Pointer(d.ptr)
}

func (d *NfcDevice) InitiatorListPassiveTargets(m Modulation) ([]NfcTarget, error) {
	return nil, errors.New("work in progress")
}

type nativeNfcDevice struct {
	ctx      nfcContextPtr
	drv      nfcDriverPtr
	drvData  unsafe.Pointer
	chipData unsafe.Pointer
	//remaining not needed
}

func (d *NfcDevice) DriverPtr() unsafe.Pointer {
	if d.ptr == nil {
		return nil
	}

	nd := (*nativeNfcDevice)(d.ptr)
	return unsafe.Pointer(nd.drv)
}

func (d *NfcDevice) GoDriver() *Driver {
	drvPtr := nfcDriverPtr(d.DriverPtr())
	if drvPtr == nil {
		return nil
	}

	for _, d := range goDrivers {
		if d.ndPtr == drvPtr {
			return d.driver
		}
	}
	return nil
}
