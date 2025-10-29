package evdi

import (
	"fmt"
	"log/slog"
	"unsafe"

	"github.com/maitredede/puregolibs/libc"
	"golang.org/x/sys/unix"
)

type Card struct {
	dispose func()

	evdiHandle   evdiHandle
	eventContext *evdiEventContext
	mode         EvdiMode

	buffers         []*Buffer
	bufferRequested *Buffer

	acquireFramebufferHandler func(buffer *Buffer)
}

func NewCard(device int32) (*Card, error) {
	h := libEvdiOpen(device)
	if h == nil {
		return nil, fmt.Errorf("failed topp open card '/dev/dri/card%d'", device)
	}

	c := &Card{}
	c.evdiHandle = h

	e := Events{
		ModeChanged: cardCModeHandler,
		UpdateReady: defaultUpdateReadyHandler,
		CursorSet:   cardCCursorSetHandler,
		CursorMove:  cardCCursorMoveHandler,
		Dpms:        dpmsHandler,
	}

	c.eventContext, c.dispose = buildNativeEvents(e, unsafe.Pointer(c))
	return c, nil
}

func (c *Card) Close() error {
	if c.evdiHandle != nil {
		c.clearBuffers()
		libEvdiClose(c.evdiHandle)
	}
	c.evdiHandle = nil
	return nil
}

func (c *Card) Connect(edid []byte, pixelAreaLimit uint32, pixelPerSecondLimit uint32) {
	var edidPtr *byte
	var edidLen uint32
	if len(edid) > 0 {
		edidPtr = &edid[0]
		edidLen = uint32(len(edid))
	}
	libEvdiConnect2(c.evdiHandle, edidPtr, edidLen, pixelAreaLimit, pixelPerSecondLimit)
}

func (c *Card) Disconnect() {
	libEvdiDisconnect(c.evdiHandle)
}

func (c *Card) GetMode() EvdiMode {
	return c.mode
}

func (c *Card) HandleEvents(waitingTime int) {
	fd := libEvdiGetEventReady(c.evdiHandle)
	//TODO double check

	fds := []unix.PollFd{
		{Fd: int32(fd), Events: unix.POLLIN},
	}

	c.requestUpdate()

	n, err := unix.Poll(fds, waitingTime)
	if err != nil {
		panic(err)
	}
	if n > 0 {
		libEvdiHandleEvents(c.evdiHandle, c.eventContext)
	}
}

func (c *Card) requestUpdate() {
	if c.bufferRequested != nil {
		return
	}

	// for _, i := range c.buffers {
	// 	if i.useCount == 1 {
	// 		c.bufferRequested = i
	// 		break
	// 	}
	// }

	if c.bufferRequested == nil {
		return
	}

	updateReady := libevdiRequestUpdate(c.evdiHandle, c.bufferRequested.buffer.id)
	if updateReady {
		c.grabPixels()
	}
}

func (c *Card) grabPixels() {
	if c.bufferRequested == nil {
		return
	}

	//mStat.grabPixels(c.evdiHandle, c.bufferRequested.buffer.rects, &c.bufferRequested.buffer.rectCount)
	libevdiGrabPixels(c.evdiHandle, &c.bufferRequested.buffer.rects, &c.bufferRequested.buffer.rectCount)

	if c.acquireFramebufferHandler != nil {
		c.acquireFramebufferHandler(c.bufferRequested)
	}
	c.bufferRequested = nil

	c.requestUpdate()
}

func (c *Card) EnableCursorEvents(enable bool) {
	libEvdiEnableCursorEvents(c.evdiHandle, enable)
}

func (c *Card) clearBuffers() {
	// c.bufferRequested.reset()
	c.bufferRequested = nil
	//c.buffers.clear()
	for _, b := range c.buffers {
		b.Close()
	}
	c.buffers = nil
}

func (c *Card) makeBuffers(count int) {
	c.clearBuffers()
	c.buffers = make([]*Buffer, count)
	for i := 0; i < count; i++ {
		b := NewBuffer(c.mode, c.evdiHandle)
		c.buffers[i] = b
	}
}

func (c *Card) setMode(mode EvdiMode) {
	c.mode = mode
}

func defaultUpdateReadyHandler(bufferToUpdate int32, userData unsafe.Pointer) {
	card := (*Card)(userData)
	if card.bufferRequested.buffer.id != bufferToUpdate {
		panic("wrong buffer")
	}
	card.grabPixels()
}

func cardCModeHandler(mode EvdiMode, userdata unsafe.Pointer) {
	slog.Info("got mode changed")
	card := (*Card)(userdata)

	card.setMode(mode)
	card.makeBuffers(2)

	// if card.modeHandler != nil {
	// 	card.modeHandler(mode)
	// }

	card.requestUpdate()
}

func cardCCursorSetHandler(cursorSet CursorSet, userData unsafe.Pointer) {
	slog.Info("got cursor set event")
	card := (*Card)(userData)
	_ = card
	libc.Free(unsafe.Pointer(cursorSet.Buffer))
}

func cardCCursorMoveHandler(cursorMove CursorMove, userData unsafe.Pointer) {
	slog.Info("got cursor move event")
	card := (*Card)(userData)
	_ = card
}

func dpmsHandler(dpmsMode int32, userData unsafe.Pointer) {
	// card := (*Card)(userData)
	slog.Info(fmt.Sprintf("got dpms signal %v", dpmsMode))
}
