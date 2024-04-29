package main

import (
	"github.com/guihai/ghtcpgs/conf"
	"github.com/guihai/ghtcpgs/demo/router"
	"github.com/guihai/ghtcpgs/server"
)

func main() {
	//
	conf.NewRunConfigObj("游戏服务", "0.0.0.0", "8999", "tcp4",
		100, 2048, "v9", 8, 1024)
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
