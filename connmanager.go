package roshan

import (
	"net"
	"sync"
)

type connManager struct {
	connections map[net.Conn]*connContext
	mutex       sync.Mutex
}

func newConnManager() *connManager {
	var manager connManager
	manager.connections = make(map[net.Conn]*connContext)
	return &manager
}

func (manager *connManager) Add(cContext *connContext) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.connections[cContext.conn] = cContext
}

func (manager *connManager) Remove(cContext *connContext) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.connections, cContext.conn)
}

func (manager *connManager) Get(conn net.Conn) *connContext {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	return manager.connections[conn]
}
