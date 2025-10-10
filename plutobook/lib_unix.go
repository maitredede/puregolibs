//go:build darwin || freebsd || linux || netbsd

package plutobook

import "github.com/ebitengine/purego"

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		name := getSystemLibrary()

		initPtr, initError = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if initError != nil {
			panic(initError)
		}
		libInitFuncs()
	}
}

func getSymbol(sym string) (uintptr, error) {
	return purego.Dlsym(initPtr, sym)
}
