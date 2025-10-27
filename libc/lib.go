package libc

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	libCLck sync.Mutex
	libCPtr uintptr
	libCErr error
)

func initCLib() {
	libCLck.Lock()
	defer libCLck.Unlock()

	if libCErr != nil {
		panic(libCErr)
	}
	if libCPtr != 0 {
		return
	}

	libCPtr, libCErr = purego.Dlopen("libc.so.6", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if libCErr != nil {
		panic(libCErr)
	}

	purego.RegisterLibFunc(&libcCalloc, libCPtr, "calloc")
	purego.RegisterLibFunc(&libcFree, libCPtr, "free")
	purego.RegisterLibFunc(&libcPoll, libCPtr, "poll")
}

var (
	libcCalloc func(nmemb int32, size int32) unsafe.Pointer
	libcFree   func(ptr unsafe.Pointer)
	libcPoll   func(fds *Pollfd, nfds uint32, timeout int32) int32
)
