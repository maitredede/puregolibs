package libnfc

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

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

var (
	libNfcInit func(context *nfcContextPtr)
	libNfcExit func(context nfcContextPtr)

	libNfcFree  func(ptr unsafe.Pointer)
	libStrError func(pnd nfcDevicePtr) string

	// nfc_device *nfc_open(nfc_context *context, const nfc_connstring connstring)
	libNfcOpen func(context nfcContextPtr, connstring unsafe.Pointer) nfcDevicePtr
	// void nfc_close(nfc_device *pnd);
	libNfcClose func(pnd nfcDevicePtr)
	// int nfc_abort_command(nfc_device *pnd)
	libNfcAbortCommand func(pnd nfcDevicePtr) int32
	// size_t nfc_list_devices(nfc_context *context, nfc_connstring connstrings[], size_t connstrings_len)
	libNfcListDevices func(context nfcContextPtr, connstrings unsafe.Pointer, constringsLen uint32) uint32
	// int nfc_idle(nfc_device *pnd)
	libNfcIdle func(pnd nfcDevicePtr) int32

	// int nfc_initiator_init(nfc_device *pnd)
	libNfcInitiatorInit func(pnd nfcDevicePtr) int32
	// int nfc_initiator_init_secure_element(nfc_device *pnd)
	libNfcInitiatorInitSecureElement func(pnd nfcDevicePtr) int32
	// int nfc_initiator_select_passive_target(nfc_device *pnd, const nfc_modulation nm, const uint8_t *pbtInitData, const size_t szInitData, nfc_target *pnt)
	libNfcInitiatorSelectPassiveTarget func(pnd nfcDevicePtr, nm ModulationType, initData *byte, szInitData int32, target nfcTargetPtr) int32
	// int nfc_initiator_list_passive_targets(nfc_device *pnd, const nfc_modulation nm, nfc_target ant[], const size_t szTargets)
	libNfcInitiatorListPassiveTargets func(pnd nfcDevicePtr, nm ModulationType, ant *uintptr, szTargets int32) int32
	// int nfc_initiator_poll_target(nfc_device *pnd, const nfc_modulation *pnmTargetTypes, const size_t szTargetTypes, const uint8_t uiPollNr, const uint8_t uiPeriod, nfc_target *pnt)
	libNfcInitiatorPollTarget func(pnd nfcDevicePtr, targetTypes uintptr, szTargetTypes int32, pollNr byte, period byte, target nfcTargetPtr) int32

	// int nfc_initiator_select_dep_target(nfc_device *pnd, const nfc_dep_mode ndm, const nfc_baud_rate nbr, const nfc_dep_info *pndiInitiator, nfc_target *pnt, const int timeout)
	libNfcInitiatorSelectDepTarget func(pnd nfcDevicePtr, m DepMode, r BaudRate, initiator uintptr, target nfcTargetPtr, timeout int16) int32
	// int nfc_initiator_poll_dep_target(nfc_device *pnd, const nfc_dep_mode ndm, const nfc_baud_rate nbr, const nfc_dep_info *pndiInitiator, nfc_target *pnt, const int timeout)
	libNfcInitiatorPollDepTarget func(pnd nfcDevicePtr) int32
	// int nfc_initiator_deselect_target(nfc_device *pnd)
	libNfcInitiatorDeselectTarget func(pnd nfcDevicePtr) int32
	// int nfc_initiator_transceive_bytes(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, uint8_t *pbtRx, const size_t szRx, int timeout)
	libNfcInitiatorTransceiveBytes func(pnd nfcDevicePtr) int32
	// int nfc_initiator_transceive_bits(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar, uint8_t *pbtRx, const size_t szRx, uint8_t *pbtRxPar)
	libNfcInitiatorTransceiveBits func(pnd nfcDevicePtr) int32
	// int nfc_initiator_transceive_bytes_timed(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, uint8_t *pbtRx, const size_t szRx, uint32_t *cycles)
	libNfcInitiatorTransceiveBytesTimed func(pnd nfcDevicePtr) int32
	// int nfc_initiator_transceive_bits_timed(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar, uint8_t *pbtRx, const size_t szRx, uint8_t *pbtRxPar, uint32_t *cycles)
	libNfcInitiatorTransceiveBitsTimed func(pnd nfcDevicePtr) int32
	// int nfc_initiator_target_is_present(nfc_device *pnd, const nfc_target *pnt)
	libNfcInitiatorTargetIsPresent func(pnd nfcDevicePtr, target nfcTargetPtr) int32

	// int nfc_target_init(nfc_device *pnd, nfc_target *pnt, uint8_t *pbtRx, const size_t szRx, int timeout)
	libNfcTargetInit func(pnd nfcDevicePtr, target nfcTargetPtr) int32
	// int nfc_target_send_bytes(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, int timeout)
	libNfcTargetSendBytes func(pnd nfcDevicePtr) int32
	// int nfc_target_receive_bytes(nfc_device *pnd, uint8_t *pbtRx, const size_t szRx, int timeout)
	libNfcTargetReceiveBytes func(pnd nfcDevicePtr) int32
	// int nfc_target_send_bits(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar)
	libNfcTargetSendBits func(pnd nfcDevicePtr) int32
	// int nfc_target_receive_bits(nfc_device *pnd, uint8_t *pbtRx, const size_t szRx, uint8_t *pbtRxPar)
	libNfcTargetReceiveBits func(pnd nfcDevicePtr) int32

	// const char *nfc_device_get_name(nfc_device *pnd)
	libNfcDeviceGetName func(pnd nfcDevicePtr) string
	// const char *nfc_device_get_connstring(nfc_device *pnd)
	libNfcDeviceGetConnString func(pnd nfcDevicePtr) string
	// int nfc_device_get_last_error(const nfc_device *pnd)
	libNfcDeviceGetLastError func(pnd nfcDevicePtr) int32

	// int nfc_device_set_property_int(nfc_device *pnd, const nfc_property property, const int value)
	libNfcDeviceSetPropertyInt func(pnd nfcDevicePtr, property Property, value int16) int32
	// int nfc_device_set_property_bool(nfc_device *pnd, const nfc_property property, const bool bEnable)
	libNfcDeviceSetPropertyBool func(pnd nfcDevicePtr, property Property, value bool) int32

	// int nfc_device_get_information_about(nfc_device *pnd, char **buf)
	libNfcDeviceGetInformationAbout func(pnd nfcDevicePtr, strinfo *unsafe.Pointer) int32

	// int nfc_register_driver(const nfc_driver *driver)
	libNfcRegisterDriver func(driver nfcDriverPtr) int32
)
