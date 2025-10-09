//go:build windows

package sane

import (
	"errors"
)

var (
	errLibNotAvailable error = errors.New("SANE is not available for windows")
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	_ = libInitFuncs

	panic(errLibNotAvailable)
}

func getSymbol(_ /*sym*/ string) (uintptr, error) {
	return 0, errLibNotAvailable
}
