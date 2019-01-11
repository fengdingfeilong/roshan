package handler

import (
	"net"
	"roshan/message"
	"sync"
)

//Handler process package
type Handler interface {
	//GetBase get base type
	GetBase() *Base
	//Execute process package
	Execute(buff []byte)
	//Receive receive message from other handler
	Receive(para *CommObj)
}

//Base handler base
type Base struct {
	conn net.Conn
}

//GetBase handler base virtual function of GetBase
func (h *Base) GetBase() *Base {
	return nil
}

//Execute handler base virtual function of Execute
func (h *Base) Execute(buff []byte) {
}

//Receive handler base virtual function of Receive
func (h *Base) Receive(para *CommObj) {
}

//Conn get conn
func (h *Base) Conn() net.Conn {
	return h.conn
}

//SetConn set conn
func (h *Base) SetConn(c net.Conn) {
	h.conn = c
}

//Dispose clean your resource if need
func (h *Base) Dispose() {

}

//CommObj use for communication of handlers
type CommObj struct {
	Src message.CmdType
	Dst message.CmdType
	Obj []interface{}
}

//NewCommObj ...
func NewCommObj(s message.CmdType, d message.CmdType, o ...interface{}) *CommObj {
	return &CommObj{Src: s, Dst: d, Obj: o}
}

//Manager ...
type Manager struct {
	sync.Mutex
	handlers map[message.CmdType]Handler
}

//Add add handlers, when msgType is 0, it means this is a data packet
func (manager *Manager) Add(msgType message.CmdType, handler Handler) {
	manager.Lock()
	defer manager.Unlock()
	if manager.handlers == nil {
		manager.handlers = make(map[message.CmdType]Handler)
	}
	manager.handlers[msgType] = handler
}

//Get get the handler
func (manager *Manager) Get(msgType message.CmdType) Handler {
	manager.Lock()
	defer manager.Unlock()
	return manager.handlers[msgType]
}

//Foreach traveral handlers
func (manager *Manager) Foreach(f func(message.CmdType, Handler)) {
	manager.Lock()
	defer manager.Unlock()
	for t, h := range manager.handlers {
		f(t, h)
	}
}
