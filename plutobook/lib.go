package plutobook

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "/usr/lib/libSystem.B.dylib"
	case "linux":
		return "libplutobook.so"
	case "windows":
		return "libplutobook.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func libInitFuncs() {
	purego.RegisterLibFunc(&libVersion, initPtr, "plutobook_version")
	purego.RegisterLibFunc(&libVersionString, initPtr, "plutobook_version_string")
	purego.RegisterLibFunc(&libBuildInfo, initPtr, "plutobook_build_info")

	purego.RegisterLibFunc(&libGetErrorMessage, initPtr, "plutobook_get_error_message")
	purego.RegisterLibFunc(&libClearErrorMessage, initPtr, "plutobook_clear_error_message")

	// libCreateSym = mustGetSymbol("plutobook_create")
	// purego.RegisterLibFunc(&libCreate, initPtr, "plutobook_create")
	registerFFICreate()

	purego.RegisterLibFunc(&libDestroy, initPtr, "plutobook_destroy")

	//purego.RegisterLibFunc(&libGetPageSize, initPtr, "plutobook_get_page_size")
	registerFFIGetPageSize()
	registerFFIGetPageSizeAt()
	registerFFIGetPageMargins()
	purego.RegisterLibFunc(&libGetMediaType, initPtr, "plutobook_get_media_type")
}
