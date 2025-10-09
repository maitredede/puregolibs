package cec

import "time"

type Keypress struct {
	Keycode  UserControlCode
	Duration time.Duration
}

type nativeKeyPress struct {
	keycode  UserControlCode
	duration uint32
}

func (n nativeKeyPress) Go() Keypress {
	key := Keypress{
		Keycode:  n.keycode,
		Duration: time.Duration(n.duration) * time.Millisecond,
	}
	return key
}
