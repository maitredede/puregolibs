package gousb

import "fmt"

// EndpointAddress is a unique identifier for the endpoint, combining the endpoint number and direction.
type EndpointAddress uint8

// String implements the Stringer interface.
func (a EndpointAddress) String() string {
	return fmt.Sprintf("0x%02x", uint8(a))
}
