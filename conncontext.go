package roshan

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
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
	sk               string
	rmutex           sync.Mutex
	wmutex           sync.Mutex
	cmutex           sync.Mutex
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

func (cc *connContext) syncRead(buf []byte) (int, error) {
	l := len(buf)
	var t int
	for {
		r, err := cc.conn.Read(buf[t:])
		if err != nil {
			return r, err
		}
		t += r
		if t >= l {
			return t, nil
		}
	}
}

func (cc *connContext) Read(b []byte) (int, error) {
	cc.rmutex.Lock()
	defer cc.rmutex.Unlock()
	db := make([]byte, len(b))
	r, err := cc.syncRead(db)
	if err != nil {
		if err != nil && cc.socketErrOccured != nil {
			cc.socketErrOccured(cc, err)
		}
		return r, err
	}
	tb := cc.aesCtrCrypt(db)
	copy(b, tb)
	cc.lastReceiveTime = time.Now()
	return r, err
}

func (cc *connContext) Write(b []byte) (int, error) {
	cc.wmutex.Lock()
	defer cc.wmutex.Unlock()
	eb := cc.aesCtrCrypt(b[:2]) //encrypt header two bytes
	copy(b[:2], eb)
	if len(b) > 2 {
		eb = cc.aesCtrCrypt(b[2:6]) //encrypt length
		copy(b[2:6], eb)
		if len(b) > 6 {
			eb = cc.aesCtrCrypt(b[6:])
			copy(b[6:], eb)
		}
	}
	n, err := cc.conn.Write(b)
	cc.lastSendTime = time.Now()
	if err != nil && cc.socketErrOccured != nil {
		cc.socketErrOccured(cc, err)
	}
	return n, err
}

func (cc *connContext) Close() error {
	cc.cmutex.Lock()
	defer cc.cmutex.Unlock()
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

func (cc *connContext) SetSK(key string) {
	cc.sk = key
}

//encrypt and decrypt data (use aes,ctr mode)
func (cc *connContext) aesCtrCrypt(data []byte) []byte {
	h := sha256.New()
	h.Write([]byte(cc.sk))
	key := h.Sum(nil)
	h.Reset()
	h.Write([]byte(cc.sk))
	iv := h.Sum(nil)
	c, _ := aes.NewCipher(key)
	dec := cipher.NewCTR(c, iv[:c.BlockSize()])
	out := make([]byte, len(data))
	dec.XORKeyStream(out, data)
	return out
}
