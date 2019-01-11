package roshan

import (
	"fmt"
	"net"
	"strconv"

	"github.com/fengdingfeilong/roshan/handler"
)

//Client tcpserver client
type Client struct {
	tcpSCBase
	//SocketConnected socket connected
	SocketConnected func(conn net.Conn)
	//BeforeClose before close the socket
	BeforeClose func(conn net.Conn)

	workConn net.Conn
}

//NewClient create client
func NewClient() *Client {
	var c Client
	c.handlerManager = &handler.Manager{}
	c.HBSendInterval = 5
	c.HBTimeout = 10
	return &c
}

//Connect connect to server
func (client *Client) Connect(ip string, port int) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		loginfo(err.Error(), err)
		return
	}
	cc := newConnContext(conn)
	cc.socketErrOccured = client.handleSocErr
	client.workConn = cc
	if client.SocketConnected != nil {
		client.SocketConnected(cc)
	}
	go client.handleConn(cc)
}

func (client *Client) handleSocErr(cc *connContext, err error) {
	loginfo(fmt.Sprintf("socket error: %s", err.Error()), err)
	client.closeSocket(cc)
}

// func (client *Client) closeSocket(cc *connContext) {
// 	cc.Close()
// }

//Close close the connection
func (client *Client) Close() {
	if client.BeforeClose != nil {
		temp := client.workConn
		client.BeforeClose(temp)
	}
	client.workConn.Close()
	client.workConn = nil
}
