package main

import (
	"github.com/guihai/ghtcpgs/demo/router"
	"github.com/guihai/ghtcpgs/server"
)

func main() {
	//
	s := server.NewServer()

	/*
		添加多个路由，使用msg id 区分路由匹配
	*/

	s.SetOnConnStart(router.DoConnectionBegin)
	s.SetOnConnStop(router.DoConnectionLost)

	s.AddRouter(0, &router.PingHandle{})
	s.AddRouter(1, &router.HelloHandle{})

	s.Run()
}
