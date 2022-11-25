package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/nitingwoshuoa/shv-zinx_tcp/utils"
	"github.com/nitingwoshuoa/shv-zinx_tcp/ziface"
)

/*
链接模块
*/
type Connection struct {
	//当前conn隶属于哪个server
	TcpServer ziface.Iserver

	//当前链接的socket TCP
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	// handleAPI ziface.HandleFunc

	ExitChan chan bool // 如果想退出链接要通过该管道告知    告知当前链接已经退出的/停止 channel (由reader告知writer退出)

	//无缓冲的管道，用于读写goroutine之间的通信
	msgChan chan []byte

	// 该链接处理的方法Router
	MsgHandle ziface.IMsgHandle

	// 链接属性集合
	property map[string]interface{}

	// 链接属性的保护锁
	propertyLock sync.RWMutex
}

//初始化链接模块的方法

func NewConnection(server ziface.Iserver, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		// handleAPI: callback_api,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
		MsgHandle: msgHandler,
		property:  make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 客户端连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[reader goroutinue is running]")

	defer fmt.Println("connID = ", c.ConnID, " [reder is exit, remote addr is]", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆包的对象
		dp := NewDataPack()

		//读取客户端的 msg head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}
		// 拆包， 得到msgID和msgDatalen   放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err ", err)
		}
		// 根据dataLen   再次读取data 放在 msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err ", err)
				break
			}
		}

		msg.SetData(data)
		//得到当前conn数据的Request请求数据

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//判断已经开启了工作池机制
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到注册绑定的conn对应的router调用
			//根据绑定好的MsgID， 找到对应处理的api业务
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}

}

// 写消息的goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Print(c.RemoteAddr().String(), "[conn writer exit!]")
	//不断的阻塞等待channel的消息 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error ", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出，此时Writer也要退出

			return
		}
	}
}

// 启动链接  让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("ConnStart() ... ConnID = ", c.ConnID)
	//启动从当前链接的读数据的业务

	go c.StartReader()
	go c.StartWriter()

	//todo 启动从当前链接写数据的业务
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	//调用开发者注册的业务
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知writer关闭
	c.ExitChan <- true
	//回收资源

	//将当前链接从connmgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态  IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()

}

// 提供一个SendMsg方法  将我们要发送给客户端得数据，先进行封包，再发送

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包
	dp := NewDataPack()

	// binaryMsg 是已经打包好的 MsgDataLen/MsgID/Data格式
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg")
	}

	// 将数据发送给管道  不用直接发给客户端了

	c.msgChan <- binaryMsg

	// if _, err := c.Conn.Write(binaryMsg); err != nil {
	// 	fmt.Println("write msg id ", msgID, " error : ", err)
	// 	return errors.New("conn write error ")
	// }

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemovePropertty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
