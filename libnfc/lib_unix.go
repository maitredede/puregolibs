//go:build darwin || freebsd || linux || netbsd

package libnfc

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		name := getSystemLibrary()

		initPtr, initError = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if initError != nil {
			err := fmt.Errorf("error loading library %s: %w", name, initError)
			panic(err)
		}
	}

	libInitFuncs()
}

func getSymbol(sym string) (uintptr, error) {
	return purego.Dlsym(initPtr, sym)
}
