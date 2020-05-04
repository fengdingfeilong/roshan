package roshan

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/fengdingfeilong/roshan/handler"
	"github.com/fengdingfeilong/roshan/message"
)

//tcpSCBase tcpserver and tcpclient base
type tcpSCBase struct {
	handlerManager *handler.Manager
	//CmdMessageReceived received command message
	CmdMessageReceived func(conn net.Conn, t message.CmdType)
	//HBSendInterval interval for sending heartbeat(unit:second)
	HBSendInterval int
	//HBTimeout timeout for receiving heartbeat(unit:second)
	HBTimeout int
	//socket disconnect
	SocketDisconnect func(conn net.Conn)
	//handle ErrorOccurred error occurred
	ErrorOccurred func(conn net.Conn, err error)
	//security key(use aes, ctr mode)
	key    string
	cmutex sync.Mutex
}

//set security key(use aes, ctr mode)
func (sc *tcpSCBase) SetSK(key string) {
	sc.key = key
}

//AddHandler add handler for message
func (sc *tcpSCBase) AddHandler(msgType message.CmdType, handler handler.Handler) {
	sc.handlerManager.Add(msgType, handler)
}

func (sc *tcpSCBase) handleConn(cc *connContext) {
	sc.handlerManager.Foreach(func(t message.CmdType, h handler.Handler) {
		h.GetBase().SetConn(cc)
	})
	go sc.sendHB(cc)
	go sc.checkConnection(cc)
	for {
		r, err := message.ParsePacket(cc, message.ParseCallback(sc.handlePacket))
		if !r {
			if err != io.EOF {
				loginfo(fmt.Sprintf("tcpscbase handle conn parse error: %s", err.Error()), err)
				if sc.ErrorOccurred != nil {
					sc.ErrorOccurred(cc, err)
				}
				// sc.closeSocket(cc)
			}
			return
		}
	}
}

func (sc *tcpSCBase) handlePacket(conn net.Conn, pac *message.Packet) {
	cc, ok := conn.(*connContext)
	if !ok {
		return
	}
	switch pac.Type {
	case message.HeartBeat:
		return
	case message.Command:
		{
			t := binary.BigEndian.Uint32(pac.Payload[:4])
			mt := message.CmdType(t)
			if sc.CmdMessageReceived != nil {
				sc.CmdMessageReceived(cc, mt)
			}
			h := sc.handlerManager.Get(mt)
			if h == nil {
				loginfo("can not find the handler, message type : "+strconv.Itoa(int(t)), nil)
				return
			}
			if cc.cancelHandle {
				return
			}
			h.Execute(pac.Payload[4:])
		}
	case message.Data:
		h := sc.handlerManager.Get(message.CmdType(0))
		if h == nil {
			loginfo("can not find the Data handler ", nil)
			return
		}
		if cc.cancelHandle {
			return
		}
		h.Execute(pac.Payload)
	default:
		return
	}

}

func (sc *tcpSCBase) sendHB(cc *connContext) {
	for {
		buff := message.GetHeartBeatBytes()
		_, err := cc.Write(buff)
		if err != nil {
			break
		}
		for {
			time.Sleep(time.Second)
			if int(time.Now().Sub(cc.lastSendTime).Seconds()) > sc.HBSendInterval {
				break
			}
		}
	}
}

func (sc *tcpSCBase) checkConnection(cc *connContext) {
	for {
		time.Sleep(time.Second)
		if int(time.Now().Sub(cc.lastReceiveTime).Seconds()) > sc.HBTimeout {
			sc.closeSocket(cc)
			break
		}
	}
}

func (sc *tcpSCBase) closeSocket(cc *connContext) {
	sc.cmutex.Lock()
	defer sc.cmutex.Unlock()
	if cc.isClosed {
		return
	}
	cc.Close()
	sc.handlerManager.Foreach(func(t message.CmdType, h handler.Handler) {
		b := h.GetBase()
		if b != nil {
			b.Dispose()
		}
	})
	if sc.SocketDisconnect != nil {
		sc.SocketDisconnect(cc)
	}
}

//Transmit send message to handler
func (sc *tcpSCBase) Transmit(para *handler.CommObj) {
	h := sc.handlerManager.Get(para.Dst)
	if h != nil {
		h.Receive(para)
	}
}
