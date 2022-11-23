/*
package :: main
author :: shv
data :: 2022年11月24日03:45:22

***
包含了所有的路由注册， rpc服务注册
***
*/

package api

import (
	"zinx/shv-zinx_tcp/ziface"
	"zinx/shv-zinx_tcp/znet"
)

type C2S_Login struct {
	znet.BaseRouter
}

// handle regsiter
// type HandleRegsiter{
// 	znet.BaseRouter
// }

func HandleRegsiter(server ziface.Iserver) {
	server.AddRouter(1001, &C2S_Login{})
}
