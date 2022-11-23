package znet

import "zinx/ziface"

type BaseRouter struct {
}

//所以Router全部继承BaseRouter的好处就是，  不需要实现PreHandle, PostHandle

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
