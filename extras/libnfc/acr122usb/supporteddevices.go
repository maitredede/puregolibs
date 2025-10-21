package acr122usb

import "github.com/maitredede/puregolibs/gousb"

type Acr122UsbSupportedDevice struct {
	VID  gousb.ID
	PID  gousb.ID
	Name string
}

var UsbSupportedDevices []Acr122UsbSupportedDevice = []Acr122UsbSupportedDevice{
	{VID: 0x072F, PID: 0x2200, Name: "ACS ACR122"},
	{VID: 0x072F, PID: 0x90CC, Name: "Touchatag"},
	{VID: 0x072F, PID: 0x2214, Name: "ACS ACR1222"},
}
