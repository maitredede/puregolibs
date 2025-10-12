package gousb

import (
	"fmt"
	"sort"
)

// InterfaceSetting contains information about a USB interface with a particular
// alternate setting, extracted from the descriptor.
type InterfaceSetting struct {
	// Number is the number of this interface, the same as in InterfaceDesc.
	Number int
	// Alternate is the number of this alternate setting.
	Alternate int
	// Class is the USB-IF (Implementers Forum) class code, as defined by the USB spec.
	Class Class
	// SubClass is the USB-IF (Implementers Forum) subclass code, as defined by the USB spec.
	SubClass Class
	// Protocol is USB protocol code, as defined by the USB spe.c
	Protocol Protocol
	// Endpoints enumerates the endpoints available on this interface with
	// this alternate setting.
	Endpoints map[EndpointAddress]EndpointDesc

	iInterface int // index of a string descriptor describing this interface.
}

func (a InterfaceSetting) sortedEndpointIds() []string {
	var eps []string
	for _, ei := range a.Endpoints {
		eps = append(eps, fmt.Sprintf("%s(%d,%s)", ei.Address, ei.Number, ei.Direction))
	}
	sort.Strings(eps)
	return eps
}

// String returns a human-readable description of the particular
// alternate setting of an interface.
func (a InterfaceSetting) String() string {
	return fmt.Sprintf("Interface %d alternate setting %d (available endpoints: %v)", a.Number, a.Alternate, a.sortedEndpointIds())
}
