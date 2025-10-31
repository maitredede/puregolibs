package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type DeckInfo int32

const (
	CEC_DECK_INFO_PLAY                 DeckInfo = 0x11
	CEC_DECK_INFO_RECORD               DeckInfo = 0x12
	CEC_DECK_INFO_PLAY_REVERSE         DeckInfo = 0x13
	CEC_DECK_INFO_STILL                DeckInfo = 0x14
	CEC_DECK_INFO_SLOW                 DeckInfo = 0x15
	CEC_DECK_INFO_SLOW_REVERSE         DeckInfo = 0x16
	CEC_DECK_INFO_FAST_FORWARD         DeckInfo = 0x17
	CEC_DECK_INFO_FAST_REVERSE         DeckInfo = 0x18
	CEC_DECK_INFO_NO_MEDIA             DeckInfo = 0x19
	CEC_DECK_INFO_STOP                 DeckInfo = 0x1A
	CEC_DECK_INFO_SKIP_FORWARD_WIND    DeckInfo = 0x1B
	CEC_DECK_INFO_SKIP_REVERSE_REWIND  DeckInfo = 0x1C
	CEC_DECK_INFO_INDEX_SEARCH_FORWARD DeckInfo = 0x1D
	CEC_DECK_INFO_INDEX_SEARCH_REVERSE DeckInfo = 0x1E
	CEC_DECK_INFO_OTHER_STATUS         DeckInfo = 0x1F
	CEC_DECK_INFO_OTHER_STATUS_LG      DeckInfo = 0x20
)

func (i DeckInfo) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := unsafe.Pointer(&buff[0])
	libCecDeckStatusToString(i, buffPtr, buffSize)
	return strings.GoStringN((*byte)(buffPtr), int(buffSize))
}
