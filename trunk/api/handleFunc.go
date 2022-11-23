/*
package :: main
author :: shv
data :: 2022年11月24日04:07:47

***
路由函数，业务功能
***
*/

package api

import (
	"fmt"
	"zinx/shv-zinx_tcp/ziface"
	_ "zinx/shv-zinx_tcp/znet"
)

func (this *C2S_Login) Handle(request ziface.IRequest) {
	//test
	fmt.Println("user login succ")

	//todo 集成protobuf 序列化与反序列化  完成登录逻辑
}
