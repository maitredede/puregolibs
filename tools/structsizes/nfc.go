//go:build linux && cgo

package main

// #include <nfc/nfc.h>
// #include <freefare.h>
/*
typedef enum {
  NOT_INTRUSIVE,
  INTRUSIVE,
  NOT_AVAILABLE,
} scan_type_enum;

struct nfc_driver {
  const char *name;
  const scan_type_enum scan_type;
  size_t (*scan)(const nfc_context *context, nfc_connstring connstrings[], const size_t connstrings_len);
  struct nfc_device *(*open)(const nfc_context *context, const nfc_connstring connstring);
  void (*close)(struct nfc_device *pnd);
  const char *(*strerror)(const struct nfc_device *pnd);

  int (*initiator_init)(struct nfc_device *pnd);
  int (*initiator_init_secure_element)(struct nfc_device *pnd);
  int (*initiator_select_passive_target)(struct nfc_device *pnd,  const nfc_modulation nm, const uint8_t *pbtInitData, const size_t szInitData, nfc_target *pnt);
  int (*initiator_poll_target)(struct nfc_device *pnd, const nfc_modulation *pnmModulations, const size_t szModulations, const uint8_t uiPollNr, const uint8_t btPeriod, nfc_target *pnt);
  int (*initiator_select_dep_target)(struct nfc_device *pnd, const nfc_dep_mode ndm, const nfc_baud_rate nbr, const nfc_dep_info *pndiInitiator, nfc_target *pnt, const int timeout);
  int (*initiator_deselect_target)(struct nfc_device *pnd);
  int (*initiator_transceive_bytes)(struct nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, uint8_t *pbtRx, const size_t szRx, int timeout);
  int (*initiator_transceive_bits)(struct nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar, uint8_t *pbtRx, uint8_t *pbtRxPar);
  int (*initiator_transceive_bytes_timed)(struct nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, uint8_t *pbtRx, const size_t szRx, uint32_t *cycles);
  int (*initiator_transceive_bits_timed)(struct nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar, uint8_t *pbtRx, uint8_t *pbtRxPar, uint32_t *cycles);
  int (*initiator_target_is_present)(struct nfc_device *pnd, const nfc_target *pnt);

  int (*target_init)(struct nfc_device *pnd, nfc_target *pnt, uint8_t *pbtRx, const size_t szRx, int timeout);
  int (*target_send_bytes)(struct nfc_device *pnd, const uint8_t *pbtTx, const size_t szTx, int timeout);
  int (*target_receive_bytes)(struct nfc_device *pnd, uint8_t *pbtRx, const size_t szRxLen, int timeout);
  int (*target_send_bits)(struct nfc_device *pnd, const uint8_t *pbtTx, const size_t szTxBits, const uint8_t *pbtTxPar);
  int (*target_receive_bits)(struct nfc_device *pnd, uint8_t *pbtRx, const size_t szRxLen, uint8_t *pbtRxPar);

  int (*device_set_property_bool)(struct nfc_device *pnd, const nfc_property property, const bool bEnable);
  int (*device_set_property_int)(struct nfc_device *pnd, const nfc_property property, const int value);
  int (*get_supported_modulation)(struct nfc_device *pnd, const nfc_mode mode, const nfc_modulation_type **const supported_mt);
  int (*get_supported_baud_rate)(struct nfc_device *pnd, const nfc_mode mode, const nfc_modulation_type nmt, const nfc_baud_rate **const supported_br);
  int (*device_get_information_about)(struct nfc_device *pnd, char **buf);

  int (*abort_command)(struct nfc_device *pnd);
  int (*idle)(struct nfc_device *pnd);
  int (*powerdown)(struct nfc_device *pnd);
};
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/libnfc"
)

func dumpNfc() {
	var nti C.nfc_target_info
	fmt.Printf("nfc: sizeof nfc_target_info: %d\n", unsafe.Sizeof(nti))

	var nDrv C.nfc_driver
	var goDrv libnfc.NativeDriver
	fmt.Printf("nfc: sizeof C.nfc_driver: %d go.NativeDriver: %d\n", unsafe.Sizeof(nDrv), unsafe.Sizeof(goDrv))

	var ffTag C.MifareTag
	fmt.Printf("freefare: sizeof MifareTag: %d\n", unsafe.Sizeof(ffTag))
}
