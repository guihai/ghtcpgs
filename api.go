package api

import (
	"github.com/guihai/ghtcpgs/conf"
	"github.com/guihai/ghtcpgs/server"
)

// 对外暴露接口
func NewServer() *server.Server {
	return &server.Server{
		Name:      conf.GO.Name,
		IP:        conf.GO.IP,
		Port:      conf.GO.Port,
		Tcp:       conf.GO.Tcp,
		MsgHandle: server.NewMsgHandle(),
		ConnMan:   server.NewConnManager(), // 链接管理器
	}
}
