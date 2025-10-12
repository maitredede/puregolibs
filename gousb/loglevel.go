package gousb

type LogLevel int32

const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)
