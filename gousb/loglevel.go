package gousb

import "fmt"

type LogLevel int32

const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelNone:
		return "none"
	case LogLevelError:
		return "error"
	case LogLevelWarning:
		return "warning"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	}
	return fmt.Sprintf("?%d", int(l))
}
