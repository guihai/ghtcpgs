package router

import (
	"fmt"
	"github.com/guihai/ghtcpgs/impl"
	"github.com/guihai/ghtcpgs/server"
)

type PingHandle struct {
	*server.BaseRouter
}

// 业务处理
func (s *PingHandle) Handle(request impl.IRequest) {
	fmt.Println("开始 PingHandle  Handle 。。。")
	//request.GetConn().GetTcpCon().Write([]byte("发送 Handle   \n"))

	// 获取消息
	msg := request.GetMsg()
	fmt.Println("收到客户端消息=>id:", msg.GetDataId(), "=>长度:", msg.GetDataLen(), "=>内容", string(msg.GetData()))

	request.GetConn().SetProperty("tom", "汤姆")
	request.GetConn().SetProperty("jack", 19)
	// 发送消息
	str := "ping ok"

	request.GetConn().SendMsg(100, []byte(str))

}

// 启动钩子
// 创建连接的时候执行
func DoConnectionBegin(conn impl.IConn) {
	fmt.Println("DoConnecionBegin is Called ... ")
	err := conn.SendMsg(301, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn impl.IConn) {
	fmt.Println("DoConneciotnLost is Called ... ")
}

type HelloHandle struct {
	*server.BaseRouter
}

// 业务处理
func (s *HelloHandle) Handle(request impl.IRequest) {
	fmt.Println("开始 HelloHandle  Handle 。。。")
	//request.GetConn().GetTcpCon().Write([]byte("发送 Handle   \n"))

	// 获取消息
	msg := request.GetMsg()
	fmt.Println("收到客户端消息=>id:", msg.GetDataId(), "=>长度:", msg.GetDataLen(), "=>内容", string(msg.GetData()))

	tom, _ := request.GetConn().GetProperty("tom")
	jack, _ := request.GetConn().GetProperty("jack")
	fmt.Println("tom=", tom)
	fmt.Println("jack=", jack)
	// 发送消息
	str := "hello ok"

	request.GetConn().SendMsg(101, []byte(str))

}
