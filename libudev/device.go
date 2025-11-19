package libudev

import "unsafe"

type Device unsafe.Pointer

func DeviceNewFromSyspath(udev UDev, syspath string) Device {
	initLib()

	return libudevDeviceNewFromSyspath(udev, syspath)
}

func DeviceRef(device Device) Device {
	initLib()

	return libudevDeviceRef(device)
}

func DeviceUnref(device Device) Device {
	initLib()

	return libudevDeviceUnref(device)
}

func DeviceGetDevNode(device Device) string {
	initLib()

	return libudevDeviceGetDevNode(device)
}

func DeviceGetAction(device Device) string {
	initLib()

	return libudevDeviceGetAction(device)
}

func DeviceGetSysAttrValue(device Device, sysattr string) string {
	initLib()

	return libudevDeviceGetSysAttrValue(device, sysattr)
}

func DeviceGetSubsystem(device Device) string {
	initLib()

	return libudevDeviceGetSubsystem(device)
}

func DeviceGetDevType(device Device) string {
	initLib()

	return libudevDeviceGetDevType(device)
}
