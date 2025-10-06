//go:build windows

package cec

import (
	"errors"
)

var (
	notAvailable error = errors.New("SANE is not available for windows")
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	_ = libInitFuncs

	panic(notAvailable)
}

func getSymbol(_ /*sym*/ string) (uintptr, error) {
	return 0, notAvailable
}
