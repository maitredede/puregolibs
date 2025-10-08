package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type UserControlCode int32

const (
	CEC_USER_CONTROL_CODE_SELECT        UserControlCode = 0x00
	CEC_USER_CONTROL_CODE_UP            UserControlCode = 0x01
	CEC_USER_CONTROL_CODE_DOWN          UserControlCode = 0x02
	CEC_USER_CONTROL_CODE_LEFT          UserControlCode = 0x03
	CEC_USER_CONTROL_CODE_RIGHT         UserControlCode = 0x04
	CEC_USER_CONTROL_CODE_RIGHT_UP      UserControlCode = 0x05
	CEC_USER_CONTROL_CODE_RIGHT_DOWN    UserControlCode = 0x06
	CEC_USER_CONTROL_CODE_LEFT_UP       UserControlCode = 0x07
	CEC_USER_CONTROL_CODE_LEFT_DOWN     UserControlCode = 0x08
	CEC_USER_CONTROL_CODE_ROOT_MENU     UserControlCode = 0x09
	CEC_USER_CONTROL_CODE_SETUP_MENU    UserControlCode = 0x0A
	CEC_USER_CONTROL_CODE_CONTENTS_MENU UserControlCode = 0x0B
	CEC_USER_CONTROL_CODE_FAVORITE_MENU UserControlCode = 0x0C
	CEC_USER_CONTROL_CODE_EXIT          UserControlCode = 0x0D
	// reserved: 0x0E 0x0F
	CEC_USER_CONTROL_CODE_TOP_MENU UserControlCode = 0x10
	CEC_USER_CONTROL_CODE_DVD_MENU UserControlCode = 0x11
	// reserved: 0x12 ... 0x1C
	CEC_USER_CONTROL_CODE_NUMBER_ENTRY_MODE   UserControlCode = 0x1D
	CEC_USER_CONTROL_CODE_NUMBER11            UserControlCode = 0x1E
	CEC_USER_CONTROL_CODE_NUMBER12            UserControlCode = 0x1F
	CEC_USER_CONTROL_CODE_NUMBER0             UserControlCode = 0x20
	CEC_USER_CONTROL_CODE_NUMBER1             UserControlCode = 0x21
	CEC_USER_CONTROL_CODE_NUMBER2             UserControlCode = 0x22
	CEC_USER_CONTROL_CODE_NUMBER3             UserControlCode = 0x23
	CEC_USER_CONTROL_CODE_NUMBER4             UserControlCode = 0x24
	CEC_USER_CONTROL_CODE_NUMBER5             UserControlCode = 0x25
	CEC_USER_CONTROL_CODE_NUMBER6             UserControlCode = 0x26
	CEC_USER_CONTROL_CODE_NUMBER7             UserControlCode = 0x27
	CEC_USER_CONTROL_CODE_NUMBER8             UserControlCode = 0x28
	CEC_USER_CONTROL_CODE_NUMBER9             UserControlCode = 0x29
	CEC_USER_CONTROL_CODE_DOT                 UserControlCode = 0x2A
	CEC_USER_CONTROL_CODE_ENTER               UserControlCode = 0x2B
	CEC_USER_CONTROL_CODE_CLEAR               UserControlCode = 0x2C
	CEC_USER_CONTROL_CODE_NEXT_FAVORITE       UserControlCode = 0x2F
	CEC_USER_CONTROL_CODE_CHANNEL_UP          UserControlCode = 0x30
	CEC_USER_CONTROL_CODE_CHANNEL_DOWN        UserControlCode = 0x31
	CEC_USER_CONTROL_CODE_PREVIOUS_CHANNEL    UserControlCode = 0x32
	CEC_USER_CONTROL_CODE_SOUND_SELECT        UserControlCode = 0x33
	CEC_USER_CONTROL_CODE_INPUT_SELECT        UserControlCode = 0x34
	CEC_USER_CONTROL_CODE_DISPLAY_INFORMATION UserControlCode = 0x35
	CEC_USER_CONTROL_CODE_HELP                UserControlCode = 0x36
	CEC_USER_CONTROL_CODE_PAGE_UP             UserControlCode = 0x37
	CEC_USER_CONTROL_CODE_PAGE_DOWN           UserControlCode = 0x38
	// reserved: 0x39 ... 0x3F
	CEC_USER_CONTROL_CODE_POWER        UserControlCode = 0x40
	CEC_USER_CONTROL_CODE_VOLUME_UP    UserControlCode = 0x41
	CEC_USER_CONTROL_CODE_VOLUME_DOWN  UserControlCode = 0x42
	CEC_USER_CONTROL_CODE_MUTE         UserControlCode = 0x43
	CEC_USER_CONTROL_CODE_PLAY         UserControlCode = 0x44
	CEC_USER_CONTROL_CODE_STOP         UserControlCode = 0x45
	CEC_USER_CONTROL_CODE_PAUSE        UserControlCode = 0x46
	CEC_USER_CONTROL_CODE_RECORD       UserControlCode = 0x47
	CEC_USER_CONTROL_CODE_REWIND       UserControlCode = 0x48
	CEC_USER_CONTROL_CODE_FAST_FORWARD UserControlCode = 0x49
	CEC_USER_CONTROL_CODE_EJECT        UserControlCode = 0x4A
	CEC_USER_CONTROL_CODE_FORWARD      UserControlCode = 0x4B
	CEC_USER_CONTROL_CODE_BACKWARD     UserControlCode = 0x4C
	CEC_USER_CONTROL_CODE_STOP_RECORD  UserControlCode = 0x4D
	CEC_USER_CONTROL_CODE_PAUSE_RECORD UserControlCode = 0x4E
	// reserved: 0x4F
	CEC_USER_CONTROL_CODE_ANGLE                     UserControlCode = 0x50
	CEC_USER_CONTROL_CODE_SUB_PICTURE               UserControlCode = 0x51
	CEC_USER_CONTROL_CODE_VIDEO_ON_DEMAND           UserControlCode = 0x52
	CEC_USER_CONTROL_CODE_ELECTRONIC_PROGRAM_GUIDE  UserControlCode = 0x53
	CEC_USER_CONTROL_CODE_TIMER_PROGRAMMING         UserControlCode = 0x54
	CEC_USER_CONTROL_CODE_INITIAL_CONFIGURATION     UserControlCode = 0x55
	CEC_USER_CONTROL_CODE_SELECT_BROADCAST_TYPE     UserControlCode = 0x56
	CEC_USER_CONTROL_CODE_SELECT_SOUND_PRESENTATION UserControlCode = 0x57
	// reserved: 0x58 ... 0x5F
	CEC_USER_CONTROL_CODE_PLAY_FUNCTION               UserControlCode = 0x60
	CEC_USER_CONTROL_CODE_PAUSE_PLAY_FUNCTION         UserControlCode = 0x61
	CEC_USER_CONTROL_CODE_RECORD_FUNCTION             UserControlCode = 0x62
	CEC_USER_CONTROL_CODE_PAUSE_RECORD_FUNCTION       UserControlCode = 0x63
	CEC_USER_CONTROL_CODE_STOP_FUNCTION               UserControlCode = 0x64
	CEC_USER_CONTROL_CODE_MUTE_FUNCTION               UserControlCode = 0x65
	CEC_USER_CONTROL_CODE_RESTORE_VOLUME_FUNCTION     UserControlCode = 0x66
	CEC_USER_CONTROL_CODE_TUNE_FUNCTION               UserControlCode = 0x67
	CEC_USER_CONTROL_CODE_SELECT_MEDIA_FUNCTION       UserControlCode = 0x68
	CEC_USER_CONTROL_CODE_SELECT_AV_INPUT_FUNCTION    UserControlCode = 0x69
	CEC_USER_CONTROL_CODE_SELECT_AUDIO_INPUT_FUNCTION UserControlCode = 0x6A
	CEC_USER_CONTROL_CODE_POWER_TOGGLE_FUNCTION       UserControlCode = 0x6B
	CEC_USER_CONTROL_CODE_POWER_OFF_FUNCTION          UserControlCode = 0x6C
	CEC_USER_CONTROL_CODE_POWER_ON_FUNCTION           UserControlCode = 0x6D
	// reserved: 0x6E ... 0x70
	CEC_USER_CONTROL_CODE_F1_BLUE   UserControlCode = 0x71
	CEC_USER_CONTROL_CODE_F2_RED    UserControlCode = 0x72
	CEC_USER_CONTROL_CODE_F3_GREEN  UserControlCode = 0x73
	CEC_USER_CONTROL_CODE_F4_YELLOW UserControlCode = 0x74
	CEC_USER_CONTROL_CODE_F5        UserControlCode = 0x75
	CEC_USER_CONTROL_CODE_DATA      UserControlCode = 0x76
	// reserved: 0x77 ... 0xFF
	CEC_USER_CONTROL_CODE_AN_RETURN        UserControlCode = 0x91 // return (Samsung)
	CEC_USER_CONTROL_CODE_AN_CHANNELS_LIST UserControlCode = 0x96 // channels list (Samsung)
	CEC_USER_CONTROL_CODE_MAX              UserControlCode = 0x96
	CEC_USER_CONTROL_CODE_UNKNOWN          UserControlCode = 0xFF
)

func (c UserControlCode) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecUserControlKeyToString(c, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}
