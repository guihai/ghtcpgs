package impl

import (
	"net"
)

/*
处理tcp链接 接口

*/

type IConn interface {

	// 启动链接
	Start()

	// 关闭链接
	Stop()

	// 获取链接对象
	GetTcpCon() *net.TCPConn

	// 获取链接地址
	GetConAddr() net.Addr

	// 获取链接id
	GetConId() uint32

	//直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error

	//直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgId uint32, data []byte) error //添加带缓冲发送消息接口

	//设置链接属性
	SetProperty(key string, value interface{})

	//获取链接属性
	GetProperty(key string) (interface{}, error)

	//移除链接属性
	RemoveProperty(key string)

	//返回ctx，用于用户自定义的go程获取连接退出状态
	// todo zinx 这里修改了
	//Context() context.Context
}
