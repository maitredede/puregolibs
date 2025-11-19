package libudev

import "unsafe"

type Monitor unsafe.Pointer

func MonitorNewFromNetlink(udev UDev, name string) Monitor {
	initLib()

	return libudevMonitorNewFromNetlink(udev, name)
}

func MonitorRef(monitor Monitor) Monitor {
	initLib()

	return libudevMonitorRef(monitor)
}

func MonitorUnref(monitor Monitor) Monitor {
	initLib()

	return libudevMonitorUnref(monitor)
}

func MonitorFilterAddMatchSubsystemDevType(mon Monitor, subsystem, devType string) int32 {
	initLib()

	return libudevMonitorFilterAddMatchSubsystemDevType(mon, subsystem, devType)
}

func MonitorEnableReceiving(mon Monitor) int32 {
	initLib()

	return libudevMonitorEnableReceiving(mon)
}

func MonitorGetFd(mon Monitor) uintptr {
	initLib()

	return libudevMonitorGetFd(mon)
}

func MonitorReceiveDevice(mon Monitor) Device {
	initLib()

	return libudevMonitorReceiveDevice(mon)
}
