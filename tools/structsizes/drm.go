package main

/*
#include <drm/drm.h>
*/
import "C"
import (
	"fmt"

	"github.com/maitredede/puregolibs/evdi/libevdi/drm"
)

func dumpDrm() {
	items := []drmCodeEntry{
		{nv: uint32(C.DRM_IOCTL_VERSION), nn: "DRM_IOCTL_VERSION", gv: drm.IOCTLVersion, gn: "IOCTLVersion"},
		{nv: uint32(C.DRM_IOCTL_GET_UNIQUE), nn: "DRM_IOCTL_GET_UNIQUE", gv: drm.IOCTLGetUnique, gn: "IOCTLGetUnique"},
		{nv: uint32(C.DRM_IOCTL_GET_MAGIC), nn: "DRM_IOCTL_GET_MAGIC", gv: drm.IOCTLGetMagic, gn: "IOCTLGetMagic"},
		{nv: uint32(C.DRM_IOCTL_IRQ_BUSID), nn: "DRM_IOCTL_IRQ_BUSID", gv: drm.IOCTLIrqBusID, gn: "IOCTLIrqBusID"},

		{nv: uint32(C.DRM_IOCTL_AUTH_MAGIC), nn: "DRM_IOCTL_AUTH_MAGIC", gv: drm.IOCTLAuthMagic, gn: "IOCTLAuthMagic"},
	}
	for _, e := range items {
		fmt.Printf("drm: C.%s=0x%08x drm.%s=0x%08x eq=%v\n", e.nn, e.nv, e.gn, e.gv, e.nv == e.gv)
	}
}

type drmCodeEntry struct {
	nv uint32
	nn string
	gv uint32
	gn string
}
