package libnfc

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

var (
	// void nfc_init(nfc_context **context)
	libNfcInit func(context *uintptr)
	// void nfc_exit(nfc_context *context)
	libNfcExit func(context uintptr)
)

type NfcContext struct {
	ptr uintptr
}

func InitContext() (*NfcContext, error) {
	libInit()

	var ptr uintptr
	libNfcInit(&ptr)
	if ptr == 0 {
		return nil, errors.New("unable to init libnfc (malloc)")
	}
	c := &NfcContext{
		ptr: ptr,
	}
	return c, nil
}

func (c *NfcContext) Close() error {
	if c.ptr == 0 {
		return ErrContextClosed
	}
	libNfcExit(c.ptr)
	c.ptr = 0
	return nil
}

func (c *NfcContext) OpenDefault() (*NfcDevice, error) {
	return c.openReal(nil)
}

func (c *NfcContext) OpenDevice(connstring string) (*NfcDevice, error) {
	return c.openReal(&connstring)
}

func (c *NfcContext) openReal(connstring *string) (*NfcDevice, error) {
	if c.ptr == 0 {
		return nil, ErrContextClosed
	}
	var pnd uintptr
	if connstring == nil {
		pnd = libNfcOpen(c.ptr, 0)
	} else {
		cstr := strings.CString(*connstring)
		p := uintptr(unsafe.Pointer(cstr))
		pnd = libNfcOpen(c.ptr, p)
	}
	if pnd == 0 {
		return nil, errors.New("unable to open NFC device")
	}
	if isLibError(pnd) {
		return nil, libNfcError(pnd).Error()
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

func (c *NfcContext) ListDevices() ([]string, error) {
	if c.ptr == 0 {
		return nil, ErrContextClosed
	}

	var constrings [maxDeviceCount]nfcConnString
	p := uintptr(unsafe.Pointer(&constrings))
	deviceFound := libNfcListDevices(c.ptr, p, maxDeviceCount)
	slog.Debug(fmt.Sprintf("found %d devices", deviceFound))
	result := make([]string, deviceFound)
	for i := 0; i < int(deviceFound); i++ {
		result[i] = constrings[i].String()
	}
	return result, nil
}
