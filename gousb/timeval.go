package gousb

import "time"

type NativeTimeval struct {
	Sec  int64
	NSec int64
}

func (t NativeTimeval) Time() time.Time {
	return time.Unix(t.Sec, t.NSec)
}
