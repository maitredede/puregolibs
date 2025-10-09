package cec

import "time"

type nativeCommand struct {
	initiator       LogicalAddress
	destination     LogicalAddress
	ack             int8
	eom             int8
	opcode          OpCode
	parameters      nativeDataPacket
	opcodeSet       int8
	transmitTimeout int32 //millis
}

func (n nativeCommand) Go() Command {
	c := Command{
		Initiator:       n.initiator,
		Destination:     n.destination,
		Ack:             n.ack != 0,
		Eom:             n.eom != 0,
		OpCode:          n.opcode,
		Parameters:      n.parameters.Slice(),
		OpCodeSet:       n.opcodeSet != 0,
		TransmitTimeout: time.Duration(n.transmitTimeout) * time.Millisecond,
	}
	return c
}

type Command struct {
	Initiator       LogicalAddress
	Destination     LogicalAddress
	Ack             bool
	Eom             bool
	OpCode          OpCode
	Parameters      []byte
	OpCodeSet       bool
	TransmitTimeout time.Duration
}
