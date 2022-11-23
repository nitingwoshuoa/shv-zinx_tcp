package znet

import "zinx/shv-zinx_tcp/ziface"

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) HostHandle(request ziface.IRequest) {}
