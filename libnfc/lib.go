package libnfc

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libnfc.dylib"
	case "linux":
		return "libnfc.so"
	case "windows":
		return "nfc.dll"
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
	libNfcFree func(ptr uintptr)
)

func libInitFuncs() {
	// Library initialization/deinitialization
	purego.RegisterLibFunc(&libNfcInit, initPtr, "nfc_init")
	purego.RegisterLibFunc(&libNfcExit, initPtr, "nfc_exit")
	purego.RegisterLibFunc(&libNfcRegisterDriver, initPtr, "nfc_register_driver")

	// NFC Device/Hardware manipulation
	purego.RegisterLibFunc(&libNfcOpen, initPtr, "nfc_open")
	purego.RegisterLibFunc(&libNfcClose, initPtr, "nfc_close")
	purego.RegisterLibFunc(&libNfcAbortCommand, initPtr, "nfc_abort_command")
	purego.RegisterLibFunc(&libNfcListDevices, initPtr, "nfc_list_devices")
	purego.RegisterLibFunc(&libNfcIdle, initPtr, "nfc_idle")

	// NFC initiator: act as "reader"
	purego.RegisterLibFunc(&libNfcInitiatorInit, initPtr, "nfc_initiator_init")
	purego.RegisterLibFunc(&libNfcInitiatorInitSecureElement, initPtr, "nfc_initiator_init_secure_element")
	purego.RegisterLibFunc(&libNfcInitiatorSelectPassiveTarget, initPtr, "nfc_initiator_select_passive_target")
	purego.RegisterLibFunc(&libNfcInitiatorListPassiveTargets, initPtr, "nfc_initiator_list_passive_targets")
	purego.RegisterLibFunc(&libNfcInitiatorPollTarget, initPtr, "nfc_initiator_poll_target")
	purego.RegisterLibFunc(&libNfcInitiatorSelectDepTarget, initPtr, "nfc_initiator_select_dep_target")
	purego.RegisterLibFunc(&libNfcInitiatorPollDepTarget, initPtr, "nfc_initiator_poll_dep_target")
	purego.RegisterLibFunc(&libNfcInitiatorDeselectTarget, initPtr, "nfc_initiator_deselect_target")
	purego.RegisterLibFunc(&libNfcInitiatorTransceiveBytes, initPtr, "nfc_initiator_transceive_bytes")
	purego.RegisterLibFunc(&libNfcInitiatorTransceiveBits, initPtr, "nfc_initiator_transceive_bits")
	purego.RegisterLibFunc(&libNfcInitiatorTransceiveBytesTimed, initPtr, "nfc_initiator_transceive_bytes_timed")
	purego.RegisterLibFunc(&libNfcInitiatorTransceiveBitsTimed, initPtr, "nfc_initiator_transceive_bits_timed")
	purego.RegisterLibFunc(&libNfcInitiatorTargetIsPresent, initPtr, "nfc_initiator_target_is_present")

	// NFC target: act as tag (i.e. MIFARE Classic) or NFC target device.
	purego.RegisterLibFunc(&libNfcTargetInit, initPtr, "nfc_target_init")
	purego.RegisterLibFunc(&libNfcTargetSendBytes, initPtr, "nfc_target_send_bytes")
	purego.RegisterLibFunc(&libNfcTargetReceiveBytes, initPtr, "nfc_target_receive_bytes")
	purego.RegisterLibFunc(&libNfcTargetSendBits, initPtr, "nfc_target_send_bits")
	purego.RegisterLibFunc(&libNfcTargetReceiveBits, initPtr, "nfc_target_receive_bits")

	// Error reporting
	purego.RegisterLibFunc(&libStrError, initPtr, "nfc_strerror")
	// strerror_r
	// perror
	purego.RegisterLibFunc(&libNfcDeviceGetLastError, initPtr, "nfc_device_get_last_error")

	// Special data accessors
	purego.RegisterLibFunc(&libNfcDeviceGetName, initPtr, "nfc_device_get_name")
	purego.RegisterLibFunc(&libNfcDeviceGetConnString, initPtr, "nfc_device_get_connstring")
	// get supported modulation
	// get supported baud rate
	// get supported baud rate target mode

	// Properties accessors
	purego.RegisterLibFunc(&libNfcDeviceSetPropertyInt, initPtr, "nfc_device_set_property_int")
	purego.RegisterLibFunc(&libNfcDeviceSetPropertyBool, initPtr, "nfc_device_set_property_bool")

	// Misc. functions
	// TODO

	purego.RegisterLibFunc(&libNfcFree, initPtr, "nfc_free")
	purego.RegisterLibFunc(&libVersion, initPtr, "nfc_version")
	purego.RegisterLibFunc(&libNfcDeviceGetInformationAbout, initPtr, "nfc_device_get_information_about")

	// String converter functions
	// str_nfc_modulation_type
	// str_nfc_baud_rate
	// str_nfc_target
}
