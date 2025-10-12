package gousb

import "fmt"

// ID represents a vendor or product ID.
type ID uint16

// String returns a hexadecimal ID.
func (id ID) String() string {
	return fmt.Sprintf("%04x", int(id))
}
