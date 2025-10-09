package cec

const (
	// Maximum size of a data packet
	CEC_MAX_DATA_PACKET_SIZE = 16 * 4
)

type nativeDataPacket struct {
	data [CEC_MAX_DATA_PACKET_SIZE]uint8
	size uint8
}

func (n nativeDataPacket) Slice() []byte {
	return n.data[:n.size]
}
