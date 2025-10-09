package gousb

type RequestType int32

const (
	// "Standard" is explicitly omitted, as functionality of standard requests
	// is exposed through higher level operations of gousb.
	// ControlStandard RequestType = (0x00 << 5) /* LIBUSB_REQUEST_TYPE_STANDARD  */
	ControlClass  RequestType = (0x01 << 5) /* LIBUSB_REQUEST_TYPE_CLASS */
	ControlVendor RequestType = (0x02 << 5) /* LIBUSB_REQUEST_TYPE_VENDOR */
	// "Reserved" is explicitly omitted, should not be used.
	// ControlReserved RequestType = (0x03 << 5) /* LIBUSB_REQUEST_TYPE_RESERVED */
)
