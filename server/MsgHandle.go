package server

import (
	"fmt"
	"github.com/guihai/ghtcpgs/impl"
	"github.com/guihai/ghtcpgs/utils"
)

type MsgHandle struct {
	// 消息id和路由映射 map  根据映射 获取路由
	Apis map[uint32]impl.IRouter

	// 创建协程池 工作单位+任务队列
	// 工作单位数量
	WorkPoolSize uint32
	// 任务队列 管道 切面，数据就是请求
	TaskQueue []chan impl.IRequest

	// WorkPoolIsOn
	PoolOn bool
}

func NewMsgHandle() *MsgHandle {
	// 必须初始化 map
	return &MsgHandle{
		Apis:         make(map[uint32]impl.IRouter),
		WorkPoolSize: utils.GO.WorkPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan impl.IRequest, utils.GO.WorkPoolSize),

		PoolOn: false, // 协程池未启动
	}
}

// 添加路由，也就是 添加 map数据
func (s *MsgHandle) AddRouter(msgId uint32, router impl.IRouter) {
	_, ok := s.Apis[msgId]
	if ok {
		// id 不能重复添加
		fmt.Println("消息类型", msgId, "已经添加过了，不能重复添加")
		return
	}
	// 数据不存在
	s.Apis[msgId] = router
	fmt.Println("消息类型", msgId, "添加路由成功")
	return

}

// 获取请求 处理 请求，请求（消息）=》 路由  协程处理
func (s *MsgHandle) DoMsgHandler(request impl.IRequest) {

	// 1，根据请求中的消息获取 路由
	r, ok := s.Apis[request.GetMsg().GetDataId()]
	if !ok {
		fmt.Println("消息类型", request.GetMsg().GetDataId(), "没有添加路由")
		return
	}

	// 获取路由成功 按照顺序执行
	r.PreHandle(request)
	r.Handle(request)
	r.PostHandle(request)
}

/*
初始化协程池，服务已开启就要启动
*/
func (s *MsgHandle) StartWorkerPool() {

	for i := 0; i < int(s.WorkPoolSize); i++ {

		// 初始化任务队列的管道
		s.TaskQueue[i] = make(chan impl.IRequest, utils.GO.TaskQueueMaxSize)

		// 启动一个工作单位
		go s.StartWorker(i, s.TaskQueue[i])
	}

	s.PoolOn = true
}

func (s *MsgHandle) StartWorker(i int, requests chan impl.IRequest) {
	fmt.Println("协程工作编号==", i, "==创建完成")

	//不断的等待队列中的消息,然后进行路由处理
	for {
		select {
		case re := <-requests:
			s.DoMsgHandler(re)
		}
	}

}

// 将消息交给TaskQueue,由worker进行处理
func (s *MsgHandle) SendMsgToTaskQueue(request impl.IRequest) {
	// 采用 链接id 取余数 放入对应的消息队列
	cid := request.GetConn().GetConId()
	mid := request.GetMsg().GetDataId()

	wid := cid % s.WorkPoolSize

	fmt.Println("添加 ConnID=", cid, " 请求 msgID=", mid, "到 workerID=", wid)

	s.TaskQueue[wid] <- request
}

// 获取协程池是否开启标志
func (s *MsgHandle) WorkPoolIsOn() bool {
	return s.PoolOn
}
