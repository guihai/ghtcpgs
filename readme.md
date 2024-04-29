tcp长连接服务
===
 引入代码
```go
package main

import (
	"github.com/guihai/ghtcpgs"
	"github.com/guihai/ghtcpgs/conf"
	"github.com/guihai/ghtcpgs/demo/router"
)


func main() {
	//

	conf.NewRunConfigObj("游戏服务", "0.0.0.0", "8999", "tcp4",
		100, 2048, "v9", 8, 1024)

	s := ghtcpgs.NewServer()

	/*
		添加多个路由，使用msg id 区分路由匹配
	*/

	s.SetOnConnStart(router.DoConnectionBegin)
	s.SetOnConnStop(router.DoConnectionLost)

	s.AddRouter(0, &router.PingHandle{})
	s.AddRouter(1, &router.HelloHandle{})

	s.Run()

}

```