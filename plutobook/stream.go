package plutobook

import (
	"fmt"
	"io"
	"log/slog"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type streamWriteData struct {
	output io.Writer
	err    error
}

type StreamStatus int32

const (
	StreamStatusSuccess    StreamStatus = 0
	StreamStatusReadError  StreamStatus = 10
	StreamStatusWriteError StreamStatus = 11
)

func streamWriteCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	argArr := unsafe.Slice(args, cif.NArgs)
	closurePtr := argArr[0]
	dataPtr := argArr[1]
	lgPtr := argArr[2]

	lg := *(*uint32)(lgPtr)
	stream := *(**streamWriteData)(closurePtr)
	data := unsafe.Slice((*byte)(dataPtr), lg)

	slog.Warn(fmt.Sprintf("streamWrite: lg=%v", lg))

	for i := uint32(0); i < lg; i++ {
		b := data[i]

		_ = b
	}

	//do the thing
	*(*StreamStatus)(ret) = StreamStatusWriteError
	n, err := stream.output.Write(data)
	if err != nil {
		stream.err = err
		return 0
	}
	if uint32(n) != lg {
		stream.err = fmt.Errorf("data write mismatch: length=%v written=%v", lg, n)
		return 0
	}
	*(*StreamStatus)(ret) = StreamStatusSuccess
	return 0
}
