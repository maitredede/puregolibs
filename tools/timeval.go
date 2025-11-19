//go:build (amd64 || arm64) && linux

package tools

import (
	"time"

	"golang.org/x/sys/unix"
)

func DurationToTimeVal(dur time.Duration) unix.Timeval {
	sec := int64(dur.Seconds())
	uSec := dur.Microseconds() - (sec * 1_000_000)

	tv := unix.Timeval{Sec: sec, Usec: uSec}

	return tv
}
