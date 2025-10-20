package libnfc

import "unsafe"

type nfcTargetPtr unsafe.Pointer

type NfcTarget interface {
	ptr() nfcTargetPtr
}
