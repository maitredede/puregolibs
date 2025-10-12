package plutobook

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type streamWriteData struct {
	output io.Writer
	err    error
}

type streamReaderData struct {
	input io.Reader
}

type StreamStatus int32

const (
	StreamStatusSuccess    StreamStatus = 0
	StreamStatusReadError  StreamStatus = 10
	StreamStatusWriteError StreamStatus = 11
)

func streamWriteCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	argArr := unsafe.Slice(args, cif.NArgs)
	closureArgPtr := argArr[0]
	dataArgPtr := argArr[1]
	lgArgPtr := argArr[2]

	closurePtr := *(*unsafe.Pointer)(closureArgPtr)
	dataPtr := *(*unsafe.Pointer)(dataArgPtr)

	lg := *(*uint32)(lgArgPtr)
	stream := (*streamWriteData)(closurePtr)
	data := unsafe.Slice((*byte)(dataPtr), lg)

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

func streamReadCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	argArr := unsafe.Slice(args, cif.NArgs)
	closureArgPtr := argArr[0]
	dataArgPtr := argArr[1]
	lgArgPtr := argArr[2]

	closurePtr := *(*unsafe.Pointer)(closureArgPtr)
	dataPtr := *(*unsafe.Pointer)(dataArgPtr)

	lg := *(*uint32)(lgArgPtr)
	stream := (*streamReaderData)(closurePtr)
	// data := unsafe.Slice((*byte)(dataPtr), lg)

	_ = dataPtr
	_ = lg
	_ = stream

	//do the thing
	*(*StreamStatus)(ret) = StreamStatusWriteError
	// n, err := stream.output.Write(data)
	// if err != nil {
	// 	stream.err = err
	// 	return 0
	// }
	// if uint32(n) != lg {
	// 	stream.err = fmt.Errorf("data write mismatch: length=%v written=%v", lg, n)
	// 	return 0
	// }
	// *(*StreamStatus)(ret) = StreamStatusSuccess
	return 0
}
