package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/nitingwoshuoa/shv-zinx_tcp/ziface"
)

/*
	链接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的链接集合
	connLock    sync.RWMutex                  // 保护链接集合的读写锁
}

// 创建链接的方法

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//将conn加入到connManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connecion add to connmanager succ : conn num = ", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connecion remove from connmanager succ : conn num = ", connMgr.Len())
}

// 根据id 获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not Found")
	}
}

// 得到链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("clear all connections succ")
}
