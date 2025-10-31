//go:build linux && cgo

package main

/*
#include <drm/drm.h>
#include <drm/drm_fourcc.h>
*/
import "C"
import (
	"fmt"

	"github.com/maitredede/puregolibs/drm"
)

func dumpDrm() {
	items := []valueEntry{
		{nv: uint32(C.DRM_IOCTL_VERSION), nn: "DRM_IOCTL_VERSION", gv: drm.IOCTLVersion, gn: "IOCTLVersion"},
		{nv: uint32(C.DRM_IOCTL_GET_UNIQUE), nn: "DRM_IOCTL_GET_UNIQUE", gv: drm.IOCTLGetUnique, gn: "IOCTLGetUnique"},
		{nv: uint32(C.DRM_IOCTL_GET_MAGIC), nn: "DRM_IOCTL_GET_MAGIC", gv: drm.IOCTLGetMagic, gn: "IOCTLGetMagic"},
		{nv: uint32(C.DRM_IOCTL_IRQ_BUSID), nn: "DRM_IOCTL_IRQ_BUSID", gv: drm.IOCTLIrqBusID, gn: "IOCTLIrqBusID"},

		{nv: uint32(C.DRM_IOCTL_AUTH_MAGIC), nn: "DRM_IOCTL_AUTH_MAGIC", gv: drm.IOCTLAuthMagic, gn: "IOCTLAuthMagic"},

		{nv: uint32(C.DRM_IOCTL_GET_CAP), nn: "DRM_IOCTL_GET_CAP", gv: drm.IOCTLGetCap, gn: "IOCTLGetCap"},

		{nv: uint32(C.DRM_IOCTL_SET_MASTER), nn: "DRM_IOCTL_SET_MASTER", gv: drm.IOCTLSetMaster, gn: "IOCTLSetMaster"},
		{nv: uint32(C.DRM_IOCTL_DROP_MASTER), nn: "DRM_IOCTL_DROP_MASTER", gv: drm.IOCTLDropMaster, gn: "IOCTLDropMaster"},
	}
	for _, e := range items {
		fmt.Printf("drm: C.%s=0x%08x drm.%s=0x%08x eq=%v\n", e.nn, e.nv, e.gn, e.gv, e.nv == e.gv)
	}

	fourccs := []valueEntry{
		{nv: C.DRM_FORMAT_ARGB8888, nn: "DRM_FORMAT_ARGB8888", gv: drm.DRM_FORMAT_ARGB8888, gn: "DRM_FORMAT_ARGB8888"},
	}
	for _, e := range fourccs {
		fmt.Printf("drm: C.%s=0x%08x drm.%s=0x%08x eq=%v\n", e.nn, e.nv, e.gn, e.gv, e.nv == e.gv)
	}
}

type valueEntry struct {
	nv uint32
	nn string
	gv uint32
	gn string
}
