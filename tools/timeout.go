package tools

import (
	"errors"
	"os"
)

type Timeout interface {
	Timeout() bool
}

func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if os.IsTimeout(err) {
		return true
	}
	if to, ok := err.(Timeout); ok {
		return to.Timeout()
	}
	w := errors.Unwrap(err)
	return IsTimeout(w)
}
