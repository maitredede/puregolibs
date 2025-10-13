package gousb

import "fmt"

type libusbInterface struct {
	altsetting    *libusbInterfaceDescriptor
	numAltsetting int32
}

// Interface is a representation of a claimed interface with a particular setting.
// To access device endpoints use InEndpoint() and OutEndpoint() methods.
// The interface should be Close()d after use.
type Interface struct {
	Setting InterfaceSetting

	config *Config
}

func (i *Interface) String() string {
	return fmt.Sprintf("%s,if=%d,alt=%d", i.config, i.Setting.Number, i.Setting.Alternate)
}

// Close releases the interface.
func (i *Interface) Close() {
	if i.config == nil {
		return
	}
	//i.config.dev.ctx.libusb.release(i.config.dev.handle, uint8(i.Setting.Number))
	libusbReleaseInterface(i.config.dev.handle, int32(i.Setting.Number))

	i.config.mu.Lock()
	defer i.config.mu.Unlock()
	delete(i.config.claimed, i.Setting.Number)
	i.config = nil
}

func (i *Interface) openEndpoint(epAddr EndpointAddress) (*endpoint, error) {
	var ep EndpointDesc
	ep, ok := i.Setting.Endpoints[epAddr]
	if !ok {
		return nil, fmt.Errorf("%s does not have endpoint with address %s. Available endpoints: %v", i, epAddr, i.Setting.sortedEndpointIds())
	}
	return &endpoint{
		InterfaceSetting: i.Setting,
		Desc:             ep,
		h:                i.config.dev.handle,
		ctx:              i.config.dev.ctx,
	}, nil
}

// InEndpoint prepares an IN endpoint for transfer.
func (i *Interface) InEndpoint(epNum int) (*InEndpoint, error) {
	if i.config == nil {
		return nil, fmt.Errorf("InEndpoint(%d) called on %s after Close", epNum, i)
	}
	ep, err := i.openEndpoint(EndpointAddress(0x80 | epNum))
	if err != nil {
		return nil, err
	}
	return &InEndpoint{
		endpoint: ep,
	}, nil
}

// OutEndpoint prepares an OUT endpoint for transfer.
func (i *Interface) OutEndpoint(epNum int) (*OutEndpoint, error) {
	if i.config == nil {
		return nil, fmt.Errorf("OutEndpoint(%d) called on %s after Close", epNum, i)
	}
	ep, err := i.openEndpoint(EndpointAddress(epNum))
	if err != nil {
		return nil, err
	}
	return &OutEndpoint{
		endpoint: ep,
	}, nil
}
