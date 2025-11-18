package libudev

import "unsafe"

type Enumerate unsafe.Pointer

func EnumerateNew(udev UDev) Enumerate {
	initLib()

	return libudevEnumerateNew(udev)
}

func EnumerateUnref(enumerate Enumerate) {
	initLib()

	libudevEnumerateUnref(enumerate)
}

func EnumerateScanDevices(enumerate Enumerate) int32 {
	initLib()

	return libudevEnumerateScanDevices(enumerate)
}

func EnumerateAddMatchSubsystem(enumerate Enumerate, subsystem string) {
	initLib()

	libudevEnumerateAddMatchSubsystem(enumerate, subsystem)
}

func EnumerateGetListEntry(enumerate Enumerate) ListEntry {
	initLib()

	return libudevEnumerateGetListEntry(enumerate)
}
