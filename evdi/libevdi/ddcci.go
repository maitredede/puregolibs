package libevdi

type DdcciData struct {
	Address uint16
	Flags   uint16
	// bufferLength uint32
	// buffer       *byte
	BufferData []byte
}

func toEvdiDdcciData(event drmEvdiEventDdcciData) DdcciData {
	data := DdcciData{
		Address:    event.address,
		Flags:      event.flags,
		BufferData: event.buffer[:event.bufferLength],
	}
	return data
}
