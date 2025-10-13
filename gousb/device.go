package gousb

import (
	"fmt"
	"sync"
	"unsafe"
)

type libusbDevice unsafe.Pointer
type libusbDeviceHandle unsafe.Pointer

type Device struct {
	ctx    *Context
	handle libusbDeviceHandle

	// Embed the device information for easy access
	Desc *DeviceDesc

	// Claimed config
	mu      sync.Mutex
	claimed *Config
}

func (d *Device) Close() error {
	// if !d.handleValid {
	if d.handle == nil {
		return ErrInvalidDevice
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.claimed != nil {
		return fmt.Errorf("can't release the device %s, it has an open config %d", d, d.claimed.Desc.Number)
	}
	d.ctx.closeDev(d)
	d.handle = nil
	return nil
}

// String represents a human readable representation of the device.
func (d *Device) String() string {
	return fmt.Sprintf("vid=%s,pid=%s,bus=%d,addr=%d", d.Desc.Vendor, d.Desc.Product, d.Desc.Bus, d.Desc.Address)
}
