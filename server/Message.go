package server

type Message struct {
	// 消息长度
	dataLen uint32

	// 消息id
	dataId uint32

	// 消息内容
	data []byte
}

func NewMessage(dataId uint32, data []byte) *Message {
	return &Message{
		dataLen: uint32(len(data)),
		dataId:  dataId,
		data:    data,
	}
}

//获取消息长度
func (s *Message) GetDataLen() uint32 {
	return s.dataLen
}

// 获取消息id
func (s *Message) GetDataId() uint32 {
	return s.dataId
}

// 获取消息内容
func (s *Message) GetData() []byte {
	return s.data
}

// 设定消息id
func (s *Message) SetDataId(dataId uint32) {
	s.dataId = dataId
}

// 设置消息长度
func (s *Message) SetDataLen(dataLen uint32) {
	s.dataLen = dataLen
}

// 设置消息内容
func (s *Message) SetData(data []byte) {
	s.data = data
}
