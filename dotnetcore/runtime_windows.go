//go:build windows

package dotnetcore

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

type Runtime struct {
	libFile string
	libDll  *syscall.DLL

	hostfxrInitializeForDotnetCommandLine func(argc int32, argv unsafe.Pointer, parameters unsafe.Pointer, hostCohostContextHandle *hostContextHandle) int32
	hostfxrInitializeForRuntimeConfig     func()
	hostfxrGetRuntimeDelegate             func()
	hostfxrRunApp                         func()
	hostfxrClose                          func(handle hostContextHandle) int32

	handle hostContextHandle
}

func InitializeRuntime(sdkPath string) (*Runtime, error) {
	libFile := filepath.Join(sdkPath, "coreclr.dll")
	if _, err := os.Stat(libFile); err != nil {
		return nil, err
	}
	libDll, err := syscall.LoadDLL(libFile)
	if err != nil {
		return nil, err
	}

	_ = libDll

	panic("WIP")
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

	if r.libDll != nil {
		if err := r.libDll.Release(); err != nil {
			return err
		}
		r.libDll = nil
	}
	return nil
}
