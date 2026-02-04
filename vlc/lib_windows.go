//go:build windows

package vlc

import (
	"golang.org/x/sys/windows"
)

var theDLL = windows.NewLazySystemDLL("libvlc.dll")

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		initError = theDLL.Load()
		if initError != nil {
			panic(initError)
		}
		initPtr = theDLL.Handle()
		libInitFuncs()
	}
}

func getSymbol(sym string) (uintptr, error) {
	proc := theDLL.NewProc(sym)
	if err := proc.Find(); err != nil {
		return 0, err
	}
	addr := proc.Addr()
	return addr, nil
}
