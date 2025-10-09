package gousb

import "strconv"

// Protocol is the interface class protocol, qualified by the values
// of interface class and subclass.
type Protocol uint8

func (p Protocol) String() string {
	return strconv.Itoa(int(p))
}
