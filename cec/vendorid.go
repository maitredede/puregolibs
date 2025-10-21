package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type VendorID int32

const (
	CEC_VENDOR_TOSHIBA        VendorID = 0x000039
	CEC_VENDOR_SAMSUNG        VendorID = 0x0000F0
	CEC_VENDOR_DENON          VendorID = 0x0005CD
	CEC_VENDOR_MARANTZ        VendorID = 0x000678
	CEC_VENDOR_LOEWE          VendorID = 0x000982
	CEC_VENDOR_ONKYO          VendorID = 0x0009B0
	CEC_VENDOR_MEDION         VendorID = 0x000CB8
	CEC_VENDOR_TOSHIBA2       VendorID = 0x000CE7
	CEC_VENDOR_APPLE          VendorID = 0x0010FA
	CEC_VENDOR_PULSE_EIGHT    VendorID = 0x001582
	CEC_VENDOR_HARMAN_KARDON2 VendorID = 0x001950
	CEC_VENDOR_GOOGLE         VendorID = 0x001A11
	CEC_VENDOR_AKAI           VendorID = 0x0020C7
	CEC_VENDOR_AOC            VendorID = 0x002467
	CEC_VENDOR_PANASONIC      VendorID = 0x008045
	CEC_VENDOR_PHILIPS        VendorID = 0x00903E
	CEC_VENDOR_DAEWOO         VendorID = 0x009053
	CEC_VENDOR_YAMAHA         VendorID = 0x00A0DE
	CEC_VENDOR_GRUNDIG        VendorID = 0x00D0D5
	CEC_VENDOR_PIONEER        VendorID = 0x00E036
	CEC_VENDOR_LG             VendorID = 0x00E091
	CEC_VENDOR_SHARP          VendorID = 0x08001F
	CEC_VENDOR_SONY           VendorID = 0x080046
	CEC_VENDOR_BROADCOM       VendorID = 0x18C086
	CEC_VENDOR_SHARP2         VendorID = 0x534850
	CEC_VENDOR_VIZIO          VendorID = 0x6B746D
	CEC_VENDOR_BENQ           VendorID = 0x8065E9
	CEC_VENDOR_HARMAN_KARDON  VendorID = 0x9C645E
	CEC_VENDOR_UNKNOWN        VendorID = 0
)

func (v VendorID) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := unsafe.Pointer(&buff[0])
	libCecVendorIDToString(v, buffPtr, buffSize)
	return strings.GoStringN((*byte)(buffPtr), int(buffSize))
}
