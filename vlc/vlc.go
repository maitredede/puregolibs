package vlc

import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

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
	ptr := libvlcNew(int32(argc), unsafe.Pointer(&argv[0]))

	if ptr == nil {
		return nil, fmt.Errorf("new instance error: %s", libvlcErrmsg())
	}

	return &Instance{ptr: ptr}, nil
}

func (v *Instance) Close() error {
	libInit()

	libvlcRelease(v.ptr)

	return nil
}
