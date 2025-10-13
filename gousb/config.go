package gousb

import (
	"fmt"
	"sync"
)

// Config represents a USB device set to use a particular configuration.
// Only one Config of a particular device can be used at any one time.
// To access device endpoints, claim an interface and it's alternate
// setting number through a call to Interface().
type Config struct {
	Desc ConfigDesc

	dev *Device

	// Claimed interfaces
	mu      sync.Mutex
	claimed map[int]bool
}

// Close releases the underlying device, allowing the caller to switch the device to a different configuration.
func (c *Config) Close() error {
	if c.dev == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.claimed) > 0 {
		var ifs []int
		for k := range c.claimed {
			ifs = append(ifs, k)
		}
		return fmt.Errorf("failed to release %s, interfaces %v are still open", c, ifs)
	}
	c.dev.mu.Lock()
	defer c.dev.mu.Unlock()
	c.dev.claimed = nil
	c.dev = nil
	return nil
}

// String returns the human-readable description of the configuration.
func (c *Config) String() string {
	return fmt.Sprintf("%s,config=%d", c.dev.String(), c.Desc.Number)
}

// Interface claims and returns an interface on a USB device.
// num specifies the number of an interface to claim, and alt specifies the
// alternate setting number for that interface.
func (c *Config) Interface(num, alt int) (*Interface, error) {
	if c.dev == nil {
		return nil, fmt.Errorf("Interface(%d, %d) called on %s after Close", num, alt, c)
	}

	intf, err := c.Desc.intfDesc(num)
	if err != nil {
		return nil, fmt.Errorf("descriptor of interface %d in %s: %v", num, c, err)
	}
	altInfo, err := intf.altSetting(alt)
	if err != nil {
		return nil, fmt.Errorf("descriptor of alternate setting %d of interface %d in %s: %v", alt, num, c, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.claimed[num] {
		return nil, fmt.Errorf("interface %d on %s is already claimed", num, c)
	}

	// Claim the interface
	// if err := c.dev.ctx.libusb.claim(c.dev.handle, uint8(num)); err != nil {
	ret := libusbClaimInterface(c.dev.handle, int32(num))
	if err := errorFromRet(ret); err != nil {
		return nil, fmt.Errorf("failed to claim interface %d on %s: %v", num, c, err)
	}

	// Select an alternate setting if needed (device has multiple alternate settings).
	if len(intf.AltSettings) > 1 {
		//if err := c.dev.ctx.libusb.setAlt(c.dev.handle, uint8(num), uint8(alt)); err != nil {
		ret := libusbSetInterfaceAltSetting(c.dev.handle, int32(num), int32(alt))
		if err := errorFromRet(ret); err != nil {
			// c.dev.ctx.libusb.release(c.dev.handle, uint8(num))
			libusbReleaseInterface(c.dev.handle, int32(num))
			return nil, fmt.Errorf("failed to set alternate config %d on interface %d of %s: %v", alt, num, c, err)
		}
	}

	c.claimed[num] = true
	return &Interface{
		Setting: *altInfo,
		config:  c,
	}, nil
}
