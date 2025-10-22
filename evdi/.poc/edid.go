package evdi

import (
	"encoding/binary"
	"io"
)

type EDID struct {
	EDIDHeader
	EDIDBasicDisplayParameters
}

func (e *EDID) WriteTo(w io.Writer) (int64, error) {

	sequence := []io.WriterTo{
		&e.EDIDHeader,
		&e.EDIDBasicDisplayParameters,
	}

	var nt int64
	for _, item := range sequence {
		n, err := item.WriteTo(w)
		nt += n
		if err != nil {
			return nt, err
		}
	}
	return nt, nil
}

type EDIDHeader struct {
	ManufacturerID  uint16
	ProductCode     uint16
	SerialNumber    uint32
	ManufactureWeek byte
	ManufactureYear byte

	// EDIDVersion  byte
	// EDIDRevision byte
}

func (e *EDIDHeader) WriteTo(w io.Writer) (int64, error) {

	data := make([]byte, 0, 1024)

	// 0-7 Fixed header pattern
	fixedheader := []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00}
	data = append(data, fixedheader...)

	// 8-9 ManufacturerID
	data = binary.BigEndian.AppendUint16(data, e.ManufacturerID)

	// 10-11 ManufacturerProductCode
	data = binary.LittleEndian.AppendUint16(data, e.ProductCode)

	// 12-15 Serial number
	data = binary.LittleEndian.AppendUint32(data, e.SerialNumber)

	// 16 Week of manufacture
	data = append(data, e.ManufactureWeek)

	// 17 year of manufacture
	data = append(data, e.ManufactureYear)

	// 18 EDID Version
	data = append(data, 0x01)

	// 19 EDID Revision
	data = append(data, 0x04)

	nt, err := w.Write(data)
	return int64(nt), err
}

type InputMode byte

const (
	InputModeAnalog  InputMode = 0
	InputModeDigital InputMode = 1
)

type DigitalBitDepth byte

const (
	BitDepthUndefined DigitalBitDepth = 0b000
	BitDepth6         DigitalBitDepth = 0b001
	BitDepth8         DigitalBitDepth = 0b010
	BitDepth10        DigitalBitDepth = 0b011
	BitDepth12        DigitalBitDepth = 0b100
	BitDepth14        DigitalBitDepth = 0b101
	BitDepth16        DigitalBitDepth = 0b110
	BitDepthReserved  DigitalBitDepth = 0b111
)

type DigitalVideoInterface byte

const (
	InterfaceUndefined   DigitalVideoInterface = 0b0000
	InterfaceDVI         DigitalVideoInterface = 0b0001
	InterfaceHDMIa       DigitalVideoInterface = 0b0010
	InterfaceHDMIb       DigitalVideoInterface = 0b0011
	InterfaceMDDI        DigitalVideoInterface = 0b0100
	InterfaceDisplayPort DigitalVideoInterface = 0b0101
)

type EDIDBasicDisplayParameters struct {
	InputMode InputMode

	DigitalBitDepth       DigitalBitDepth
	DigitalVideoInterface DigitalVideoInterface

	AnalogLevel                     AnalogLevel
	AnalogBlankToBlackSetupExpected bool
	AnalogSeparateSyncSupported     bool
	AnalogCompositeSyncSupported    bool
	AnalogSyncOnGreenSupported      bool
	AnalogVSyncSerrated             bool
}

type AnalogLevel byte

const (
	AnalogLevel_07_03     AnalogLevel = 0b00
	AnalogLevel_0714_0286 AnalogLevel = 0b01
	AnalogLevel_10_04     AnalogLevel = 0b10
	AnalogLevel_07_0      AnalogLevel = 0b11
)

func (e *EDIDBasicDisplayParameters) WriteTo(w io.WriterTo) (int64, error) {

}
