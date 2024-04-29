package impl

/*
接口定义规则 ，服务接口
*/
type IServer interface {
	// 开启服务
	Start()

	// 停止服务
	Stop()

	// 运行服务。main中调用
	Run()

	// 添加路由
	AddRouter(uint32, IRouter)

	//得到链接管理
	GetConnMan() IConnManager

	//设置该Server的连接创建时Hook函数   函数和变量一样也可以作为其他函数的参数
	SetOnConnStart(func(IConn))

	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConn))

	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConn)

	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConn)
}
