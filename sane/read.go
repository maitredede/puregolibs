package sane

import (
	"io"

	"github.com/maitredede/puregolibs/sane/internal"
)

func (h *Handle) Read(buf []byte) (int, error) {

	nBuf := (*byte)(&buf[0])
	nMaxLen := internal.SANE_Int(len(buf))
	var nLen internal.SANE_Int

	ret := libSaneRead(h.h, nBuf, nMaxLen, &nLen)
	switch ret {
	case StatusGood:
		return int(nLen), nil
	case StatusEOF:
		return int(nLen), io.EOF
	}

	return 0, mkError(ret)
}
