package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
	消息处理模块
*/

type MsgHandle struct {
	//存放每个msgid 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker读取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的数量
	WorkerPoolSize uint32
}

// 初始化 创建 MsgHandle方法

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置的参数中获取，也可以在配置文件中让用户进行配置
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 从 request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msg ID =", request.GetMsgID(), " is not found ! need register")
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
	// 根据msgid 调度对应router业务即可
}
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//当前的msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	mh.Apis[msgID] = router
	fmt.Println("add api Msg id = ", msgID, " succ!")
}

// 启动一个worker工作池  一个server节点开启工作池的动作只能发生一次，  对外暴露
func (mh *MsgHandle) StartWorkPool() {
	//根据workerPoolSze分别开启worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 当前的worker对应的channel消息队列  开辟空间  第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

// 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started ..")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}

}

// 将消息交给TaskQueue， 由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同的worker， 按照轮询
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connID = ", request.GetConnection().GetConnID(), " request Msgid = ", request.GetMsgID(), " to worker id = ", workerID, " is started")
	mh.TaskQueue[workerID] <- request
}
