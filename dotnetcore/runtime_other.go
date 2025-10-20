//go:build !windows

package dotnetcore

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/maitredede/puregolibs/strings"
)

type Runtime struct {
	libFile string
	libPtr  uintptr

	hostfxrInitializeForDotnetCommandLine func(argc int32, argv unsafe.Pointer, parameters unsafe.Pointer, hostCohostContextHandle *hostContextHandle) int32
	hostfxrInitializeForRuntimeConfig     func()
	hostfxrGetRuntimeDelegate             func()
	hostfxrRunApp                         func()
	hostfxrClose                          func(handle hostContextHandle) int32

	handle hostContextHandle
}

func InitializeRuntime(sdkPath string) (*Runtime, error) {
	const libNetHostName = "libnethost.so"
	var libsNetHost []string
	err := filepath.WalkDir(sdkPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == libNetHostName {
			info, err := d.Info()
			if err != nil {
				return err
			}
			libsNetHost = append(libsNetHost, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(libsNetHost) == 0 {
		return nil, fmt.Errorf("missing %s", libNetHostName)
	}
	libNetHostPtr, err := purego.Dlopen(libsNetHost[0], purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}
	var libNetHostGetHostfxrPath func(buffer unsafe.Pointer, bufferSize uint32, _ unsafe.Pointer) int32
	purego.RegisterLibFunc(&libNetHostGetHostfxrPath, libNetHostPtr, "get_hostfxr_path")
	buffer := make([]byte, 4096)
	rc := libNetHostGetHostfxrPath(unsafe.Pointer(&buffer[0]), uint32(len(buffer)), nil)
	if rc != 0 {
		return nil, fmt.Errorf("err %x", rc)
	}
	libFile := strings.GoStringN((uintptr)(unsafe.Pointer(&buffer[0])), len(buffer))

	// libFile := filepath.Join(sdkPath, "libhostfxr.so")
	if _, err := os.Stat(libFile); err != nil {
		return nil, err
	}
	libPtr, err := purego.Dlopen(libFile, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}

	r := &Runtime{
		libFile: libFile,
		libPtr:  libPtr,
	}

	purego.RegisterLibFunc(&r.hostfxrInitializeForDotnetCommandLine, r.libPtr, "hostfxr_initialize_for_dotnet_command_line")
	purego.RegisterLibFunc(&r.hostfxrInitializeForRuntimeConfig, r.libPtr, "hostfxr_initialize_for_runtime_config")
	purego.RegisterLibFunc(&r.hostfxrGetRuntimeDelegate, r.libPtr, "hostfxr_get_runtime_delegate")
	purego.RegisterLibFunc(&r.hostfxrRunApp, r.libPtr, "hostfxr_run_app")
	purego.RegisterLibFunc(&r.hostfxrClose, r.libPtr, "hostfxr_close")

	ret := r.hostfxrInitializeForDotnetCommandLine(0, nil, nil, &r.handle)
	_ = ret
	return r, nil
}

func (r *Runtime) Close() error {
	var err error
	if r.handle != nil {
		ret := r.hostfxrClose(r.handle)
		//TODO handle ret
		_ = ret
		if err != nil {
			return err
		}
	}

	if r.libPtr != 0 {
		if err := purego.Dlclose(r.libPtr); err != nil {
			return err
		}
		r.libPtr = 0
	}
	return nil
}
