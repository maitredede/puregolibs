package acr122usb

import "github.com/maitredede/puregolibs/libnfc"

type ACR122Device interface {
	GetLedState() (LedState, error)
	SetLed(state LedState) (LedState, error)

	Name() string
	Close() error
}

type realACR122Device struct {
	libnfc.NfcDevice
}

var _ ACR122Device = (*realACR122Device)(nil)

func FromDevice(device *libnfc.NfcDevice) ACR122Device {

	godrv := device.GoDriver()
	if godrv == nil {
		return nil
	}
	if godrv != drv {
		return nil
	}

	real := &realACR122Device{
		NfcDevice: *device,
	}
	return real
}

func (d *realACR122Device) GetLedState() (LedState, error) {
	panic("TODO")
}

func (d *realACR122Device) SetLed(state LedState) (LedState, error) {
	panic("TODO")
}
