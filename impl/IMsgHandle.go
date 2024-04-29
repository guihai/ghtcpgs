package impl

/*
 消息管理器，用来存放 消息id和路由的映射
连接=》请求=》拆包=》消息=》消息管理=》路由=》业务
*/

type IMsgHandle interface {
	// 添加路由，也就是 添加 map数据
	AddRouter(msgId uint32, router IRouter)

	// 获取请求 处理 请求，请求（消息）=》 路由  协程处理
	DoMsgHandler(request IRequest)

	//启动worker工作池
	StartWorkerPool()

	//将消息交给TaskQueue,由worker进行处理
	SendMsgToTaskQueue(request IRequest)

	// 获取协程池是否开启标志
	WorkPoolIsOn() bool
}
