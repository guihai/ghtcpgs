package impl

/*
路由封装
请求-分发路由-具体业务
*/
type IRouter interface {
	// 业务处理前
	PreHandle(IRequest)

	// 业务处理
	Handle(IRequest)

	// 业务处理后
	PostHandle(IRequest)
}
