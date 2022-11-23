package znet

import (
	"fmt"
	"net"
	"zinx/shv-zinx_tcp/utils"
	"zinx/shv-zinx_tcp/ziface"
)

//Iserver 的接口实现，  定义一个Server的服务器模块

type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//当前的server 消息管理模块， 用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	//该server的连接管理器
	ConnMgr ziface.IConnManager
	//Hook
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

// 定义当前客户端连接所绑定的handle api  目前这个handle是写死的，以后优化应该由用户自定义handle方法
// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	//回显的业务
// 	fmt.Println("[Conn Handle] CallBackToClient ..")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("write back buf err", err)
// 		return errors.New("CallBackToCliente error")
// 	}
// 	return nil
// }

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name :%s listenner at IP : %s, Port : %d is starting\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPackectSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[start]  Server Listernner at IP: %s, Port %d is string \n", s.IP, s.Port)
	//1. 获取一个tcp的addr

	go func() {
		//0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("reslove tcp addr error :", err)
			return
		}
		//2. 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen :", s.IPVersion, "err ", err)
			return
		}

		var cid uint32 = 0

		fmt.Println("start Zinx server success", s.Name, "succ, listening...")

		//3. 阻塞的等待客户端进行连接，处理客户端连接业务（读写）
		for {
			//如果有客户端连接过来， 阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept err ", err)
				continue
			}

			//设置最大连接数，超出连接最大数关闭新的连接

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//todo : 返回一个最大连接的错误包
				fmt.Println("too many connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// 将处理新链接的业务方法和conn进行绑定得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的链接业务
			go dealConn.Start()
			// //已经与客户端建立了连接， 做一个最基本的512字节长度的回显业务
			// go func() {
			// 	for {
			// 		buf := make([]byte, 512)
			// 		cnt, err := conn.Read(buf)
			// 		if err != nil {
			// 			fmt.Println("recv buf err", err)
			// 			continue
			// 		}

			// 		fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)

			// 		//回显功能

			// 		if _, err := conn.Write(buf[:cnt]); err != nil {
			// 			fmt.Println("write back buf err", err)
			// 			continue
			// 		}

			// 	}
			// }()
		}
	}()

}

func (s *Server) Stop() {
	// todo  将一些服务器的资源，状态，或者一些已经开辟的链接信息 进行停止或者回收
	fmt.Println("[Stop] zinx server name = ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	//Todo  做一些启动服务器之后的额外业务
	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("add router succ !!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
  初始化Server模块的方法
*/

func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}

	return s
}

func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----->Call OnConnStart()")
		s.OnConnStart(connection)
	}
}
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----->Call OnConnStop()")
		s.OnConnStop(connection)
	}
}
