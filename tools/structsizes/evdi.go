//go:build linux && cgo

package main

// #cgo CFLAGS: -I./reallibevdi/
// #  define __user
// #include <drm.h>
// #include "module/evdi_drm.h"
// #include "library/evdi_lib.h"
import "C"
import (
	"fmt"

	"github.com/maitredede/puregolibs/libevdi"
)

func dumpEvdi() {
	var v C.struct_evdi_lib_version
	C.evdi_get_lib_version(&v)

	fmt.Printf("evdi: version=%d.%d.%d\n", v.version_major, v.version_minor, v.version_patchlevel)

	items := []valueEntry{
		{nv: uint32(C.DRM_IOCTL_EVDI_CONNECT), nn: "DRM_IOCTL_EVDI_CONNECT", gv: libevdi.DRM_IOCTL_EVDI_CONNECT, gn: "DRM_IOCTL_EVDI_CONNECT"},
		{nv: uint32(C.DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS), nn: "DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS", gv: libevdi.DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS, gn: "DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS"},
	}
	for _, e := range items {
		fmt.Printf("evdi: C.%s=0x%08x drm.%s=0x%08x eq=%v\n", e.nn, e.nv, e.gn, e.gv, e.nv == e.gv)
	}
}
