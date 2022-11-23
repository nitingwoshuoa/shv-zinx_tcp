package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/shv-zinx_tcp/ziface"
)

type GlobalObj struct {
	/*
		Server
	*/

	TcpServer ziface.Iserver //Zinx全局的Server对象
	Host      string         //服务器主机监听的ip
	TcpPort   int            //服务器主机监听的端口号
	Name      string         //服务器名称
	/*
		Zinx
	*/
	Version          string //版本号
	MaxConn          int    //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作池的Gorotinue数量
	MaxWorkerTaskLen uint32 //允许的用户最多开辟多少个worker
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &GlobalObject)

}

func init() {
	//default value
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.5",
		TcpPort:          8999,
		Host:             "127.0.0.1",
		MaxConn:          3,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	// GlobalObject.Reload()
}
