//go:build darwin || freebsd || linux || netbsd

package vlc

import (
	"github.com/ebitengine/purego"
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		initPtr, initError = purego.Dlopen("libvlc.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if initError != nil {
			panic(initError)
		}
		libInitFuncs()
	}
}

func getSymbol(sym string) (uintptr, error) {
	return purego.Dlsym(initPtr, sym)
}
