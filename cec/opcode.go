package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type OpCode int32

const (
	CEC_OPCODE_ACTIVE_SOURCE                 OpCode = 0x82
	CEC_OPCODE_IMAGE_VIEW_ON                 OpCode = 0x04
	CEC_OPCODE_TEXT_VIEW_ON                  OpCode = 0x0D
	CEC_OPCODE_INACTIVE_SOURCE               OpCode = 0x9D
	CEC_OPCODE_REQUEST_ACTIVE_SOURCE         OpCode = 0x85
	CEC_OPCODE_ROUTING_CHANGE                OpCode = 0x80
	CEC_OPCODE_ROUTING_INFORMATION           OpCode = 0x81
	CEC_OPCODE_SET_STREAM_PATH               OpCode = 0x86
	CEC_OPCODE_STANDBY                       OpCode = 0x36
	CEC_OPCODE_RECORD_OFF                    OpCode = 0x0B
	CEC_OPCODE_RECORD_ON                     OpCode = 0x09
	CEC_OPCODE_RECORD_STATUS                 OpCode = 0x0A
	CEC_OPCODE_RECORD_TV_SCREEN              OpCode = 0x0F
	CEC_OPCODE_CLEAR_ANALOGUE_TIMER          OpCode = 0x33
	CEC_OPCODE_CLEAR_DIGITAL_TIMER           OpCode = 0x99
	CEC_OPCODE_CLEAR_EXTERNAL_TIMER          OpCode = 0xA1
	CEC_OPCODE_SET_ANALOGUE_TIMER            OpCode = 0x34
	CEC_OPCODE_SET_DIGITAL_TIMER             OpCode = 0x97
	CEC_OPCODE_SET_EXTERNAL_TIMER            OpCode = 0xA2
	CEC_OPCODE_SET_TIMER_PROGRAM_TITLE       OpCode = 0x67
	CEC_OPCODE_TIMER_CLEARED_STATUS          OpCode = 0x43
	CEC_OPCODE_TIMER_STATUS                  OpCode = 0x35
	CEC_OPCODE_CEC_VERSION                   OpCode = 0x9E
	CEC_OPCODE_GET_CEC_VERSION               OpCode = 0x9F
	CEC_OPCODE_GIVE_PHYSICAL_ADDRESS         OpCode = 0x83
	CEC_OPCODE_GET_MENU_LANGUAGE             OpCode = 0x91
	CEC_OPCODE_REPORT_PHYSICAL_ADDRESS       OpCode = 0x84
	CEC_OPCODE_SET_MENU_LANGUAGE             OpCode = 0x32
	CEC_OPCODE_DECK_CONTROL                  OpCode = 0x42
	CEC_OPCODE_DECK_STATUS                   OpCode = 0x1B
	CEC_OPCODE_GIVE_DECK_STATUS              OpCode = 0x1A
	CEC_OPCODE_PLAY                          OpCode = 0x41
	CEC_OPCODE_GIVE_TUNER_DEVICE_STATUS      OpCode = 0x08
	CEC_OPCODE_SELECT_ANALOGUE_SERVICE       OpCode = 0x92
	CEC_OPCODE_SELECT_DIGITAL_SERVICE        OpCode = 0x93
	CEC_OPCODE_TUNER_DEVICE_STATUS           OpCode = 0x07
	CEC_OPCODE_TUNER_STEP_DECREMENT          OpCode = 0x06
	CEC_OPCODE_TUNER_STEP_INCREMENT          OpCode = 0x05
	CEC_OPCODE_DEVICE_VENDOR_ID              OpCode = 0x87
	CEC_OPCODE_GIVE_DEVICE_VENDOR_ID         OpCode = 0x8C
	CEC_OPCODE_VENDOR_COMMAND                OpCode = 0x89
	CEC_OPCODE_VENDOR_COMMAND_WITH_ID        OpCode = 0xA0
	CEC_OPCODE_VENDOR_REMOTE_BUTTON_DOWN     OpCode = 0x8A
	CEC_OPCODE_VENDOR_REMOTE_BUTTON_UP       OpCode = 0x8B
	CEC_OPCODE_SET_OSD_STRING                OpCode = 0x64
	CEC_OPCODE_GIVE_OSD_NAME                 OpCode = 0x46
	CEC_OPCODE_SET_OSD_NAME                  OpCode = 0x47
	CEC_OPCODE_MENU_REQUEST                  OpCode = 0x8D
	CEC_OPCODE_MENU_STATUS                   OpCode = 0x8E
	CEC_OPCODE_USER_CONTROL_PRESSED          OpCode = 0x44
	CEC_OPCODE_USER_CONTROL_RELEASE          OpCode = 0x45
	CEC_OPCODE_GIVE_DEVICE_POWER_STATUS      OpCode = 0x8F
	CEC_OPCODE_REPORT_POWER_STATUS           OpCode = 0x90
	CEC_OPCODE_FEATURE_ABORT                 OpCode = 0x00
	CEC_OPCODE_ABORT                         OpCode = 0xFF
	CEC_OPCODE_GIVE_AUDIO_STATUS             OpCode = 0x71
	CEC_OPCODE_GIVE_SYSTEM_AUDIO_MODE_STATUS OpCode = 0x7D
	CEC_OPCODE_REPORT_AUDIO_STATUS           OpCode = 0x7A
	CEC_OPCODE_SET_SYSTEM_AUDIO_MODE         OpCode = 0x72
	CEC_OPCODE_SYSTEM_AUDIO_MODE_REQUEST     OpCode = 0x70
	CEC_OPCODE_SYSTEM_AUDIO_MODE_STATUS      OpCode = 0x7E
	CEC_OPCODE_SET_AUDIO_RATE                OpCode = 0x9A

	/* CEC 1.4 */
	CEC_OPCODE_REPORT_SHORT_AUDIO_DESCRIPTORS  OpCode = 0xA3
	CEC_OPCODE_REQUEST_SHORT_AUDIO_DESCRIPTORS OpCode = 0xA4
	CEC_OPCODE_START_ARC                       OpCode = 0xC0
	CEC_OPCODE_REPORT_ARC_STARTED              OpCode = 0xC1
	CEC_OPCODE_REPORT_ARC_ENDED                OpCode = 0xC2
	CEC_OPCODE_REQUEST_ARC_START               OpCode = 0xC3
	CEC_OPCODE_REQUEST_ARC_END                 OpCode = 0xC4
	CEC_OPCODE_END_ARC                         OpCode = 0xC5
	CEC_OPCODE_CDC                             OpCode = 0xF8
	/* when this opcode is set no opcode will be sent to the device. this is one of the reserved numbers */
	CEC_OPCODE_NONE OpCode = 0xFD
)

func (c OpCode) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecOpCodeToString(c, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}
