package message

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

//Packet is the tcp packet data format
//| description | version | type | length | payload |
//| size(byte)  |    1    |  1   |   4    | length  |
//when the type is heartbeat(0x00), there is no Length and Payload field
type Packet struct {
	Version byte
	Type    PacketType
	Length  int32
	Payload []byte
}

//PacketType raw packet type
type PacketType byte

//HeaderLen packet header length
const HeaderLen int32 = 6

//PacketVersion the default version of packet
const PacketVersion byte = 0x2B

const (
	HeartBeat = PacketType(iota)
	Command
	Data
)

func read(reader io.Reader, buf []byte) (bool, error) {
	l := len(buf)
	var t int
	for {
		r, err := reader.Read(buf[t:])
		if err != nil {
			return false, err
		}
		t += r
		if t >= l {
			return true, nil
		}
	}
}

type ParseCallback func(net.Conn, *Packet)

func ParsePacket(cc net.Conn, callback ParseCallback) (bool, error) {
	head := [HeaderLen]byte{}
	b, err := read(cc, head[:2])
	if !b && err != nil {
		return false, err
	}
	var pac Packet
	pac.Version = head[0]
	pac.Type = PacketType(head[1])
	if pac.Type != HeartBeat {
		b, err = read(cc, head[2:])
		if !b && err != nil {
			return false, err
		}
		pac.Length = int32(binary.BigEndian.Uint32(head[2:]))
		pac.Payload = make([]byte, pac.Length)
		b, err = read(cc, pac.Payload)
		if !b && err != nil {
			return false, err
		}
	}
	if callback != nil {
		callback(cc, &pac)
	}
	return true, nil
}

//GetBytes get the bytes of packet
func (pac *Packet) getBytes() []byte {
	buflen := pac.Length + HeaderLen
	if pac.Type == HeartBeat {
		buflen = 2
	}
	buf := make([]byte, buflen)
	buf[0] = pac.Version
	buf[1] = byte(pac.Type)
	if pac.Type != HeartBeat {
		if pac.Length > 0 {
			binary.BigEndian.PutUint32(buf[2:6], uint32(pac.Length))
			copy(buf[6:], pac.Payload)
		} else {
			panic("packet length should big than 0 if the type is not heartbeat")
		}
	}
	return buf
}

//GetHeartBeatBytes ...
func GetHeartBeatBytes() []byte {
	buflen := 2
	buf := make([]byte, buflen)
	buf[0] = PacketVersion
	buf[1] = byte(HeartBeat)
	return buf
}

//GetCommandBytes ...
func GetCommandBytes(t CmdType, o interface{}) []byte {
	var data []byte
	if s, ok := o.(string); ok {
		data = []byte(s)
	} else {
		data, _ = json.Marshal(o)
	}
	payloadLen := 4 + int32(len(data))
	buflen := HeaderLen + payloadLen
	buf := make([]byte, buflen)
	buf[0] = PacketVersion
	buf[1] = byte(Command)
	binary.BigEndian.PutUint32(buf[2:6], uint32(payloadLen))
	binary.BigEndian.PutUint32(buf[6:10], uint32(t))
	copy(buf[10:], data)
	return buf
}

//GetDataMsgBytes ...
func GetDataMsgBytes(id []byte, offset int64, data []byte) []byte {
	payloadLen := 16 + 8 + int32(len(data))
	buflen := HeaderLen + payloadLen
	buf := make([]byte, buflen)
	buf[0] = PacketVersion
	buf[1] = byte(Data)
	binary.BigEndian.PutUint32(buf[2:6], uint32(payloadLen))
	copy(buf[6:22], id)
	binary.BigEndian.PutUint64(buf[22:30], uint64(offset))
	copy(buf[30:], data)
	return buf
}
