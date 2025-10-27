package evdi

type evdiDdcciData struct {
	address      uint16
	flags        uint16
	bufferLength uint32
	buffer       *byte
}

type DdcciData struct {
	Address uint16
	Flags   uint16
	Buffer  []byte
}
