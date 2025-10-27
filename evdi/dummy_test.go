package evdi

import (
	"context"
	"testing"
	"time"
	"unsafe"

	"github.com/maitredede/puregolibs/libc"
)

func TestDummy(t *testing.T) {
	initLib()

	SetLogging(func(s string) { t.Log(s) })

	device, err := OpenAttachedToNone()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { device.Close() })

	var bufferId int32 = 0
	var width int32 = 1980
	var height int32 = 1200
	var bufferMemory unsafe.Pointer = libc.CAlloc(height*width*4, 1)
	if bufferMemory == nil {
		t.Fatal("memory allocation failed")
	}

	edid := EDIDv1_1280x800
	eventsHandler := Events{
		ModeChanged: func(mode EvdiMode, userdata unsafe.Pointer) {
			t.Logf("Mode changed: %dx%d @%dHz (bpp: %d, format: 0x%x)\n",
				mode.Width, mode.Height, mode.RefreshRate,
				mode.BitsPerPixel, mode.PixelFormat)

			// 1. Free previous buffer
			libevdiUnregisterBuffer(device.h, bufferId)
			if bufferMemory != nil {
				libc.Free(bufferMemory)
			}

			// 2. Update dimensions
			width = mode.Width
			height = mode.Height

			// 3. Allocate new buffer (black)
			bufferMemory = libc.CAlloc(width*height*4, 1)
			if bufferMemory == nil {
				t.Fatal("memory allocation failed")
			}

			// 4. Register new buffer
			buffer := evdiBuffer{
				id:        bufferId,
				buffer:    (*byte)(bufferMemory),
				width:     width,
				height:    height,
				stride:    width * 4,
				rects:     nil,
				rectCount: 0,
			}

			libevdiRegisterBuffer(device.h, &buffer)
		},
		UpdateReady: func(bufferToUpdate int32, userData unsafe.Pointer) {
			t.Logf("updateReady: %d", bufferToUpdate)

			// rects := make([]evdiRect, 16)
			// numRects := int32(len(rects))
			var rectsPtr *evdiRect
			var numRects int32

			libevdiGrabPixels(device.h, &rectsPtr, &numRects)
			if numRects > 0 {
				t.Logf("Update received for %d rects", numRects)
			}

			buffer := evdiBuffer{
				id:        bufferId,
				buffer:    (*byte)(bufferMemory),
				width:     width,
				height:    height,
				stride:    width * 4,
				rects:     rectsPtr,
				rectCount: numRects,
			}

			libevdiRegisterBuffer(device.h, &buffer)
		},

		Dpms: func(dpmsMode int32, userData unsafe.Pointer) {
			t.Logf("dpms: 0x%x", dpmsMode)
		},
		CrtcState: func(state int32, userData unsafe.Pointer) {
			t.Logf("crtc: 0x%x", state)
		},
		CursorSet: func(cursorSet CursorSet, userData unsafe.Pointer) {
			t.Logf("cursorSet: %+v", cursorSet)
		},
		CursorMove: func(cursorMove CursorMove, userData unsafe.Pointer) {
			t.Logf("cursorMove: %+v", cursorMove)
		},
		DdcCiData: func(data DdcciData, userData unsafe.Pointer) {
			t.Logf("ddciData: %+v", data)
		},
	}

	nativeEvts, dispose := buildNativeEvents(eventsHandler)
	t.Cleanup(dispose)

	t.Log("evdi: connecting")
	maxUint32 := uint32(0xFFFFFFFF)
	libEvdiConnect(device.h, &edid[0], uint32(len(edid)), maxUint32)
	t.Log("evdi: connected")
	t.Cleanup(func() { libEvdiDisconnect(device.h) })

	buffer := evdiBuffer{
		id:        bufferId,
		buffer:    (*byte)(bufferMemory),
		width:     width,
		height:    height,
		stride:    width * 4, // 4 bytes par pixel (BGRA)
		rects:     nil,
		rectCount: 0,
	}
	libevdiRegisterBuffer(device.h, &buffer)
	t.Logf("registered buffer %dx%d (black)", width, height)

	newSelectable := libEvdiGetEventReady(device.h)
	t.Logf("evdi: eventGetReady: %v", newSelectable)

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	// epollFd, err := unix.EpollCreate1(0)
	// if err != nil {
	// 	t.Fatalf("failed to create epoll: %v", err)
	// }
	// defer unix.Close(epollFd)

	// Configurer epoll pour surveiller les événements READ sur ce fd
	// event := unix.EpollEvent{
	// 	Events: unix.EPOLLIN, // Équivalent de ev::READ
	// 	Fd:     int32(newSelectable),
	// }

	// if err := unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, int(newSelectable), &event); err != nil {
	// 	t.Fatalf("failed to add fd to epoll: %v", err)
	// }
	// defer func() {
	// 	// Retirer le fd d'epoll
	// 	if err := unix.EpollCtl(epollFd, unix.EPOLL_CTL_DEL, int(newSelectable), nil); err != nil {
	// 		t.Fatalf("failed to remove fd from epoll: %v", err)
	// 	}
	// }()

	// events := []unix.EpollEvent{
	// 	event,
	// }
	var frameCount int
	for {
		select {
		case <-ctx.Done():
			t.Log(ctx.Err())
			return
		default:
		}

		fds := []libc.Pollfd{
			libc.Pollfd{Fd: int32(newSelectable), Events: libc.POLLIN},
		}

		n := libc.Poll(fds, 100)
		// Attendre les événements (timeout de 100ms pour vérifier done)
		// n, err := unix.EpollWait(epollFd, events, 100)
		// if err != nil {
		// 	if err == unix.EINTR {
		// 		continue
		// 	}
		// 	t.Logf("epoll wait error: %v", err)
		// 	continue
		// }

		// time.Sleep(1 * time.Millisecond)
		isEPollIn := (fds[0].Revents & libc.POLLIN) == libc.POLLIN
		t.Logf("f=%010d polled n=%v isEPollIn=%v", frameCount, n, isEPollIn)

		if n > 0 {
			libEvdiHandleEvents(device.h, nativeEvts)
		}
		frameCount++
		if frameCount%60 == 0 {
			t.Logf("f=%010d requestUpdate", frameCount)
			libevdiRequestUpdate(device.h, bufferId)
		}
	}
}
