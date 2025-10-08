package cec

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	gostrings "strings"

	"github.com/maitredede/puregolibs/strings"
)

const (
	// default connection timeout in milliseconds
	CEC_DEFAULT_CONNECT_TIMEOUT               = 10000
	DefaultConnectTimeout       time.Duration = CEC_DEFAULT_CONNECT_TIMEOUT * time.Millisecond
)

type Conn struct {
	ptr uintptr
	cfg *NativeConfiguration
}

// func Initialise(cfg *Configuration) (*Conn, error) {
// 	libInit()

// 	var cfgRaw NativeConfiguration
// 	libCecClearConfiguration(&cfgRaw)

// 	if cfg != nil {
// 		cfgRaw.ClientVersion = cfg.ClientVersion
// 		for i := 0; i < OSDNameSize; i++ {
// 			cfgRaw.DeviceName[i] = 0
// 		}
// 		for i, c := range cfg.DeviceName {
// 			if i >= OSDNameSize {
// 				break
// 			}
// 			cfgRaw.DeviceName[i] = byte(c)
// 		}
// 		cfgRaw.DeviceName[OSDNameSize-1] = 0
// 		cfgRaw.DeviceTypes = cfg.DeviceTypes
// 		cfgRaw.AutodectAddress = toBool(cfg.AutodectAddress)
// 		cfgRaw.PhysicalAddress = uint16(cfg.PhysicalAddress)
// 		cfgRaw.BaseDevice = cfg.BaseDevice
// 		cfgRaw.HDMIPort = cfg.HDMIPort
// 		cfgRaw.TVVendor = cfg.TVVendor
// 		cfgRaw.WakeDevices = cfg.WakeDevices
// 		cfgRaw.PowerDevices = cfg.PowerDevices

// 		cfgRaw.ServerVersion = cfg.ServerVersion

// 		cfgRaw.GetSettingsFromROM = toBool(cfg.GetSettingsFromROM)
// 		cfgRaw.ActivateSource = toBool(cfg.ActivateSource)
// 		cfgRaw.PowerOffOnStandby = toBool(cfg.PowerOffOnStandby)
// 		//TODO callbacks
// 		cfgRaw.LogicalAddress = cfg.LogicalAddress
// 		cfgRaw.FirmwareVersion = cfg.FirmwareVersion
// 	}

// 	ret := libCecInitialise(&cfgRaw)
// 	if ret == 0 {
// 		return nil, fmt.Errorf("can't initialise")
// 	}

// 	cfgRet := libCecGetCurrentConfiguration(ret, &cfgRaw)

// 	_ = cfgRet

// 	c := &Conn{
// 		ptr: ret,
// 		cfg: &cfgRaw,
// 	}
// 	return c, nil
// }

func (c *Conn) Close() error {
	if c.ptr == 0 {
		return ErrConnectionIsClosed
	}
	libCecDestroy(c.ptr)
	c.ptr = 0
	return nil
}

func (c *Conn) GetLibInfo() (string, error) {
	if c.ptr == 0 {
		return "", ErrConnectionIsClosed
	}
	info := libCecGetLibInfo(c.ptr)
	return info, nil
}

// Open - open a new connection to the CEC device with the given name
func Open(name string, deviceName string, printLogs bool) (*Conn, error) {
	libInit()

	var cfgRaw NativeConfiguration
	libCecClearConfiguration(&cfgRaw)
	cfgRaw.ClientVersion = VersionCurrent
	cfgRaw.DeviceTypes[0] = DeviceTypeRecordingDevice
	cfgRaw.DeviceName = CDeviceNameString(deviceName)

	ptr := libCecInitialise(&cfgRaw)
	if ptr == 0 {
		return nil, errors.New("cec init failed")
	}

	adapter, err := getAdapter(ptr, name)
	if err != nil {
		defer libCecDestroy(ptr)
		return nil, err
	}

	err = openAdapter(ptr, adapter)
	if err != nil {
		defer libCecDestroy(ptr)
		return nil, err
	}

	cfgRet := libCecGetCurrentConfiguration(ptr, &cfgRaw)
	if cfgRet == 1 {
		//OK
	} else {
		//ERR
	}
	c := Conn{
		ptr: ptr,
		cfg: &cfgRaw,
	}

	return &c, nil
}

func getAdapter(connection uintptr, name string) (Adapter, error) {
	var adapter Adapter

	var deviceList [10]nativeAdapter
	// devicesFound := libcec_find_adapters(connection, &deviceList[0], 10, nil)
	devicesFound := int(libCecFindAdapters(connection, &deviceList[0], byte(len(deviceList)), nil))

	for i := 0; i < devicesFound; i++ {
		device := deviceList[i]
		adapter.Path = strings.GoStringN(uintptr(unsafe.Pointer(&device.path[0])), 1024)
		adapter.Comm = strings.GoStringN(uintptr(unsafe.Pointer(&device.comm[0])), 1024)

		if gostrings.Contains(adapter.Path, name) || gostrings.Contains(adapter.Comm, name) {
			return adapter, nil
		}
	}

	return adapter, errors.New("no Device Found")
}

func openAdapter(connection uintptr, adapter Adapter) error {
	libCecInitVideoStandalone(connection)

	cPort := strings.CString(adapter.Comm)
	result := libCecOpen(connection, uintptr(unsafe.Pointer(cPort)), CEC_DEFAULT_CONNECT_TIMEOUT)
	if result < 1 {
		return fmt.Errorf("adapter '%s' open failed", adapter.Comm)
	}
	return nil
}

// List - list active devices (returns a map of Devices)
func (c *Conn) List() map[string]Device {
	devices := make(map[string]Device)

	activeDevices := c.GetActiveDevices()

	for address, active := range activeDevices {
		if active {
			// var dev Device

			// dev.LogicalAddress = address
			// dev.PhysicalAddress = c.GetDevicePhysicalAddress(address)
			// dev.OSDName = c.GetDeviceOSDName(address)
			// dev.PowerStatus = c.GetDevicePowerStatus(address)
			// dev.ActiveSource = c.IsActiveSource(address)
			// dev.CECVersion = c.GetDeviceCecVersion(address)
			// dev.Language = c.GetDeviceMenuLanguage(address)
			// dev.Vendor = GetVendorByID(c.GetDeviceVendorID(address))

			// devices[logicalNames[address]] = dev

			_ = address
			panic("WIP")
		}
	}
	return devices
}

// GetActiveDevices - returns an array of active devices
func (c *Conn) GetActiveDevices() [16]bool {
	result := libCecGetActiveDevices(c.ptr)
	var devices [16]bool
	// result := C.libcec_get_active_devices(c.connection)
	for i := 0; i < 16; i++ {
		if int(result.addresses[i]) > 0 {
			devices[i] = true
		}
	}

	return devices
}

// GetDevicePhysicalAddress - Get the physical address of the device at
// the given logical address
func (c *Conn) GetDevicePhysicalAddress(address int) string {
	result := libCecGetDevicePhysicalAddress(c.ptr, LogicalAddress(address))

	return fmt.Sprintf("%x.%x.%x.%x", (uint(result)>>12)&0xf, (uint(result)>>8)&0xf, (uint(result)>>4)&0xf, uint(result)&0xf)
}

// GetDeviceOSDName - get the OSD name of the specified device
func (c *Conn) GetDeviceOSDName(address int) string {
	// name := make([]byte, 14)
	// C.libcec_get_device_osd_name(c.connection, C.cec_logical_address(address), (*C.char)(unsafe.Pointer(&name[0])))

	// return string(name)
	panic("WIP")
}

// GetDevicePowerStatus - Get the power status of the device at the
// given address
func (c *Conn) GetDevicePowerStatus(address int) string {
	// result := C.libcec_get_device_power_status(c.connection, C.cec_logical_address(address))

	// // C.CEC_POWER_STATUS_UNKNOWN == error

	// if int(result) == C.CEC_POWER_STATUS_ON {
	// 	return "on"
	// } else if int(result) == C.CEC_POWER_STATUS_STANDBY {
	// 	return "standby"
	// } else if int(result) == C.CEC_POWER_STATUS_IN_TRANSITION_STANDBY_TO_ON {
	// 	return "starting"
	// } else if int(result) == C.CEC_POWER_STATUS_IN_TRANSITION_ON_TO_STANDBY {
	// 	return "shutting down"
	// } else {
	// 	return ""
	// }
	panic("WIP")
}

func (c *Conn) GetDeviceCecVersion(address int) string {
	// result := int(C.libcec_get_device_cec_version(c.connection, C.cec_logical_address(address)))

	// if result == C.CEC_VERSION_1_2 {
	// 	return "1.2"
	// } else if result == C.CEC_VERSION_1_2A {
	// 	return "1.2a"
	// } else if result == C.CEC_VERSION_1_3 {
	// 	return "1.3"
	// } else if result == C.CEC_VERSION_1_3A {
	// 	return "1.3a"
	// } else if result == C.CEC_VERSION_1_4 {
	// 	return "1.4"
	// } else if result == C.CEC_VERSION_2_0 {
	// 	return "2.0"
	// } else if result == C.CEC_VERSION_UNKNOWN {
	// 	return "unknown"
	// }

	// return ""
	panic("WIP")
}
