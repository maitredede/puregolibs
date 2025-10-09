package gousb

// Control request type bit fields as defined in the USB spec. All values are
// of uint8 type.  These constants can be used with Device.Control() method to
// specify the type and destination of the control request, e.g.
// `dev.Control(ControlOut|ControlVendor|ControlDevice, ...)`.
const (
	ControlDevice    = 0x00 /* LIBUSB_RECIPIENT_DEVICE */
	ControlInterface = 0x01 /* LIBUSB_RECIPIENT_INTERFACE */
	ControlEndpoint  = 0x02 /* LIBUSB_RECIPIENT_ENDPOINT */
	ControlOther     = 0x03 /* LIBUSB_RECIPIENT_OTHER */
)
