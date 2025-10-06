package sane

import "io"

func Read(h SANE_Handle, buf []byte) (int, error) {

	nBuf := (*byte)(&buf[0])
	nMaxLen := SANE_Int(len(buf))
	var nLen SANE_Int

	ret := libSaneRead(h, nBuf, nMaxLen, &nLen)
	switch ret {
	case StatusGood:
		return int(nLen), nil
	case StatusEOF:
		return int(nLen), io.EOF
	}

	return 0, mkError(ret)
}
