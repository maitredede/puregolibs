//go:build linux && cgo

package main

// // #cgo pkg-config: libdrm
// // #cgo LDFLAGS: -levdi
/*
#cgo CFLAGS: -I../structsizes/reallibevdi/
#cgo LDFLAGS: -levdi
// #  define __user
// // #include <drm.h>
// // #include "module/evdi_drm.h"
// #include "library/evdi_lib.h"
// #include <stdarg.h>
// #include <stdio.h>
// #include <stdlib.h>
#include "main.h"
*/
import "C"
import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/maitredede/puregolibs/drm"
	"github.com/maitredede/puregolibs/libevdi"
	"github.com/maitredede/puregolibs/resources"
	"golang.org/x/sys/unix"
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	// C.set_evdi_log(nil)

	h := C.evdi_open_attached_to_fixed(nil, 0)
	if h == nil {
		slog.Error("failed to open evdi device")
		os.Exit(1)
	}
	defer C.evdi_close(h)

	edid := resources.EDIDv1_1280x800
	edidPtr := (*C.uchar)(unsafe.Pointer(&edid[0]))
	edidLen := C.uint(len(edid))
	C.evdi_connect(h, edidPtr, edidLen, 1280*800)
	defer C.evdi_disconnect(h)

	C.evdi_enable_cursor_events(h, true)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	ctx, stop := signal.NotifyContext(ctxTimeout, os.Interrupt, os.Kill)
	defer stop()

	mainCtx(ctx, h)

	filename := "screenshot.png"
	out, err := os.Create(filename)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer out.Close()

	if err := png.Encode(out, myDisplay); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	out.Sync()

	slog.Info(fmt.Sprintf("written %s", filename))
}

func mainCtx(ctx context.Context, h C.evdi_handle) {
	events := C.my_events(unsafe.Pointer(h))
	pollable := C.evdi_get_event_ready(h)
	for {
		select {
		case <-ctx.Done():
			slog.Warn(ctx.Err().Error())
			return
		default:
		}
		fds := []unix.PollFd{
			{Fd: int32(pollable), Events: unix.POLLIN},
		}
		n, err := unix.Poll(fds, 1000)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			slog.Debug(fmt.Sprintf("poll n=%v err=%v", n, err))
		}
		if n > 0 {
			C.evdi_handle_events(h, &events)
		}
	}
}

//export go_dpms_handler
func go_dpms_handler(cDpmsMode C.int, userData unsafe.Pointer) {
	goDpsMode := libevdi.DpmsMode(cDpmsMode)
	slog.Info(fmt.Sprintf("dpms: %v", goDpsMode))
}

var (
	numerator            int
	myBuffer             *C.struct_evdi_buffer
	myBufferSize         int
	myDisplay            *image.RGBA
	myDisplayPixelFormat drm.Fourcc
	myGrabber            pixelsGrabber
)

//export go_mode_changed_handler
func go_mode_changed_handler(mode C.struct_evdi_mode, userData unsafe.Pointer) {
	slog.Info(fmt.Sprintf("mode: %+v", mode))

	myDisplayPixelFormat = drm.Fourcc(mode.pixel_format)
	slog.Info(fmt.Sprintf("  pixelFormat: %+v", myDisplayPixelFormat))
	if grabber, ok := pixelsGrabberMap[myDisplayPixelFormat]; ok {
		myGrabber = grabber
	} else {
		slog.Warn(fmt.Sprintf("NO PIXEL GRABBER for %s", myDisplayPixelFormat))
		myGrabber = func(src []byte, rect image.Rectangle, dest *image.RGBA) {}
	}

	h := C.evdi_handle(userData)

	if myBuffer != nil {
		C.evdi_unregister_buffer(h, myBuffer.id)
		C.free(myBuffer.buffer)
		C.free(unsafe.Pointer(myBuffer.rects))
		myBuffer = nil
	}

	stride := mode.width
	pitch_mask := C.int(63)
	stride += pitch_mask
	stride &= ^pitch_mask
	stride *= 4
	bytes_per_pixel := mode.bits_per_pixel / 8
	myBufferSize = int(mode.width * mode.height * bytes_per_pixel)

	numerator++
	id := C.int(numerator)
	myBuffer = &C.struct_evdi_buffer{
		id:         id,
		width:      mode.width,
		height:     mode.height,
		stride:     stride,
		rect_count: 16,
		rects:      (*C.struct_evdi_rect)(C.calloc(16, C.size_t(unsafe.Sizeof(C.struct_evdi_rect{})))),
	}
	myBuffer.buffer = C.calloc(1, C.size_t(myBufferSize))
	C.evdi_register_buffer(h, *myBuffer)

	myDisplay = image.NewRGBA(image.Rect(0, 0, int(mode.width), int(mode.height)))

	request_update(h)
}

var (
	buffer_requested bool
)

func request_update(h C.evdi_handle) {
	if buffer_requested {
		return
	}
	buffer_requested = true
	updateReady := C.evdi_request_update(h, myBuffer.id)
	if updateReady {
		grab_pixels(h)
	}
}

func grab_pixels(h C.evdi_handle) {
	if !buffer_requested {
		return
	}

	C.evdi_grab_pixels(h, myBuffer.rects, &myBuffer.rect_count)
	buffer_requested = false

	if myBuffer.rect_count > 0 {
		data := unsafe.Slice((*byte)(myBuffer.buffer), myBufferSize)
		rects := unsafe.Slice(myBuffer.rects, myBuffer.rect_count)
		rectStrings := make([]string, len(rects))
		for i, r := range rects {
			myRect := image.Rect(int(r.x1), int(r.y1), int(r.x2), int(r.y2))
			rectStrings[i] = myRect.String()
			myGrabber(data, myRect, myDisplay)
		}
		slog.Info(fmt.Sprintf("grabPixels: nRects=%d rects=%s", len(rects), strings.Join(rectStrings, ",")))
	}

	request_update(h)
}

//export go_update_ready_handler
func go_update_ready_handler(bufferToBeUpdated C.int, userData unsafe.Pointer) {
	slog.Info(fmt.Sprintf("updateReady: %v", bufferToBeUpdated))
	h := C.evdi_handle(userData)
	grab_pixels(h)
}

//export go_crtc_state_handler
func go_crtc_state_handler(state C.int, userData unsafe.Pointer) {
	slog.Info(fmt.Sprintf("crtcState: %v", state))
}

//export go_cursor_set_handler
func go_cursor_set_handler(cursorSet C.struct_evdi_cursor_set, userData unsafe.Pointer) {
	slog.Info(fmt.Sprintf("cursorSet: %v", cursorSet))
	C.free(unsafe.Pointer(cursorSet.buffer))
}

//export go_cursor_move_handler
func go_cursor_move_handler(cursorMove C.struct_evdi_cursor_move, userData unsafe.Pointer) {
	slog.Info(fmt.Sprintf("cursorMove: %v", cursorMove))
}

//export go_ddcci_data_handler
func go_ddcci_data_handler(ddcciData C.struct_evdi_ddcci_data, userData unsafe.Pointer) {
	slog.Info(fmt.Sprintf("ddcciData: %v", ddcciData))
}

//export go_log
func go_log(userData unsafe.Pointer, msg *C.char) {
	realMsg := C.GoString(msg)
	slog.Info(realMsg)
}
