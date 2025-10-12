//go:build windows

package gousb

import (
	"errors"
	"fmt"
	"syscall"
)

var (
	theDLL *syscall.DLL
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		theFile := getSystemLibrary()
		var err error
		theDLL, err = syscall.LoadDLL(theFile)
		if err != nil {
			panic(fmt.Errorf("error loading lib '%s': %w", theFile, err))
		} else {
			initPtr = uintptr(theDLL.Handle)
		}
		libInitFuncs()
	}
}

func getSymbol(sym string) (uintptr, error) {
	if theDLL == nil {
		return 0, errors.New("library not initialized")
	}
	proc, err := theDLL.FindProc(sym)
	if err != nil {
		return 0, err
	}
	return proc.Addr(), nil
}
