package cec

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libcec.dylib"
	case "linux":
		return "libcec.so"
	case "windows":
		return "cec.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

var (
	libCecInitialise         func(configuration *NativeConfiguration) uintptr
	libCecDestroy            func(connection uintptr)
	libCecOpen               func(connection uintptr, port uintptr, timeout uint32) int32
	libCecClose              func(connection uintptr)
	libCecClearConfiguration func(configuration *NativeConfiguration)
	libCecSetCallbacks       func(connection uintptr, callbacks /*ICECCallbacks* */ uintptr, cbParam /*void* */ uintptr) int32
	// libCecDisableCallbacks          func(connection uintptr) int32
	// libCecEnableCallbacks           func(connection uintptr, cbParam /*void* */ uintptr, callbacks /*ICECCallbacks* */ uintptr) int32
	libCecFindAdapters             func(connection uintptr, deviceList *nativeAdapter, bufSize byte, devicePath *byte) int8
	libCecPingAdapters             func(connection uintptr) int32
	libCecStartBootloader          func(connection uintptr) int32
	libCecPowerOnDevice            func(connection uintptr, address LogicalAddress) int32
	libCecStandbyDevice            func(connection uintptr, address LogicalAddress) int32
	libCecSetActiveSource          func() int32
	libCecSetDeckControlMode       func() int32
	libCecSetDeckInfo              func() int32
	libCecSetInactiveView          func() int32
	libCecSetMenuState             func() int32
	libCecTransmit                 func() int32
	libCecSetLogicalAddress        func() int32
	libCecSetPhysicalAddress       func() int32
	libCecSetOsdString             func() int32
	libCecSwitchMonitoring         func() int32
	libCecGetDeviceCecVersion      func(connection uintptr, address LogicalAddress) CecVersion
	libCecGetDeviceMenuLanguage    func() int32
	libCecGetDeviceVendorID        func() int32
	libCecGetDevicePhysicalAddress func(connection uintptr, address LogicalAddress) uint16
	libCecGetActiveSource          func() int32
	libCecIsActiveSource           func() int32
	libCecGetDevicePowerStatus     func() int32
	libCecPollDevice               func() int32
	libCecGetActiveDevices         func(connection uintptr) *LogicalAddresses
	libCecIsActiveDevice           func() int32
	libCecIsActiveDeviceType       func() int32
	libCecSetHdmiPort              func() int32
	libCecVolumeUp                 func() int32
	libCecVolumeDown               func() int32
	libCecMuteAudio                func() int32
	libCecSendKeypress             func() int32
	libCecSendKeyRelease           func() int32
	libCecGetDeviceOsdName         func() int32
	libCecSetStreamPathLogical     func() int32
	libCecSetStreamPathPhysical    func() int32
	libCecGetLogicalAddresses      func() int32
	libCecGetCurrentConfiguration  func(connection uintptr, configuration *NativeConfiguration) int32
	libCecCanSaveConfiguration     func() int32
	// libCecCanPersistConfiguration   func() int32
	// libCecPersistConfiguration      func() int32
	libCecSetConfiguration          func() int32
	libCecRescanDevices             func()
	libCecIsLibCecActiveSource      func() int32
	libCecGetDeviceInformation      func() int32
	libCecGetLibInfo                func(connection uintptr) string
	libCecInitVideoStandalone       func(connection uintptr)
	libCecGetAdapterVendorID        func(connection uintptr) uint16
	libCecGetAdapterProductID       func(connection uintptr) uint16
	libCecAudioToggleMute           func(connection uintptr) byte
	libCecAudioMute                 func(connection uintptr) byte
	libCecAudioUnmute               func(connection uintptr) byte
	libCecAudioGetStatus            func(connection uintptr) byte
	libCecDetectAdapters            func(connection uintptr) int8
	libCecMenuStateToString         func(state MenuState, buf uintptr, bufSize int32)
	libCecCecVersionToString        func(version CecVersion, buf uintptr, bufSize int32)
	libCecPowerStatusToString       func(status PowerStatus, buf uintptr, bufSize int32)
	libCecLogicalAddressToString    func(address LogicalAddress, buf uintptr, bufSize int32)
	libCecDeckControlModeToString   func(mode DeckControlMode, buf uintptr, bufSize int32)
	libCecDeckStatusToString        func(status DeckInfo, buf uintptr, bufSize int32)
	libCecOpCodeToString            func(opcode OpCode, buf uintptr, bufSize int32)
	libCecSystemAudioStatusToString func(status SystemAudioStatus, buf uintptr, bufSize int32)
	libCecVendorIDToString          func(vendor VendorID, buf uintptr, bufSize int32)
	libCecUserControlKeyToString    func(key UserControlCode, buf uintptr, bufSize int32)
	libCecAdapterTypeToString       func(aType AdapterType, buf uintptr, bufSize int32)
	libCecVersionToString           func(version uint32, buf uintptr, bufSize int32)
)

func libInitFuncs() {
	//purego.RegisterLibFunc(&libCecInitialise, initPtr, "libcec_initialise")
	var cifInitialise ffi.Cif
	if status := ffi.PrepCif(&cifInitialise, ffi.DefaultAbi, 1, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}
	symInitialize, err := getSymbol("libcec_initialise")
	if err != nil {
		panic(err)
	}
	libCecInitialise = func(configuration *NativeConfiguration) uintptr {
		var ret uintptr

		retPtr := unsafe.Pointer(&ret)
		cfgPtr := unsafe.Pointer(configuration)

		argsPtr := []unsafe.Pointer{
			unsafe.Pointer(&cfgPtr),
		}

		ffi.Call(&cifInitialise, symInitialize, retPtr, argsPtr...)
		return ret
	}

	purego.RegisterLibFunc(&libCecDestroy, initPtr, "libcec_destroy")
	purego.RegisterLibFunc(&libCecOpen, initPtr, "libcec_open")
	purego.RegisterLibFunc(&libCecClose, initPtr, "libcec_close")
	purego.RegisterLibFunc(&libCecClearConfiguration, initPtr, "libcec_clear_configuration")
	purego.RegisterLibFunc(&libCecSetCallbacks, initPtr, "libcec_set_callbacks")
	// purego.RegisterLibFunc(&libCecDisableCallbacks, initPtr, "libcec_disable_callbacks")
	// purego.RegisterLibFunc(&libCecEnableCallbacks, initPtr, "libcec_enable_callbacks")
	purego.RegisterLibFunc(&libCecFindAdapters, initPtr, "libcec_find_adapters")
	purego.RegisterLibFunc(&libCecPingAdapters, initPtr, "libcec_ping_adapters")
	purego.RegisterLibFunc(&libCecStartBootloader, initPtr, "libcec_start_bootloader")
	purego.RegisterLibFunc(&libCecPowerOnDevice, initPtr, "libcec_power_on_devices")
	purego.RegisterLibFunc(&libCecStandbyDevice, initPtr, "libcec_standby_devices")
	purego.RegisterLibFunc(&libCecSetActiveSource, initPtr, "libcec_set_active_source")
	purego.RegisterLibFunc(&libCecSetDeckControlMode, initPtr, "libcec_set_deck_control_mode")
	purego.RegisterLibFunc(&libCecSetDeckInfo, initPtr, "libcec_set_deck_info")
	purego.RegisterLibFunc(&libCecSetInactiveView, initPtr, "libcec_set_inactive_view")
	purego.RegisterLibFunc(&libCecSetMenuState, initPtr, "libcec_set_menu_state")
	purego.RegisterLibFunc(&libCecTransmit, initPtr, "libcec_transmit")
	purego.RegisterLibFunc(&libCecSetLogicalAddress, initPtr, "libcec_set_logical_address")
	purego.RegisterLibFunc(&libCecSetPhysicalAddress, initPtr, "libcec_set_physical_address")
	purego.RegisterLibFunc(&libCecSetOsdString, initPtr, "libcec_set_osd_string")
	purego.RegisterLibFunc(&libCecSwitchMonitoring, initPtr, "libcec_switch_monitoring")
	purego.RegisterLibFunc(&libCecGetDeviceCecVersion, initPtr, "libcec_get_device_cec_version")
	purego.RegisterLibFunc(&libCecGetDeviceMenuLanguage, initPtr, "libcec_get_device_menu_language")
	purego.RegisterLibFunc(&libCecGetDeviceVendorID, initPtr, "libcec_get_device_vendor_id")
	purego.RegisterLibFunc(&libCecGetDevicePhysicalAddress, initPtr, "libcec_get_device_physical_address")
	purego.RegisterLibFunc(&libCecGetActiveSource, initPtr, "libcec_get_active_source")
	purego.RegisterLibFunc(&libCecIsActiveSource, initPtr, "libcec_is_active_source")
	purego.RegisterLibFunc(&libCecGetDevicePowerStatus, initPtr, "libcec_get_device_power_status")
	purego.RegisterLibFunc(&libCecPollDevice, initPtr, "libcec_poll_device")
	purego.RegisterLibFunc(&libCecGetActiveDevices, initPtr, "libcec_get_active_devices")
	purego.RegisterLibFunc(&libCecIsActiveDevice, initPtr, "libcec_is_active_device")
	purego.RegisterLibFunc(&libCecIsActiveDeviceType, initPtr, "libcec_is_active_device_type")
	purego.RegisterLibFunc(&libCecSetHdmiPort, initPtr, "libcec_set_hdmi_port")
	purego.RegisterLibFunc(&libCecVolumeUp, initPtr, "libcec_volume_up")
	purego.RegisterLibFunc(&libCecVolumeDown, initPtr, "libcec_volume_down")
	purego.RegisterLibFunc(&libCecMuteAudio, initPtr, "libcec_mute_audio")
	purego.RegisterLibFunc(&libCecSendKeypress, initPtr, "libcec_send_keypress")
	purego.RegisterLibFunc(&libCecSendKeyRelease, initPtr, "libcec_send_key_release")
	purego.RegisterLibFunc(&libCecGetDeviceOsdName, initPtr, "libcec_get_device_osd_name")
	purego.RegisterLibFunc(&libCecSetStreamPathLogical, initPtr, "libcec_set_stream_path_logical")
	purego.RegisterLibFunc(&libCecSetStreamPathPhysical, initPtr, "libcec_set_stream_path_physical")
	purego.RegisterLibFunc(&libCecGetLogicalAddresses, initPtr, "libcec_get_logical_addresses")
	purego.RegisterLibFunc(&libCecGetCurrentConfiguration, initPtr, "libcec_get_current_configuration")
	purego.RegisterLibFunc(&libCecCanSaveConfiguration, initPtr, "libcec_can_save_configuration")
	// purego.RegisterLibFunc(&libCecCanPersistConfiguration, initPtr, "libcec_can_persist_configuration")
	// purego.RegisterLibFunc(&libCecPersistConfiguration, initPtr, "libcec_persist_configuration")
	purego.RegisterLibFunc(&libCecSetConfiguration, initPtr, "libcec_set_configuration")
	purego.RegisterLibFunc(&libCecRescanDevices, initPtr, "libcec_rescan_devices")
	purego.RegisterLibFunc(&libCecIsLibCecActiveSource, initPtr, "libcec_is_libcec_active_source")
	purego.RegisterLibFunc(&libCecGetDeviceInformation, initPtr, "libcec_get_device_information")
	purego.RegisterLibFunc(&libCecGetLibInfo, initPtr, "libcec_get_lib_info")
	purego.RegisterLibFunc(&libCecInitVideoStandalone, initPtr, "libcec_init_video_standalone")
	purego.RegisterLibFunc(&libCecGetAdapterVendorID, initPtr, "libcec_get_adapter_vendor_id")
	purego.RegisterLibFunc(&libCecGetAdapterProductID, initPtr, "libcec_get_adapter_product_id")
	purego.RegisterLibFunc(&libCecAudioToggleMute, initPtr, "libcec_audio_toggle_mute")
	purego.RegisterLibFunc(&libCecAudioMute, initPtr, "libcec_audio_mute")
	purego.RegisterLibFunc(&libCecAudioUnmute, initPtr, "libcec_audio_unmute")
	purego.RegisterLibFunc(&libCecAudioGetStatus, initPtr, "libcec_audio_get_status")
	purego.RegisterLibFunc(&libCecDetectAdapters, initPtr, "libcec_detect_adapters")
	purego.RegisterLibFunc(&libCecMenuStateToString, initPtr, "libcec_menu_state_to_string")
	purego.RegisterLibFunc(&libCecCecVersionToString, initPtr, "libcec_cec_version_to_string")
	purego.RegisterLibFunc(&libCecPowerStatusToString, initPtr, "libcec_power_status_to_string")
	purego.RegisterLibFunc(&libCecLogicalAddressToString, initPtr, "libcec_logical_address_to_string")
	purego.RegisterLibFunc(&libCecDeckControlModeToString, initPtr, "libcec_deck_control_mode_to_string")
	purego.RegisterLibFunc(&libCecDeckStatusToString, initPtr, "libcec_deck_status_to_string")
	purego.RegisterLibFunc(&libCecOpCodeToString, initPtr, "libcec_opcode_to_string")
	purego.RegisterLibFunc(&libCecSystemAudioStatusToString, initPtr, "libcec_system_audio_status_to_string")
	purego.RegisterLibFunc(&libCecVendorIDToString, initPtr, "libcec_vendor_id_to_string")
	purego.RegisterLibFunc(&libCecUserControlKeyToString, initPtr, "libcec_user_control_key_to_string")
	purego.RegisterLibFunc(&libCecAdapterTypeToString, initPtr, "libcec_adapter_type_to_string")
	purego.RegisterLibFunc(&libCecVersionToString, initPtr, "libcec_version_to_string")
}
