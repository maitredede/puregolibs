package acr122usb

import (
	"fmt"

	"github.com/maitredede/puregolibs/gousb"
	"github.com/maitredede/puregolibs/libnfc"
)

func mkDrvScan(usb *gousb.Context, drv *libnfc.Driver) func() ([]string, libnfc.LibNfcError) {
	return func() ([]string, libnfc.LibNfcError) {
		result := make([]string, 0, 10)
		devs, err := usb.OpenDevices(func(desc *gousb.DeviceDesc) bool {
			ok := false
			var deviceInfo Acr122UsbSupportedDevice
			for _, d := range UsbSupportedDevices {
				if d.VID == desc.Vendor && d.PID == desc.Product {
					ok = true
					deviceInfo = d
					break
				}
			}
			if !ok {
				return false
			}

			// checks from libnfc
			// Make sure there are 2 endpoints available
			// with libusb-win32 we got some null pointers so be robust before looking at endpoints:
			if len(desc.Configs) == 0 {
				return false
			}
			cfg := desc.Configs[1]
			if len(cfg.Interfaces) == 0 {
				return false
			}
			iface := cfg.Interfaces[0]
			if len(iface.AltSettings) == 0 {
				return false
			}
			aset := iface.AltSettings[0]
			if len(aset.Endpoints) < 2 {
				return false
			}

			conString := fmt.Sprintf("%v:%v:%v", drv.Name, desc.Bus, desc.Address)
			_ = deviceInfo

			result = append(result, conString)
			return false
		})

		for _, d := range devs {
			d.Close()
		}
		if err != nil {
			return nil, libnfc.LibNfcESOFT
		}
		return result, libnfc.LibNfcSuccess
	}
}
