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
	"zinx/shv-zinx_tcp/trunk/api"
	_ "zinx/shv-zinx_tcp/ziface"
	"zinx/shv-zinx_tcp/znet"
)

func main() {
	//1 创建一个游戏服务器的server句柄，  使用Zinx的Api

	s := znet.NewServer("[Server Point]")

	// 注册hook函数
	s.Set
	// 给当前zinx框架添加一个自定义的router
	// s.AddRouter(0, &PingRouter{})
	api.HandleRegsiter(s)
	//2 启动Sever
	s.Serve()

}
