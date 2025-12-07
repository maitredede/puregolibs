package vlc

import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type libvlcInstance unsafe.Pointer

type Instance struct {
	ptr libvlcInstance
}

func New(args []string) (*Instance, error) {
	libInit()
	argc := len(args)
	argv := make([]unsafe.Pointer, argc)
	for i := 0; i < argc; i++ {
		v := strings.CString(args[i])
		argv[i] = unsafe.Pointer(v)
	}
	var ptr libvlcInstance
	if argc > 0 {
		ptr = libvlcNew(int32(argc), unsafe.Pointer(&argv[0]))
	} else {
		ptr = libvlcNew(int32(argc), nil)
	}

	if ptr == nil {
		return nil, fmt.Errorf("new instance error: %s", libvlcErrmsg())
	}

	return &Instance{ptr: ptr}, nil
}

func (v *Instance) Close() error {
	libvlcRelease(v.ptr)

	return nil
}
