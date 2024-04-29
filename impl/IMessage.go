package impl

/*
请求->消息->路由-> 业务

*/

type IMessage interface {
	//获取消息长度
	GetDataLen() uint32

	// 获取消息id
	GetDataId() uint32

	// 获取消息内容
	GetData() []byte

	// 设定消息id
	SetDataId(uint32)

	// 设置消息长度
	SetDataLen(uint32)

	// 设置消息内容
	SetData([]byte)
}
