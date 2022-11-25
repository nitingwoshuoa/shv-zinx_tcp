package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/nitingwoshuoa/shv-zinx_tcp/znet"
)

/*
模拟客户端
*/

func main() {

	fmt.Println("client start ...")
	// 1 直接链接远程服务器， 得到一个conn链接

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")

	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {

		// 发送封包的msg消息  msgid:0

		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1002, []byte("zinv 0.5 client test message")))

		if err != nil {
			fmt.Println("Pack error", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error")
			return
		}

		//服务器回复我们一个数据，， MsgID:1 pingpingping

		//先读取流中的head部分 得到id 和 datalen

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		//将二进制的head
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			//msg 中有数据  再根据datalen 进行第二次读取将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error ", err)
				return
			}

			fmt.Println("recv server msg : id = ", msg.Id, ", len = ", msg.DataLen, ", data = ", string(msg.Data))

		}

		// // 链接调用Write  写数据
		// _, err := conn.Write([]byte("hello zinx v0.5"))

		// if err != nil {
		// 	fmt.Println("write conn err ", err)
		// 	return
		// }

		// buf := make([]byte, 512)
		// cnt, err := conn.Read(buf)

		// if err != nil {
		// 	fmt.Println("read buf error")
		// 	return
		// }
		// fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}

}
