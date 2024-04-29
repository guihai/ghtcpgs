package impl

/*
链接->请求->解包->消息->路由
业务->消息->封包->链接
*/

type IDataPack interface {
	// 打包
	Pack(IMessage) ([]byte, error)
	// 拆包
	UnPack([]byte) (IMessage, error)
	//获取包头长度方法
	GetHeadLen() uint32
}
