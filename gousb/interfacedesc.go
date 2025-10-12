package gousb

import "fmt"

// InterfaceDesc contains information about a USB interface, extracted from
// the descriptor.
type InterfaceDesc struct {
	// Number is the number of this interface.
	Number int
	// AltSettings is a list of alternate settings supported by the interface.
	AltSettings []InterfaceSetting
}

func (i *InterfaceDesc) altSetting(alt int) (*InterfaceSetting, error) {
	alts := make([]int, len(i.AltSettings))
	for a, s := range i.AltSettings {
		if s.Alternate == alt {
			return &s, nil
		}
		alts[a] = s.Alternate
	}
	return nil, fmt.Errorf("alternate setting %d not found for %s, available alt settings: %v", alt, i, alts)
}

// String returns a human-readable description of the interface descriptor and
// its alternate settings.
func (i InterfaceDesc) String() string {
	return fmt.Sprintf("Interface %d (%d alternate settings)", i.Number, len(i.AltSettings))
}
