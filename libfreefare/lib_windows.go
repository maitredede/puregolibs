//go:build windows

package libfreefare

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var (
	theDLL = windows.NewLazySystemDLL(getSystemLibrary())
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		initError = theDLL.Load()
		if initError != nil {
			name := getSystemLibrary()
			err := fmt.Errorf("error loading library %s: %w", name, initError)
			panic(err)

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
