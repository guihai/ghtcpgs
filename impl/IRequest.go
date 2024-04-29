package impl

/*
请求封装
请求和路由，是针对具体业务的
*/

type IRequest interface {
	//得到当前链接
	GetConn() IConn

	// 获取请求数据
	GetMsg() IMessage
}
