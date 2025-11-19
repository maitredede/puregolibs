package libudev

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	libLck sync.Mutex
	libPtr uintptr
	libErr error
)

func Initialize() error {
	libLck.Lock()
	defer libLck.Unlock()

	initLibNoPanic()
	return libErr
}

func Unload() error {
	libLck.Lock()
	defer libLck.Unlock()

	if libPtr == 0 {
		return nil
	}
	if err := purego.Dlclose(libPtr); err != nil {
		return err
	}
	libPtr = 0
	libErr = nil
	return nil
}

func initLib() {
	libLck.Lock()
	defer libLck.Unlock()

	initLibNoPanic()
	if libErr != nil {
		panic(libErr)
	}
}

func initLibNoPanic() {
	if libErr != nil {
		return
	}
	if libPtr != 0 {
		return
	}

	libPtr, libErr = purego.Dlopen("libudev.so.1", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if libErr != nil {
		return
	}

	purego.RegisterLibFunc(&libudevNew, libPtr, "udev_new")
	purego.RegisterLibFunc(&libudevRef, libPtr, "udev_ref")
	purego.RegisterLibFunc(&libudevUnref, libPtr, "udev_unref")
	purego.RegisterLibFunc(&libudevSetLogFn, libPtr, "udev_set_log_fn")
	purego.RegisterLibFunc(&libudevGetLogPriority, libPtr, "udev_get_log_priority")
	purego.RegisterLibFunc(&libudevSetLogPriority, libPtr, "udev_set_log_priority")
	purego.RegisterLibFunc(&libudevGetUserData, libPtr, "udev_get_userdata")
	purego.RegisterLibFunc(&libudevSetUserData, libPtr, "udev_set_userdata")

	purego.RegisterLibFunc(&libudevEnumerateNew, libPtr, "udev_enumerate_new")
	purego.RegisterLibFunc(&libudevEnumerateRef, libPtr, "udev_enumerate_ref")
	purego.RegisterLibFunc(&libudevEnumerateUnref, libPtr, "udev_enumerate_unref")
	purego.RegisterLibFunc(&libudevEnumerateScanDevices, libPtr, "udev_enumerate_scan_devices")
	purego.RegisterLibFunc(&libudevEnumerateAddMatchSubsystem, libPtr, "udev_enumerate_add_match_subsystem")
	purego.RegisterLibFunc(&libudevEnumerateGetListEntry, libPtr, "udev_enumerate_get_list_entry")

	purego.RegisterLibFunc(&libudevListEntryGetNext, libPtr, "udev_list_entry_get_next")
	purego.RegisterLibFunc(&libudevListEntryGetName, libPtr, "udev_list_entry_get_name")
	purego.RegisterLibFunc(&libudevListEntryGetValue, libPtr, "udev_list_entry_get_value")

	purego.RegisterLibFunc(&libudevDeviceRef, libPtr, "udev_device_ref")
	purego.RegisterLibFunc(&libudevDeviceUnref, libPtr, "udev_device_unref")
	purego.RegisterLibFunc(&libudevDeviceNewFromSyspath, libPtr, "udev_device_new_from_syspath")
	purego.RegisterLibFunc(&libudevDeviceNewFromDevNum, libPtr, "udev_device_new_from_devnum")
	purego.RegisterLibFunc(&libudevDeviceNewFromSubsystemSysname, libPtr, "udev_device_new_from_subsystem_sysname")
	purego.RegisterLibFunc(&libudevDeviceNewFromDeviceID, libPtr, "udev_device_new_from_device_id")
	purego.RegisterLibFunc(&libudevDeviceNewFromEnvironment, libPtr, "udev_device_new_from_environment")

	purego.RegisterLibFunc(&libudevDeviceGetDevPath, libPtr, "udev_device_get_devpath")
	purego.RegisterLibFunc(&libudevDeviceGetSubsystem, libPtr, "udev_device_get_subsystem")
	purego.RegisterLibFunc(&libudevDeviceGetDevType, libPtr, "udev_device_get_devtype")
	purego.RegisterLibFunc(&libudevDeviceGetSysPath, libPtr, "udev_device_get_syspath")
	purego.RegisterLibFunc(&libudevDeviceGetSysName, libPtr, "udev_device_get_sysname")
	purego.RegisterLibFunc(&libudevDeviceGetSysNum, libPtr, "udev_device_get_sysnum")
	purego.RegisterLibFunc(&libudevDeviceGetDevNode, libPtr, "udev_device_get_devnode")
	purego.RegisterLibFunc(&libudevDeviceGetAction, libPtr, "udev_device_get_action")
	purego.RegisterLibFunc(&libudevDeviceGetSysAttrValue, libPtr, "udev_device_get_sysattr_value")

	purego.RegisterLibFunc(&libudevMonitorNewFromNetlink, libPtr, "udev_monitor_new_from_netlink")
	purego.RegisterLibFunc(&libudevMonitorRef, libPtr, "udev_monitor_ref")
	purego.RegisterLibFunc(&libudevMonitorUnref, libPtr, "udev_monitor_unref")
	purego.RegisterLibFunc(&libudevMonitorFilterAddMatchSubsystemDevType, libPtr, "udev_monitor_filter_add_match_subsystem_devtype")
	purego.RegisterLibFunc(&libudevMonitorEnableReceiving, libPtr, "udev_monitor_enable_receiving")
	purego.RegisterLibFunc(&libudevMonitorGetFd, libPtr, "udev_monitor_get_fd")
	purego.RegisterLibFunc(&libudevMonitorReceiveDevice, libPtr, "udev_monitor_receive_device")
}

var (
	libudevNew            func() UDev
	libudevRef            func(udev UDev) UDev
	libudevUnref          func(udev UDev)
	libudevSetLogFn       func(udev UDev, fn unsafe.Pointer)
	libudevGetLogPriority func(udev UDev) int32
	libudevSetLogPriority func(udev UDev, priority int32)
	libudevGetUserData    func(udev UDev) unsafe.Pointer
	libudevSetUserData    func(udev UDev, userData unsafe.Pointer)

	libudevEnumerateNew               func(udev UDev) Enumerate
	libudevEnumerateRef               func(enumerate Enumerate)
	libudevEnumerateUnref             func(enumerate Enumerate)
	libudevEnumerateScanDevices       func(enumerate Enumerate) int32
	libudevEnumerateAddMatchSubsystem func(enumerate Enumerate, subsystem string)
	libudevEnumerateGetListEntry      func(enumerate Enumerate) ListEntry

	libudevListEntryGetNext  func(entry ListEntry) ListEntry
	libudevListEntryGetName  func(entry ListEntry) string
	libudevListEntryGetValue func(entry ListEntry) string

	libudevDeviceRef                     func(device Device) Device
	libudevDeviceUnref                   func(device Device) Device
	libudevDeviceNewFromSyspath          func(udev UDev, syspath string) Device
	libudevDeviceNewFromDevNum           func(udev UDev, typ byte, devnum uint32) Device
	libudevDeviceNewFromSubsystemSysname func(udev UDev, subsystem string, sysname string) Device
	libudevDeviceNewFromDeviceID         func(udev UDev, id string) Device
	libudevDeviceNewFromEnvironment      func(udev UDev) Device

	libudevDeviceGetDevPath   func(device Device) string
	libudevDeviceGetSubsystem func(device Device) string
	libudevDeviceGetDevType   func(device Device) string
	libudevDeviceGetSysPath   func(device Device) string
	libudevDeviceGetSysName   func(device Device) string
	libudevDeviceGetSysNum    func(device Device) string
	libudevDeviceGetDevNode   func(device Device) string
	libudevDeviceGetAction    func(device Device) string

	libudevDeviceGetSysAttrValue func(device Device, sysattr string) string

	libudevMonitorNewFromNetlink                 func(udev UDev, name string) Monitor
	libudevMonitorRef                            func(monitor Monitor) Monitor
	libudevMonitorUnref                          func(monitor Monitor) Monitor
	libudevMonitorFilterAddMatchSubsystemDevType func(mon Monitor, subsystem string, devType string) int32
	libudevMonitorEnableReceiving                func(mon Monitor) int32
	libudevMonitorGetFd                          func(mon Monitor) uintptr
	libudevMonitorReceiveDevice                  func(mon Monitor) Device
)
