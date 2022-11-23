package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	//获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前链接模块的链接ID
	GetConnID() uint32
	// 获取远程客户端的TCP状态  IP port
	RemoteAddr() net.Addr
	//发送数据， 将数据发送给远程的客户端
	SendMsg(msgID uint32, data []byte) error
	/*
		conn member function
	*/
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemovePropertty(key string)
}

/*
handle function
*/
type HandleFunc func(*net.TCPConn, []byte, int) error
