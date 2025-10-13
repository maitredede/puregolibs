package gousb

import "context"

type endpoint struct {
	h libusbDeviceHandle

	InterfaceSetting
	Desc EndpointDesc

	ctx *Context
}

// String returns a human-readable description of the endpoint.
func (e *endpoint) String() string {
	return e.Desc.String()
}

func (e *endpoint) transfer(ctx context.Context, buf []byte) (int, error) {
	// t, err := newUSBTransfer(e.ctx, e.h, &e.Desc, len(buf))
	// if err != nil {
	// 	return 0, err
	// }
	// defer t.free()
	// if e.Desc.Direction == EndpointDirectionOut {
	// 	copy(t.data(), buf)
	// }

	// if err := t.submit(); err != nil {
	// 	return 0, err
	// }

	// n, err := t.wait(ctx)
	// if e.Desc.Direction == EndpointDirectionIn {
	// 	copy(buf, t.data())
	// }
	// if err != nil {
	// 	return n, err
	// }
	// return n, nil
	panic("WIP")
}

// InEndpoint represents an IN endpoint open for transfer.
// InEndpoint implements the io.Reader interface.
// For high-throughput transfers, consider creating a buffered read stream
// through InEndpoint.ReadStream.
type InEndpoint struct {
	*endpoint
}

// Read reads data from an IN endpoint. Read returns number of bytes obtained
// from the endpoint. Read may return non-zero length even if
// the returned error is not nil (partial read).
// It's recommended to use buffer sizes that are multiples of
// EndpointDesc.MaxPacketSize to avoid overflows.
// When a USB device receives a read request, it doesn't know the size of the
// buffer and may send too much data in one packet to fit in the buffer.
// If that happens, Read will return an error signaling an overflow.
// See http://libusb.sourceforge.net/api-1.0/libusb_packetoverflow.html
// for more details.
func (e *InEndpoint) Read(buf []byte) (int, error) {
	return e.transfer(context.Background(), buf)
}

// ReadContext reads data from an IN endpoint. ReadContext returns number of
// bytes obtained from the endpoint. ReadContext may return non-zero length
// even if the returned error is not nil (partial read).
// The passed context can be used to control the cancellation of the read. If
// the context is cancelled, ReadContext will cancel the underlying transfers,
// resulting in TransferCancelled error.
// It's recommended to use buffer sizes that are multiples of
// EndpointDesc.MaxPacketSize to avoid overflows.
// When a USB device receives a read request, it doesn't know the size of the
// buffer and may send too much data in one packet to fit in the buffer.
// If that happens, Read will return an error signaling an overflow.
// See http://libusb.sourceforge.net/api-1.0/libusb_packetoverflow.html
// for more details.
func (e *InEndpoint) ReadContext(ctx context.Context, buf []byte) (int, error) {
	return e.transfer(ctx, buf)
}

// OutEndpoint represents an OUT endpoint open for transfer.
type OutEndpoint struct {
	*endpoint
}

// Write writes data to an OUT endpoint. Write returns number of bytes comitted
// to the endpoint. Write may return non-zero length even if the returned error
// is not nil (partial write).
func (e *OutEndpoint) Write(buf []byte) (int, error) {
	return e.transfer(context.Background(), buf)
}

// WriteContext writes data to an OUT endpoint. WriteContext returns number of
// bytes comitted to the endpoint. WriteContext may return non-zero length even
// if the returned error is not nil (partial write).
// The passed context can be used to control the cancellation of the write. If
// the context is cancelled, WriteContext will cancel the underlying transfers,
// resulting in TransferCancelled error.
func (e *OutEndpoint) WriteContext(ctx context.Context, buf []byte) (int, error) {
	return e.transfer(ctx, buf)
}
