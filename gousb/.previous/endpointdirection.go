package gousb

const (
	ControlIn  = 0x80 /*LIBUSB_ENDPOINT_IN*/
	ControlOut = 0x00 /*LIBUSB_ENDPOINT_OUT*/
)

// EndpointDirection defines the direction of data flow - IN (device to host)
// or OUT (host to device).
type EndpointDirection bool

const (
	endpointNumMask       = 0x0f
	endpointDirectionMask = 0x80
	// EndpointDirectionIn marks data flowing from device to host.
	EndpointDirectionIn EndpointDirection = true
	// EndpointDirectionOut marks data flowing from host to device.
	EndpointDirectionOut EndpointDirection = false
)

var endpointDirectionDescription = map[EndpointDirection]string{
	EndpointDirectionIn:  "IN",
	EndpointDirectionOut: "OUT",
}

func (ed EndpointDirection) String() string {
	return endpointDirectionDescription[ed]
}
