package server

import "github.com/guihai/ghtcpgs/impl"

/*
路由基类，路由继承此类，实现接口方法，方法没有具体内容
*/
type BaseRouter struct {
}

// 业务处理前
func (s *BaseRouter) PreHandle(request impl.IRequest) {}

// 业务处理
func (s *BaseRouter) Handle(request impl.IRequest) {}

// 业务处理后
func (s *BaseRouter) PostHandle(request impl.IRequest) {}
