package roshan

import (
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/fengdingfeilong/roshan/handler"
	"github.com/fengdingfeilong/roshan/roshantool"
)

func loginfo(s string, err error) {
	if roshantool.InnerLog != nil {
		roshantool.InnerLog("roshan: "+s, err)
	}
}

//Server tcp server struct
type Server struct {
	tcpSCBase
	listener *net.TCPListener
	//BeforeAccept before accept socket(this is not the same with ConnectMessage), you can do something such as stop or block continue accept
	BeforeAccept func()
	//SocketAccepted accept socket
	SocketAccepted func(conn net.Conn)

	stopacc bool
	//connmanager *connManager
}

//NewServer create new tcp server
func NewServer() *Server {
	var s Server
	s.handlerManager = &handler.Manager{}
	s.stopacc = false
	//s.connmanager = newConnManager()
	s.HBSendInterval = 5
	s.HBTimeout = 10
	return &s
}

//Start start tcp server
func (server *Server) Start(port int) {
	r := server.startSever(port)
	if !r {
		loginfo("start server failed", nil)
		return
	}
	loginfo(fmt.Sprintf("start server success. listen port:%d", port), nil)
	server.startAccept()
}

//CloseListen close listen
func (server *Server) CloseListen() {
	server.listener.Close()
	server.listener = nil
}

//StopAccept stop continue accept socket
func (server *Server) StopAccept() {
	server.stopacc = true
}

func (server *Server) startSever(port int) bool {
	var err error
	server.listener, err = net.ListenTCP("tcp4", &net.TCPAddr{IP: nil, Port: port})
	if err != nil {
		loginfo(err.Error(), err)
		return false
	}
	return true
}

func (server *Server) startAccept() {
	for {
		if server.BeforeAccept != nil {
			server.BeforeAccept()
		}
		if server.stopacc {
			break
		}
		conn, err := server.listener.Accept()
		if err != nil {
			if err == syscall.EINVAL {
				loginfo("listen closed", err)
				break
			}
			loginfo(err.Error(), err)
			time.Sleep(time.Second * 3)
			continue
		}
		cc := newConnContext(conn)
		cc.socketErrOccured = server.handleSocErr
		//server.connmanager.Add(context)
		if server.SocketAccepted != nil {
			server.SocketAccepted(cc)
		}
		go server.handleConn(cc)
	}
}

func (server *Server) handleSocErr(cc *connContext, err error) {
	loginfo(fmt.Sprintf("socket error: %s", err.Error()), err)
	server.closeSocket(cc)
}

// func (server *Server) closeSocket(cc *connContext) {
// 	cc.Close()
// 	//server.connmanager.Remove(cc.conn)
// }

//StartHandlePacket continue to handle packet
//after receive some special command packet such as connectmessage, you can start or stop handle command and data packet
func (server *Server) StartHandlePacket(conn net.Conn) {
	cc, ok := conn.(*connContext)
	if ok {
		cc.cancelHandle = false
	}
}

//StopHandlePacket stop handle command and data packet
func (server *Server) StopHandlePacket(conn net.Conn) {
	cc, ok := conn.(*connContext)
	if ok {
		cc.cancelHandle = true
	}
}
