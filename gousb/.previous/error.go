// Copyright 2013 Google Inc.  All rights reserved.
// Copyright 2016 the gousb Authors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gousb

import (
	"fmt"
)

// #include <libusb.h>
import "C"

var ErrInvalidContext *ErrorInvalidContext = &ErrorInvalidContext{m: "invalid context"}

type ErrorInvalidContext struct {
	m string
}

var _ error = (*ErrorInvalidContext)(nil)

func (e ErrorInvalidContext) Error() string {
	return e.m
}

type USBErrorCode int32

// USBError is an error code from a USB operation. See the list of USBError constants below.
type USBError struct {
	code    USBErrorCode
	name    string
	message string
}

// Error implements the error interface.
func (e USBError) Error() string {
	return fmt.Sprintf("libusb: %s [code %d]", errorString[e], e)
}

// func fromErrNo(errno C.int) error {
// 	err := USBError(errno)
// 	if err == Success {
// 		return nil
// 	}
// 	return err
// }

// Defined result codes.
const (
	Success           USBErrorCode = 0  // C.LIBUSB_SUCCESS
	ErrorIO           USBErrorCode = -1 // C.LIBUSB_ERROR_IO
	ErrorInvalidParam USBErrorCode = -2 // C.LIBUSB_ERROR_INVALID_PARAM
	ErrorAccess       USBErrorCode = -3 // C.LIBUSB_ERROR_ACCESS
	ErrorNoDevice     USBErrorCode = -4 // C.LIBUSB_ERROR_NO_DEVICE
	ErrorNotFound     USBErrorCode = -5 // C.LIBUSB_ERROR_NOT_FOUND
	ErrorBusy         USBErrorCode = -6 // C.LIBUSB_ERROR_BUSY
	ErrorTimeout      USBErrorCode = -7 // C.LIBUSB_ERROR_TIMEOUT
	// ErrorOverflow indicates that the device tried to send more data than was
	// requested and that could fit in the packet buffer.
	ErrorOverflow     USBErrorCode = -8  // C.LIBUSB_ERROR_OVERFLOW
	ErrorPipe         USBErrorCode = -9  // C.LIBUSB_ERROR_PIPE
	ErrorInterrupted  USBErrorCode = -10 // C.LIBUSB_ERROR_INTERRUPTED
	ErrorNoMem        USBErrorCode = -11 // C.LIBUSB_ERROR_NO_MEM
	ErrorNotSupported USBErrorCode = -12 // C.LIBUSB_ERROR_NOT_SUPPORTED
	ErrorOther        USBErrorCode = -99 // C.LIBUSB_ERROR_OTHER
)

var errorString = map[USBErrorCode]string{
	Success:           "success",
	ErrorIO:           "i/o error",
	ErrorInvalidParam: "invalid param",
	ErrorAccess:       "bad access",
	ErrorNoDevice:     "no device",
	ErrorNotFound:     "not found",
	ErrorBusy:         "device or resource busy",
	ErrorTimeout:      "timeout",
	ErrorOverflow:     "overflow",
	ErrorPipe:         "pipe error",
	ErrorInterrupted:  "interrupted",
	ErrorNoMem:        "out of memory",
	ErrorNotSupported: "not supported",
	ErrorOther:        "unknown error",
}
