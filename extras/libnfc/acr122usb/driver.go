package acr122usb

import (
	"sync"

	"github.com/maitredede/puregolibs/gousb"
	"github.com/maitredede/puregolibs/libnfc"
)

var (
	lck sync.Mutex
	drv *libnfc.Driver
)

func RegisterGoACR122usbDriver(usb *gousb.Context) error {
	lck.Lock()
	defer lck.Unlock()

	if drv == nil {
		drv = &libnfc.Driver{
			Name:     "go_acr122usb",
			ScanMode: libnfc.ScanModeNotIntrusive,
		}
		drv.Scan = mkDrvScan(usb, drv)
		// .open                             = acr122_usb_open,
		// .close                            = acr122_usb_close,
		// .strerror                         = pn53x_strerror,

		// .initiator_init                   = pn53x_initiator_init,
		// .initiator_init_secure_element    = NULL, // No secure-element support
		// .initiator_select_passive_target  = pn53x_initiator_select_passive_target,
		// .initiator_poll_target            = pn53x_initiator_poll_target,
		// .initiator_select_dep_target      = pn53x_initiator_select_dep_target,
		// .initiator_deselect_target        = pn53x_initiator_deselect_target,
		// .initiator_transceive_bytes       = pn53x_initiator_transceive_bytes,
		// .initiator_transceive_bits        = pn53x_initiator_transceive_bits,
		// .initiator_transceive_bytes_timed = pn53x_initiator_transceive_bytes_timed,
		// .initiator_transceive_bits_timed  = pn53x_initiator_transceive_bits_timed,
		// .initiator_target_is_present      = pn53x_initiator_target_is_present,

		// .target_init           = pn53x_target_init,
		// .target_send_bytes     = pn53x_target_send_bytes,
		// .target_receive_bytes  = pn53x_target_receive_bytes,
		// .target_send_bits      = pn53x_target_send_bits,
		// .target_receive_bits   = pn53x_target_receive_bits,

		// .device_set_property_bool     = pn53x_set_property_bool,
		// .device_set_property_int      = pn53x_set_property_int,
		// .get_supported_modulation     = pn53x_get_supported_modulation,
		// .get_supported_baud_rate      = pn53x_get_supported_baud_rate,
		// .device_get_information_about = pn53x_get_information_about,

		// .abort_command  = acr122_usb_abort_command,
		// .idle           = pn53x_idle,
		// /* Even if PN532, PowerDown is not recommended on those devices */
		// .powerdown      = NULL,
	}
	return libnfc.RegisterDriver(drv)
}
