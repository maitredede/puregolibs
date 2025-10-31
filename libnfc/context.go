package libnfc

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type nfcContextPtr unsafe.Pointer

type NfcContext struct {
	ptr nfcContextPtr
}

func InitContext() (*NfcContext, error) {
	libInit()

	var ptr nfcContextPtr
	libNfcInit(&ptr)
	if ptr == nil {
		return nil, errors.New("unable to init libnfc (malloc)")
	}
	c := &NfcContext{
		ptr: ptr,
	}
	return c, nil
}

func (c *NfcContext) Close() error {
	if c.ptr == nil {
		return ErrContextClosed
	}
	libNfcExit(c.ptr)
	c.ptr = nil
	return nil
}

func (c *NfcContext) OpenDefault() (*NfcDevice, error) {
	return c.openReal(nil)
}

func (c *NfcContext) OpenDevice(connstring string) (*NfcDevice, error) {
	return c.openReal(&connstring)
}

func (c *NfcContext) openReal(connstring *string) (*NfcDevice, error) {
	if c.ptr == nil {
		return nil, ErrContextClosed
	}
	var pnd nfcDevicePtr
	if connstring == nil {
		pnd = libNfcOpen(c.ptr, nil)
	} else {
		cstr := strings.CString(*connstring)
		p := unsafe.Pointer(cstr)
		pnd = libNfcOpen(c.ptr, p)
	}
	if pnd == nil {
		return nil, errors.New("unable to open NFC device")
	}
	if isLibErrorPtr(unsafe.Pointer(pnd)) {
		return nil, LibNfcError(pnd).Error()
	}

	device := &NfcDevice{
		ctx: c,
		ptr: pnd,
	}
	return device, nil
}

const nfcBufSizeConnString = 1024
const maxDeviceCount = 16

type nfcConnString [nfcBufSizeConnString]byte

func (s nfcConnString) String() string {
	bin := s[:]
	nullByte := slices.Index(bin, 0)
	if nullByte != -1 {
		bin = s[:nullByte]
	}
	return string(bin)
}

func (s *nfcConnString) Set(value string) {
	for i := 0; i < nfcBufSizeConnString; i++ {
		s[i] = 0
	}
	binValue := []byte(value)
	m := min(len(binValue), nfcBufSizeConnString)
	for i := 0; i < m; i++ {
		s[i] = binValue[i]
	}
}

func (c *NfcContext) ListDevices() ([]string, error) {
	if c.ptr == nil {
		return nil, ErrContextClosed
	}

	var constrings [maxDeviceCount]nfcConnString
	p := unsafe.Pointer(&constrings[0])
	deviceFound := libNfcListDevices(c.ptr, p, maxDeviceCount)
	slog.Debug(fmt.Sprintf("found %d devices", deviceFound))
	result := make([]string, deviceFound)
	for i := 0; i < int(deviceFound); i++ {
		result[i] = constrings[i].String()
	}
	return result, nil
}
