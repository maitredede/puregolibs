package libevdev

func GetIDBusType(evdev Evdev) uint16 {
	initLib()

	return libevdevGetIDBusType(evdev)
}

func GetIDVendor(evdev Evdev) uint16 {
	initLib()

	return libevdevGetIDVendor(evdev)
}

func GetIDProduct(evdev Evdev) uint16 {
	initLib()

	return libevdevGetIDProduct(evdev)
}

func GetIDVersion(evdev Evdev) uint16 {
	initLib()

	return libevdevGetIDVersion(evdev)
}
