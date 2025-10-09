package cec

import (
	"time"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type LogMessage struct {
	Message string
	Level   LogLevel
	Time    time.Time
}
type nativeLogMessage struct {
	message uintptr
	level   LogLevel
	time    int64
}

func (n nativeLogMessage) Go() LogMessage {
	str := ""
	if n.message != 0 {
		ptr := *(*uintptr)(unsafe.Pointer(n.message))
		if ptr != 0 {
			str = strings.GoString(ptr)
		}
	}

	msg := LogMessage{
		Message: str,
		Level:   n.level,
		Time:    time.Unix(n.time, 0),
	}
	return msg
}
