package ziface

/*
	路由的抽象接口
	路由里的数据都是IRequest请求
*/

type IRouter interface {
	//  生命周期函数
	PreHandle(request IRequest)
	HostHandle(request IRequest)

	Handle(request IRequest)
}
