package gousb

// #include <libusb.h>
import "C"

// TransferStatus contains information about the result of a transfer.
type TransferStatus int32

// Defined Transfer status values.
const (
	TransferCompleted TransferStatus = C.LIBUSB_TRANSFER_COMPLETED
	TransferError     TransferStatus = C.LIBUSB_TRANSFER_ERROR
	TransferTimedOut  TransferStatus = C.LIBUSB_TRANSFER_TIMED_OUT
	TransferCancelled TransferStatus = C.LIBUSB_TRANSFER_CANCELLED
	TransferStall     TransferStatus = C.LIBUSB_TRANSFER_STALL
	TransferNoDevice  TransferStatus = C.LIBUSB_TRANSFER_NO_DEVICE
	TransferOverflow  TransferStatus = C.LIBUSB_TRANSFER_OVERFLOW
)

var transferStatusDescription = map[TransferStatus]string{
	TransferCompleted: "transfer completed without error",
	TransferError:     "transfer failed",
	TransferTimedOut:  "transfer timed out",
	TransferCancelled: "transfer was cancelled",
	TransferStall:     "halt condition detected (endpoint stalled) or control request not supported",
	TransferNoDevice:  "device was disconnected",
	TransferOverflow:  "device sent more data than requested",
}

// String returns a human-readable transfer status.
func (ts TransferStatus) String() string {
	return transferStatusDescription[ts]
}

// Error implements the error interface.
func (ts TransferStatus) Error() string {
	return ts.String()
}
