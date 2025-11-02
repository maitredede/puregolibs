//go:build linux

package libevdi

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/resources"
	"golang.org/x/sys/unix"
)

var EDIDv1_1280x800 = resources.EDIDv1_1280x800

type bufferData struct {
	buffer *evdiBuffer
	data   []byte
}

func OpenAttachedToNone() (*Handle, error) {
	return OpenAttachedTo("")
}

func (h *Handle) DummyEDID() []byte {
	return EDIDv1_1280x800[:]
}

func (h *Handle) RunDummy(ctx context.Context) error {
	evdiLogDebug("dummy: enabling cursor events")
	h.EnableCursorEvents(true)

	var currentBuffer int32

	eventsHandler := EventHandlers{
		ModeChanged: func(mode Mode, userdata any) {
			evdiLogInfo("event: mode: %+v", mode)
			for bid, _ /* buff */ := range h.buffersMap {
				h.UnregisterBuffer(bid)
				delete(h.buffersMap, bid)
			}

			for i := 0; i < bufferCount; i++ {
				id := int32(i)
				buffer := &evdiBuffer{
					id:     id,
					width:  mode.Width,
					height: mode.Height,
					stride: mode.BitsPerPixel / 8 * mode.Height,
				}
				// TODO check buffer alloc/free
				data := make([]byte, buffer.height*buffer.stride)
				buffer.buffer = &data[0]

				buffer.rectCount = MAX_DIRTS
				rectArr := make([]evdiRect, buffer.rectCount)
				buffer.rects = &rectArr[0]

				h.RegisterBuffer(buffer)
				b := bufferData{
					data:   data,
					buffer: buffer,
				}
				h.buffersMap[buffer.id] = b
			}
			h.bufferToUpdate = 0
			if h.RequestUpdate(h.bufferToUpdate) {
				b := h.buffersMap[h.bufferToUpdate]
				h.GrabPixels(b.buffer.rects, &b.buffer.rectCount)
			}
		},

		Dpms: func(dpmsMode DpmsMode, userData any) {
			evdiLogInfo("event: dpms: %v", dpmsMode)
		},
		UpdateReady: func(bufferToUpdate int32, userData any) {
			evdiLogInfo("event: updateReady: %d", bufferToUpdate)
			buff := h.buffersMap[h.bufferToUpdate]
			h.GrabPixels(buff.buffer.rects, &buff.buffer.rectCount)
			currentBuffer = (currentBuffer + 1) % int32(len(h.buffersMap))
			h.RequestUpdate(currentBuffer)
		},
		CrtcState: func(state int32, userData any) {
			evdiLogInfo("event: crtc: 0x%x", state)
		},
		CursorSet: func(cursorSet CursorSet, userData any) {
			evdiLogInfo("event: cursorSet: %+v", cursorSet)
		},
		CursorMove: func(cursorMove CursorMove, userData any) {
			evdiLogInfo("event: cursorMove: %+v", cursorMove)
		},
		DdcciData: func(data DdcciData, userData any) {
			evdiLogInfo("event: ddciData: %+v", data)
		},
	}

	evdiLogDebug("dummy: connecting")

	skuLimit := uint32(1280 * 800)
	h.Connect(EDIDv1_1280x800, skuLimit)
	evdiLogDebug("dummy: connected")
	defer h.Disconnect()

	newSelectable := h.GetEventReady()
	event := unix.PollFd{
		Events: unix.EPOLLIN,
		Fd:     int32(newSelectable),
	}
	events := []unix.PollFd{
		event,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// if len(h.buffersMap) > 0 {
		// 	updateReady := h.RequestUpdate(currentBuffer)
		// 	evdiLogDebug("dummy: requestUpdate buff=%d ready=%v", currentBuffer, updateReady)
		// 	if updateReady {
		// 		buf := h.buffersMap[currentBuffer]

		// 		h.GrabPixels(buf.buffer.rects, &buf.buffer.rectCount)
		// 		currentBuffer = (currentBuffer + 1) % int32(len(h.buffersMap))
		// 		continue
		// 	}
		// }

		n, err := unix.Poll(events, 1000)
		if err != nil {
			evdiLogDebug("epoll wait error: %v", err)
			if err == unix.EINTR {
				continue
			}
			continue
		}
		evdiLogDebug("polled n=%d", n)
		if n == 0 {
			continue
		}
		h.HandleEvents(eventsHandler)
	}
	return fmt.Errorf("TODO")
}

const (
	MAX_DIRTS = 16
)

func (h *Handle) GrabPixels(rects *evdiRect, numRects *int32) {
	evdiLogInfo("called grabPixels")
	destinationNode := h.frameBuffersListHead.FindItem(func(a *evdiBuffer) bool {
		return a.id == h.bufferToUpdate
	})
	if destinationNode == nil {
		evdiLogInfo("buffer %d not found. Not grabbing.", h.bufferToUpdate)
		*numRects = 0
		return
	}
	evdiLogDebug("grabpix: numRects=%v", numRects)
	destinationBuffer := destinationNode.Item

	kernelDirts := make([]drmClipRect, MAX_DIRTS)
	grab := drmEvdiGrabPix{
		mode:          EVDI_GRABPIX_MODE_DIRTY,
		bufWidth:      destinationBuffer.width,
		bufHeight:     destinationBuffer.height,
		bufByteStride: destinationBuffer.stride,
		buffer:        destinationBuffer.buffer,
		numRects:      MAX_DIRTS,
		rects:         &kernelDirts[0],
	}
	ret := doIoctl(h.fd, DRM_IOCTL_EVDI_GRABPIX, uintptr(unsafe.Pointer(&grab)), "grabpix")
	if ret == 0 {
		/*
		 * Buffer was filled by ioctl
		 * now we only have to fill the dirty rects
		 */

		rectsArr := unsafe.Slice(rects, grab.numRects)

		for r := 0; r < int(grab.numRects); r++ {
			rectsArr[r].x1 = int32(kernelDirts[r].x1)
			rectsArr[r].y1 = int32(kernelDirts[r].y1)
			rectsArr[r].x2 = int32(kernelDirts[r].x2)
			rectsArr[r].y2 = int32(kernelDirts[r].y2)
		}

		*numRects = grab.numRects
	} else {
		id := destinationBuffer.id

		evdiLogInfo("Grabbing pixels for buffer %d failed.", id)
		evdiLogInfo("Ignore if caused by change of mode.")
		*numRects = 0
	}
}
