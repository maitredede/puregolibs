package libevdi

import "log/slog"

func SetLogging(logger func(msg string)) {
	libInit()
	slog.Error("logging not configured")
}
