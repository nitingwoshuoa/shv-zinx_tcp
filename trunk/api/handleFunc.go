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
	"gopkg.in/mgo.v2/bson"
	"github.com/nitingwoshuoa/shv-zinx_tcp/zdbTool"
	"github.com/nitingwoshuoa/shv-zinx_tcp/ziface"

	_ "github.com/nitingwoshuoa/shv-zinx_tcp/znet"
)

func (handle *C2S_Login) Handle(request ziface.IRequest) {
	//test
	fmt.Println("user login succ")
	fmt.Println(zdbTool.Count("config", "System", bson.M{"_id": 1}))
	//todo 集成protobuf 序列化与反序列化  完成登录逻辑
}

func (handle *C2S_HeartPackage) Handle(request ziface.IRequest) {
	request.GetConnection().SetProperty("HeartCount", 0),
	fmt.Printf("conn == : %d ", request.GetConnection())
}
