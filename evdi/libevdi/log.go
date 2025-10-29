package libevdi

import (
	"fmt"
	"log/slog"
)

var (
	gEvdiLogging *slog.Logger = slog.Default()
)

func SetLogging(logger *slog.Logger) {
	gEvdiLogging = logger
}

func evdiLogError(format string, args ...any) {
	gEvdiLogging.Error("evdi: " + fmt.Sprintf(format, args...))
}

func evdiLogInfo(format string, args ...any) {
	gEvdiLogging.Info("evdi: " + fmt.Sprintf(format, args...))
}

func evdiLogDebug(format string, args ...any) {
	gEvdiLogging.Debug("evdi: " + fmt.Sprintf(format, args...))
}

func assert(condition bool) {
	if !condition {
		panic("assert failed")
	}
}
