package libnfc

import (
	"github.com/maitredede/puregolibs/strings"
)

var (
	// nfc_device *nfc_open(nfc_context *context, const nfc_connstring connstring)
	libNfcOpen func(context uintptr, connstring uintptr) uintptr
	// void nfc_close(nfc_device *pnd);
	libNfcClose func(pnd uintptr)
	// int nfc_abort_command(nfc_device *pnd)
	libNfcAbortCommand func(pnd uintptr) int16
	// size_t nfc_list_devices(nfc_context *context, nfc_connstring connstrings[], size_t connstrings_len)
	libNfcListDevices func(context uintptr, connstrings uintptr, constringsLen uintptr) uintptr
	// int nfc_idle(nfc_device *pnd)
	libNfcIdle func(pnd uintptr) int16

	// int nfc_initiator_init(nfc_device *pnd)
	libNfcInitiatorInit func(pnd uintptr) int16
	// int nfc_initiator_init_secure_element(nfc_device *pnd)
	libNfcInitiatorInitSecureElement func(pnd uintptr) int16
	// int nfc_initiator_select_passive_target(nfc_device *pnd, const nfc_modulation nm, const uint8_t *pbtInitData, const size_t szInitData, nfc_target *pnt)
	libNfcInitiatorSelectPassiveTarget func(pnd uintptr, nm ModulationType, initData *byte, szInitData int32, target uintptr) int16
	// int nfc_initiator_list_passive_targets(nfc_device *pnd, const nfc_modulation nm, nfc_target ant[], const size_t szTargets)
	libNfcInitiatorListPassiveTargets func(pnd uintptr, nm ModulationType, ant *uintptr, szTargets int32) int16
	// int nfc_initiator_poll_target(nfc_device *pnd, const nfc_modulation *pnmTargetTypes, const size_t szTargetTypes, const uint8_t uiPollNr, const uint8_t uiPeriod, nfc_target *pnt)
	libNfcInitiatorPollTarget func(pnd uintptr, targetTypes uintptr, szTargetTypes int32, pollNr byte, period byte, target *uintptr) int16

	// int nfc_initiator_select_dep_target(nfc_device *pnd, const nfc_dep_mode ndm, const nfc_baud_rate nbr, const nfc_dep_info *pndiInitiator, nfc_target *pnt, const int timeout)
	libNfcInitiatorSelectDepTarget func(pnd uintptr, m DepMode, r BaudRate, initiator uintptr, target *uintptr, timeout int16) int16
	// int nfc_initiator_poll_dep_target(nfc_device *pnd, const nfc_dep_mode ndm, const nfc_baud_rate nbr, const nfc_dep_info *pndiInitiator, nfc_target *pnt, const int timeout)
	libNfcInitiatorPollDepTarget func(pnd uintptr) int16
	// int nfc_initiator_deselect_target(nfc_device *pnd)
	libNfcInitiatorDeselectTarget func(pnd uintptr) int16
	// int nfc_initiator_transceive_bytes(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, uint8_t *pbtRx, const size_t szRx, int timeout)
	libNfcInitiatorTransceiveBytes func(pnd uintptr) int16
	// int nfc_initiator_transceive_bits(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar, uint8_t *pbtRx, const size_t szRx, uint8_t *pbtRxPar)
	libNfcInitiatorTransceiveBits func(pnd uintptr) int16
	// int nfc_initiator_transceive_bytes_timed(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, uint8_t *pbtRx, const size_t szRx, uint32_t *cycles)
	libNfcInitiatorTransceiveBytesTimed func(pnd uintptr) int16
	// int nfc_initiator_transceive_bits_timed(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar, uint8_t *pbtRx, const size_t szRx, uint8_t *pbtRxPar, uint32_t *cycles)
	libNfcInitiatorTransceiveBitsTimed func(pnd uintptr) int16
	// int nfc_initiator_target_is_present(nfc_device *pnd, const nfc_target *pnt)
	libNfcInitiatorTargetIsPresent func(pnd uintptr) int16

	// int nfc_target_init(nfc_device *pnd, nfc_target *pnt, uint8_t *pbtRx, const size_t szRx, int timeout)
	libNfcTargetInit func(pnd uintptr) int16
	// int nfc_target_send_bytes(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, int timeout)
	libNfcTargetSendBytes func(pnd uintptr) int16
	// int nfc_target_receive_bytes(nfc_device *pnd, uint8_t *pbtRx, const size_t szRx, int timeout)
	libNfcTargetReceiveBytes func(pnd uintptr) int16
	// int nfc_target_send_bits(nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar)
	libNfcTargetSendBits func(pnd uintptr) int16
	// int nfc_target_receive_bits(nfc_device *pnd, uint8_t *pbtRx, const size_t szRx, uint8_t *pbtRxPar)
	libNfcTargetReceiveBits func(pnd uintptr) int16

	// const char *nfc_device_get_name(nfc_device *pnd)
	libNfcDeviceGetName func(pnd uintptr) string
	// const char *nfc_device_get_connstring(nfc_device *pnd)
	libNfcDeviceGetConnString func(pnd uintptr) string
	// int nfc_device_get_last_error(const nfc_device *pnd)
	libNfcDeviceGetLastError func(pnd uintptr) int16

	// int nfc_device_set_property_int(nfc_device *pnd, const nfc_property property, const int value)
	libNfcDeviceSetPropertyInt func(pnd uintptr, property Property, value int16) int16
	// int nfc_device_set_property_bool(nfc_device *pnd, const nfc_property property, const bool bEnable)
	libNfcDeviceSetPropertyBool func(pnd uintptr, property Property, value bool) int16

	// int nfc_device_get_information_about(nfc_device *pnd, char **buf)
	libNfcDeviceGetInformationAbout func(pnd uintptr, strinfo *uintptr) uintptr
)

type NfcDevice struct {
	ctx *NfcContext

	ptr uintptr
}

func (d *NfcDevice) Close() error {
	if d.ptr == 0 {
		return ErrDeviceClosed
	}
	libNfcClose(d.ptr)
	d.ptr = 0
	return nil
}

func (d *NfcDevice) Name() (string, error) {
	if d.ptr == 0 {
		return "", ErrDeviceClosed
	}

	name := libNfcDeviceGetName(d.ptr)
	return name, nil
}

func (d *NfcDevice) ConnString() (string, error) {
	if d.ptr == 0 {
		return "", ErrDeviceClosed
	}

	name := libNfcDeviceGetConnString(d.ptr)
	return name, nil
}

func (d *NfcDevice) lastError() libNfcError {
	if d.ptr == 0 {
		panic(ErrDeviceClosed)
	}
	e := libNfcDeviceGetLastError(d.ptr)
	return libNfcError(e)
}

func (d *NfcDevice) GetInformationAbout() (string, error) {
	if d.ptr == 0 {
		return "", ErrDeviceClosed
	}

	var strinfo uintptr
	ret := libNfcDeviceGetInformationAbout(d.ptr, &strinfo)
	if isLibError(ret) {
		return "", libNfcError(ret).Error()
	}
	defer libNfcFree(strinfo)

	info := strings.GoString(strinfo)
	return info, nil
}

func (d *NfcDevice) Ptr() uintptr {
	return d.ptr
}
