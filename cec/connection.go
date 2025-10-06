package cec

import (
	"fmt"
)

type Conn struct {
	ptr uintptr
	cfg *NativeConfiguration
}

func Initialise(cfg *Configuration) (*Conn, error) {
	libInit()

	var cfgRaw NativeConfiguration
	libCecClearConfiguration(&cfgRaw)

	if cfg != nil {
		cfgRaw.ClientVersion = cfg.ClientVersion
		for i := 0; i < OSDNameSize; i++ {
			cfgRaw.DeviceName[i] = 0
		}
		for i, c := range cfg.DeviceName {
			if i >= OSDNameSize {
				break
			}
			cfgRaw.DeviceName[i] = byte(c)
		}
		cfgRaw.DeviceName[OSDNameSize-1] = 0
		cfgRaw.DeviceTypes = cfg.DeviceTypes
		cfgRaw.AutodectAddress = toBool(cfg.AutodectAddress)
		cfgRaw.PhysicalAddress = uint16(cfg.PhysicalAddress)
		cfgRaw.BaseDevice = cfg.BaseDevice
		cfgRaw.HDMIPort = cfg.HDMIPort
		cfgRaw.TVVendor = cfg.TVVendor
		cfgRaw.WakeDevices = cfg.WakeDevices
		cfgRaw.PowerDevices = cfg.PowerDevices

		cfgRaw.ServerVersion = cfg.ServerVersion

		cfgRaw.GetSettingsFromROM = toBool(cfg.GetSettingsFromROM)
		cfgRaw.ActivateSource = toBool(cfg.ActivateSource)
		cfgRaw.PowerOffOnStandby = toBool(cfg.PowerOffOnStandby)
		//TODO callbacks
		cfgRaw.LogicalAddress = cfg.LogicalAddress
		cfgRaw.FirmwareVersion = cfg.FirmwareVersion
	}

	ret := libCecInitialise(&cfgRaw)
	if ret == 0 {
		return nil, fmt.Errorf("can't initialise")
	}

	cfgRet := libCecGetCurrentConfiguration(ret, &cfgRaw)

	_ = cfgRet

	c := &Conn{
		ptr: ret,
		cfg: &cfgRaw,
	}
	return c, nil
}

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
