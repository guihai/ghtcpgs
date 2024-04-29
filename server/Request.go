package server

import "github.com/guihai/ghtcpgs/impl"

type Request struct {
	// 请求链接
	cn impl.IConn

	// 请求消息
	msg impl.IMessage
}

func NewRequest(cn impl.IConn, msg impl.IMessage) *Request {
	return &Request{
		cn: cn,
		//data: data,
		msg: msg,
	}
}

// 得到当前链接
func (s *Request) GetConn() impl.IConn {
	return s.cn
}

// 获取请求数据
func (s *Request) GetMsg() impl.IMessage {
	return s.msg
}
