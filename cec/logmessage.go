package cec

import (
	"time"

	"github.com/maitredede/puregolibs/strings"
)

type LogMessage struct {
	Message string
	Level   LogLevel
	Time    time.Time
}
type nativeLogMessage struct {
	message *byte
	level   LogLevel
	time    int64
}

func (n nativeLogMessage) Go() LogMessage {
	str := ""
	if n.message != nil {
		str = strings.GoString(n.message)
	}

	msg := LogMessage{
		Message: str,
		Level:   n.level,
		Time:    time.Unix(n.time, 0),
	}
	return msg
}
