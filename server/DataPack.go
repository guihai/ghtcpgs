package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/guihai/ghtcpgs/impl"
	"github.com/guihai/ghtcpgs/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 打包
/*
数据格式  len(4)+id(4)+data
先写 len
在写 id
在写 data
*/
func (s *DataPack) Pack(msg impl.IMessage) ([]byte, error) {

	// 读取消息长度
	msg.GetDataLen() // 返回消息长度

	if len(msg.GetData()) == 0 {
		// 数据空
		return nil, errors.New("消息内容为空")
	}

	/*
		将消息 写入 二级制字节
	*/
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 先写入 长度   uint32  占据四个字节
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 在写入 消息id  uint32  占据四个字节
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataId()); err != nil {
		return nil, err
	}

	// 以上两项 占了  4+4 = 8 个字节
	// 最后写入消息内容
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	// 写入完成返回
	return dataBuff.Bytes(), nil
}

// 拆包
/*
获取消息的 len 和 id  不获取data 有了 长度 下次解析数据到消息
数据格式  len(4)+id(4)+data
*/
func (s *DataPack) UnPack(buf []byte) (impl.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(buf)

	//只解压head的信息，得到dataLen和msgID
	var dataLen, dataId uint32
	dataLen = 0
	dataId = 0
	// 先要读取长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &dataLen); err != nil {
		return nil, err
	}

	// 读取数据id
	if err := binary.Read(dataBuff, binary.LittleEndian, &dataId); err != nil {
		return nil, err
	}

	if utils.GO.MaxPacketSize > 0 && dataLen > utils.GO.MaxPacketSize {
		return nil, errors.New("消息太长。。。。超过限度")
	}

	// 最终返回
	return &Message{
		dataLen: dataLen,
		dataId:  dataId,
		data:    []byte{},
	}, nil

}

// 获取包头长度方法
func (s *DataPack) GetHeadLen() uint32 {
	// 数据固定为  len(4)+id(4)+data
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}
