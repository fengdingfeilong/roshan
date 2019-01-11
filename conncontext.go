package roshan

import (
	"net"
	"sync"
	"time"
)

type connContext struct {
	conn             net.Conn
	cancelHandle     bool
	lastReceiveTime  time.Time
	lastSendTime     time.Time
	isClosed         bool
	socketErrOccured func(*connContext, error)
	mutex            sync.Mutex
}

func newConnContext(conn net.Conn) *connContext {
	var cc connContext
	cc.conn = conn
	cc.cancelHandle = false
	cc.lastReceiveTime = time.Now()
	cc.lastSendTime = time.Now()
	cc.isClosed = false
	return &cc
}

func (cc *connContext) Read(b []byte) (int, error) {
	n, err := cc.conn.Read(b)
	cc.lastReceiveTime = time.Now()
	if err != nil && cc.socketErrOccured != nil {
		cc.socketErrOccured(cc, err)
	}
	return n, err
}

func (cc *connContext) Write(b []byte) (int, error) {
	n, err := cc.conn.Write(b)
	cc.lastSendTime = time.Now()
	if err != nil && cc.socketErrOccured != nil {
		cc.socketErrOccured(cc, err)
	}
	return n, err
}

func (cc *connContext) Close() error {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	if cc.isClosed {
		return nil
	}
	err := cc.conn.Close()
	if err == nil {
		cc.isClosed = true
	}
	return err
}

func (cc *connContext) LocalAddr() net.Addr {
	return cc.conn.LocalAddr()
}

func (cc *connContext) RemoteAddr() net.Addr {
	return cc.conn.RemoteAddr()
}

func (cc *connContext) SetDeadline(t time.Time) error {
	return cc.conn.SetDeadline(t)
}

func (cc *connContext) SetReadDeadline(t time.Time) error {
	return cc.conn.SetReadDeadline(t)
}

func (cc *connContext) SetWriteDeadline(t time.Time) error {
	return cc.conn.SetWriteDeadline(t)
}
