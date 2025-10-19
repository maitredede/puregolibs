package gousb

// Capability Capabilities supported by an instance of libusb on the current running platform. Test if the loaded library supports a given capability by calling libusb_has_capability()
type Capability uint32

const (
	CapHasCapability              Capability = 0x0000 // The libusb_has_capability() API is available.
	CapHasHotplug                 Capability = 0x0001 // Hotplug support is available on this platform.
	CapHasHIDAccess               Capability = 0x0100 // The library can access HID devices without requiring user intervention. Note that before being able to actually access an HID device, you may still have to call additional libusb functions such as ibusb_detach_kernel_driver().
	CapSupportsDetachKernelDriver Capability = 0x0101 // The library supports detaching of the default USB driver, using libusb_detach_kernel_driver(), if one is set by the OS kernel
)

func HasCapability(cap Capability) bool {
	libInit()

	return libusbHasCapability(cap)
}
