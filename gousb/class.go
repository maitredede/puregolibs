package gousb

import "strconv"

// Class represents a USB-IF (Implementers Forum) class or subclass code.
type Class uint8

// Standard classes defined by USB spec, see https://www.usb.org/defined-class-codes
const (
	ClassPerInterface       Class = 0x00
	ClassAudio              Class = 0x01
	ClassComm               Class = 0x02
	ClassHID                Class = 0x03
	ClassPhysical           Class = 0x05
	ClassImage              Class = 0x06
	ClassPTP                Class = ClassImage // legacy name for image
	ClassPrinter            Class = 0x07
	ClassMassStorage        Class = 0x08
	ClassHub                Class = 0x09
	ClassData               Class = 0x0a
	ClassSmartCard          Class = 0x0b
	ClassContentSecurity    Class = 0x0d
	ClassVideo              Class = 0x0e
	ClassPersonalHealthcare Class = 0x0f
	ClassAudioVideo         Class = 0x10
	ClassBillboard          Class = 0x11
	ClassUSBTypeCBridge     Class = 0x12
	ClassDiagnosticDevice   Class = 0xdc
	ClassWireless           Class = 0xe0
	ClassMiscellaneous      Class = 0xef
	ClassApplication        Class = 0xfe
	ClassVendorSpec         Class = 0xff
)

var classDescription = map[Class]string{
	ClassPerInterface:       "per-interface",
	ClassAudio:              "audio",
	ClassComm:               "communications",
	ClassHID:                "human interface device",
	ClassPhysical:           "physical",
	ClassImage:              "image",
	ClassPrinter:            "printer",
	ClassMassStorage:        "mass storage",
	ClassHub:                "hub",
	ClassData:               "data",
	ClassSmartCard:          "smart card",
	ClassContentSecurity:    "content security",
	ClassVideo:              "video",
	ClassPersonalHealthcare: "personal healthcare",
	ClassAudioVideo:         "audio/video",
	ClassBillboard:          "billboard",
	ClassUSBTypeCBridge:     "USB type-C bridge",
	ClassDiagnosticDevice:   "diagnostic device",
	ClassWireless:           "wireless",
	ClassMiscellaneous:      "miscellaneous",
	ClassApplication:        "application-specific",
	ClassVendorSpec:         "vendor-specific",
}

func (c Class) String() string {
	if d, ok := classDescription[c]; ok {
		return d
	}
	return strconv.Itoa(int(c))
}
