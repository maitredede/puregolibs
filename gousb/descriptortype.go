package gousb

// DescriptorType identifies the type of a USB descriptor.
type DescriptorType uint8

// Descriptor types defined by the USB spec.
const (
	DescriptorTypeDevice               DescriptorType = 0x01 // Device descriptor. See libusb_device_descriptor
	DescriptorTypeConfig               DescriptorType = 0x02 // Configuration descriptor. See libusb_config_descriptor
	DescriptorTypeString               DescriptorType = 0x03 // String descriptor
	DescriptorTypeInterface            DescriptorType = 0x04 // Interface descriptor. See libusb_interface_descriptor.
	DescriptorTypeEndpoint             DescriptorType = 0x05 // Endpoint descriptor. See libusb_endpoint_descriptor.
	DescriptorTypeInterfaceAssotiation DescriptorType = 0x0b // Interface Association Descriptor. See libusb_interface_association_descriptor*
	DescriptorTypeBOS                  DescriptorType = 0x0f // BOS descriptor
	DescriptorTypeDeviceCapability     DescriptorType = 0x10 // Device Capability descriptor
	DescriptorTypeHID                  DescriptorType = 0x21 // HID descriptor
	DescriptorTypeReport               DescriptorType = 0x22 // HID report descriptor
	DescriptorTypePhysical             DescriptorType = 0x23 // Physical descriptor
	DescriptorTypeHub                  DescriptorType = 0x29 // Hub descriptor
	DescriptorTypeSuperspeedHub        DescriptorType = 0x2a // SuperSpeed Hub descriptor
	DescriptorTypeSsEndpointCompanion  DescriptorType = 0x30 // SuperSpeed Endpoint Companion descriptor
)

var descriptorTypeDescription = map[DescriptorType]string{
	DescriptorTypeDevice:               "device",
	DescriptorTypeConfig:               "configuration",
	DescriptorTypeString:               "string",
	DescriptorTypeInterface:            "interface",
	DescriptorTypeEndpoint:             "endpoint",
	DescriptorTypeInterfaceAssotiation: "iface association",
	DescriptorTypeDeviceCapability:     "capability",
	DescriptorTypeHID:                  "HID",
	DescriptorTypeReport:               "HID report",
	DescriptorTypePhysical:             "physical",
	DescriptorTypeHub:                  "hub",
	DescriptorTypeSuperspeedHub:        "superspeed hub",
	DescriptorTypeSsEndpointCompanion:  "superspeed endpoint companion",
}

func (dt DescriptorType) String() string {
	return descriptorTypeDescription[dt]
}
