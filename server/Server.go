package server

import (
	"fmt"
	"github.com/guihai/ghtcpgs/conf"
	"github.com/guihai/ghtcpgs/impl"
	"net"
)

type Server struct {
	// 服务名称
	Name string
	// 绑定ip
	IP string
	// 绑定端口
	Port string
	// 传输协议
	Tcp string

	// 路由管理器
	MsgHandle impl.IMsgHandle

	// 链接管理器
	ConnMan impl.IConnManager

	// 启动钩子
	OnConnStart func(conn impl.IConn)

	// 关闭钩子
	OnConnStop func(conn impl.IConn)

	// 全局运行配置
}

// 构造器
func NewServer() *Server {
	return &Server{
		Name:      conf.GO.Name,
		IP:        conf.GO.IP,
		Port:      conf.GO.Port,
		Tcp:       conf.GO.Tcp,
		MsgHandle: NewMsgHandle(),
		ConnMan:   NewConnManager(), // 链接管理器
	}
}

// 启动服务
func (s *Server) Start() {

	go func() {
		// 开启协程池，等待工作
		s.MsgHandle.StartWorkerPool()

		// 创建地址
		addr, err := net.ResolveTCPAddr(s.Tcp, s.IP+":"+s.Port)
		if err != nil {
			fmt.Println("!!启动失败", err)
			panic(err)
		}
		// 创建监听
		lis, err := net.ListenTCP(s.Tcp, addr)
		if err != nil {
			fmt.Println("!!启动监听失败", err)
			panic(err)
		}

		// 监听创建成功，开启循环接收 tcp链接
		fmt.Printf("====%s==== 服务已启动 \n", s.Name)
		fmt.Printf("启动端口%s,服务协议%s,版本号%s \n", s.Port, s.Tcp, conf.GO.Version)

		var cid uint32
		cid = 0
		for {
			// 接收tcp链接
			con, err := lis.AcceptTCP()
			if err != nil {

				fmt.Println("tcp 链接失败", err)

				continue
			}

			// 验证链接个数
			if s.ConnMan.Len() >= int(conf.GO.MaxConn) {

				// 链接数超了不能 加入
				con.Close()
				continue
			}

			// 创建 con 处理对象 传入路由
			cn := NewConn(s, con, cid, s.MsgHandle)

			cid++

			go cn.Start()

		}

	}()
}

// 关闭服务
func (s *Server) Stop() {
	// 清空链接管理器
	s.ConnMan.ClearConn()
	fmt.Printf("====%s==== 服务关闭 \n", s.Name)
}

// 运行服务
func (s *Server) Run() {
	fmt.Printf("====%s==== 服务启动中。。。 \n", s.Name)

	// 调用启动方法
	s.Start()

	// 开启阻塞
	select {}

}

// 添加路由
func (s *Server) AddRouter(msgId uint32, router impl.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)

}

// 得到链接管理
func (s *Server) GetConnMan() impl.IConnManager {
	return s.ConnMan
}

// 设置该Server的连接创建时Hook函数   函数和变量一样也可以作为其他函数的参数
func (s *Server) SetOnConnStart(fu func(conn impl.IConn)) {
	s.OnConnStart = fu
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(fu func(conn impl.IConn)) {
	s.OnConnStop = fu
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn impl.IConn) {
	if s.OnConnStart != nil {
		fmt.Println("启动钩子运行")
		s.OnConnStart(conn)
	}

}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn impl.IConn) {
	if s.OnConnStop != nil {
		fmt.Println("关闭钩子运行")
		s.OnConnStop(conn)
	}

}
