package libnfc

import (
	"errors"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
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

func (d *NfcDevice) Name() (string, error) {
	if d.ptr == nil {
		return "", ErrDeviceClosed
	}

	name := libNfcDeviceGetName(d.ptr)
	return name, nil
}

func (d *NfcDevice) ConnString() (string, error) {
	if d.ptr == nil {
		return "", ErrDeviceClosed
	}

	name := libNfcDeviceGetConnString(d.ptr)
	return name, nil
}

func (d *NfcDevice) lastError() libNfcError {
	if d.ptr == nil {
		panic(ErrDeviceClosed)
	}
	e := libNfcDeviceGetLastError(d.ptr)
	return libNfcError(e)
}

func (d *NfcDevice) GetInformationAbout() (string, error) {
	if d.ptr == nil {
		return "", ErrDeviceClosed
	}

	var strinfo unsafe.Pointer
	ret := libNfcDeviceGetInformationAbout(d.ptr, &strinfo)
	if isLibErrorInt32(ret) {
		return "", libNfcError(ret).Error()
	}
	defer libNfcFree(strinfo)

	info := strings.GoString(uintptr(strinfo))
	return info, nil
}

func (d *NfcDevice) Ptr() unsafe.Pointer {
	return unsafe.Pointer(d.ptr)
}

func (d *NfcDevice) InitiatorListPassiveTargets(m Modulation) ([]NfcTarget, error) {
	return nil, errors.New("work in progress")
}
