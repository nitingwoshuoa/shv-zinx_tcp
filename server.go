/*
package :: main
author :: shv
data :: 2022年11月24日03:36:05

***
zinx项目的启动文件
***
*/
package main

import (
	"github.com/nitingwoshuoa/shv-zinx_tcp/trunk/api"
	"github.com/nitingwoshuoa/shv-zinx_tcp/ziface"
	"github.com/nitingwoshuoa/shv-zinx_tcp/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	conn.SetProperty("HeartCount", 0)
}

func main() {
	//1 创建一个游戏服务器的server句柄，  使用Zinx的Api

	s := znet.NewServer("[Server Point]")

	// 注册hook函数
	s.SetOnConnStart(OnConnectionAdd)
	// 给当前zinx框架添加一个自定义的router
	// s.AddRouter(0, &PingRouter{})
	api.HandleRegsiter(s)
	//2 启动Sever
	s.Serve()

}
