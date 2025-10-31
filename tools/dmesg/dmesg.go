package dmesg

// import (
// 	"bytes"
// 	"errors"
// 	"os"
// 	"strconv"
// 	"syscall"
// )

// const (
// 	defaultBufSize = uint32(1 << 14) // 16KB by default
// 	levelMask      = uint64(1<<3 - 1)
// )

// type dmesg struct {
// 	msg []Msg
// }

// type Msg struct {
// 	Level      uint64            // SYSLOG lvel
// 	Facility   uint64            // SYSLOG facility
// 	Seq        uint64            // Message sequence number
// 	TsUsec     int64             // Timestamp in microsecond
// 	Caller     string            // Message caller
// 	IsFragment bool              // This message is a fragment of an early message which is not a fragment
// 	Text       string            // Log text
// 	DeviceInfo map[string]string // Device info
// }

// func DumpTail() {
// 	d := dmesg{}
// 	file, err := os.OpenFile("/dev/kmsg", syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	var conn syscall.RawConn
// 	conn, err = file.SyscallConn()
// 	if err != nil {
// 		panic(err)
// 	}

// 	d.msg = make([]Msg, 0)

// 	var syscallError error = nil
// 	err = conn.Read(func(fd uintptr) bool {
// 		for {
// 			buf := make([]byte, bufSize)
// 			_, err := syscall.Read(int(fd), buf)
// 			if err != nil {
// 				syscallError = err
// 				// EINVAL means buf is not enough, data would be truncated, but still can continue.
// 				if !errors.Is(err, syscall.EINVAL) {
// 					return true
// 				}
// 			}

// 			msg := parseData(buf)
// 			if msg == nil {
// 				continue
// 			}
// 			d.msg = append(d.msg, *msg)
// 		}
// 	})

// 	// EAGAIN means no more data, should be treated as normal.
// 	if syscallError != nil && !errors.Is(syscallError, syscall.EAGAIN) {
// 		err = syscallError
// 	}
// }

// func parseData(data []byte) *Msg {
// 	msg := Msg{}

// 	dataLen := len(data)
// 	prefixEnd := bytes.IndexByte(data, ';')
// 	if prefixEnd == -1 {
// 		return nil
// 	}

// 	for index, prefix := range bytes.Split(data[:prefixEnd], []byte(",")) {
// 		switch index {
// 		case 0:
// 			val, _ := strconv.ParseUint(string(prefix), 10, 64)
// 			msg.Level = val & levelMask
// 			msg.Facility = val & (^levelMask)
// 		case 1:
// 			val, _ := strconv.ParseUint(string(prefix), 10, 64)
// 			msg.Seq = val
// 		case 2:
// 			val, _ := strconv.ParseInt(string(prefix), 10, 64)
// 			msg.TsUsec = val
// 		case 3:
// 			msg.IsFragment = prefix[0] != '-'
// 		case 4:
// 			msg.Caller = string(prefix)
// 		}
// 	}

// 	textEnd := bytes.IndexByte(data, '\n')
// 	if textEnd == -1 || textEnd <= prefixEnd {
// 		return nil
// 	}

// 	msg.Text = string(data[prefixEnd+1 : textEnd])
// 	if textEnd == dataLen-1 {
// 		return nil
// 	}

// 	msg.DeviceInfo = make(map[string]string, 2)
// 	deviceInfo := bytes.Split(data[textEnd+1:dataLen-1], []byte("\n"))
// 	for _, info := range deviceInfo {
// 		if info[0] != ' ' {
// 			continue
// 		}

// 		kv := bytes.Split(info, []byte("="))
// 		if len(kv) != 2 {
// 			continue
// 		}

// 		msg.DeviceInfo[string(kv[0])] = string(kv[1])
// 	}

// 	return &msg
// }
