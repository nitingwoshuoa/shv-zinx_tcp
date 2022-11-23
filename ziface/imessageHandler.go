package ziface

type IMsgHandle interface {
	// router function
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	/*
		work pool function
	*/
	StartWorkPool()
	SendMsgToTaskQueue(request IRequest)
}
