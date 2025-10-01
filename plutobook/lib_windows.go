//go:build windows

package plutobook

import "golang.org/x/sys/windows"

var (
	theDLL = windows.NewLazySystemDLL(getSystemLibrary())
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		initError = theDLL.Load()
		if initError != nil {
			panic(initError)
		}
		initPtr = theDLL.Handle()
	}

	libInitFuncs()
}

func mustGetSymbol(sym string) uintptr {
	proc := theDLL.NewProc(sym)
	if err := proc.Find(); err != nil {
		panic(err)
	}
	return proc.Addr()
}
