package libudev

import (
	"testing"
	"unsafe"
)

func TestLog(t *testing.T) {
	u := New()
	defer Unref(u)

	p := GetLogPriority(u)
	t.Logf("current log priority: %v", p)
	SetLogPriority(u, LogDebug)
	p = GetLogPriority(u)
	t.Logf("new log priority: %v", p)

	SetLogFn(u, func(udev UDev, priority int32, file string, line int32, fn, format string, args unsafe.Pointer) {
		t.Logf("level=%v file=%v line=%v fn=%v format=%v args=%v", priority, file, line, fn, format, args)
	})
}
