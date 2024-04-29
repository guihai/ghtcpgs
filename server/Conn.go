package server

import (
	"errors"
	"fmt"
	"github.com/guihai/ghtcpgs/impl"
	"github.com/guihai/ghtcpgs/utils"
	"io"
	"net"
	"sync"
)

type Conn struct {
	//当前Conn属于哪个Server
	tcpServer impl.IServer //当前conn属于哪个server，在conn初始化的时候添加即可

	// 当前连接的socket TCP套接字
	cn *net.TCPConn
	// 链接id 也可以称作为SessionID，ID全局唯一
	cid uint32
	// 链接状态
	isClose bool

	// 等待链接退出的管道 ，用此管道来监听关闭
	exitChan chan bool

	// 使用具体的路由进行业务处理
	//router impl.IRouter
	// 路由管理器
	MsgHandle impl.IMsgHandle

	// 写数据通道，无缓冲，将需要写的数据放入通道，对应的监听协程就会工作
	writerChan chan []byte

	// 写数据通道，有缓冲，将需要写的数据放入通道，对应的监听协程就会工作
	writerBuffChan chan []byte

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewConn(tcpS impl.IServer, cn *net.TCPConn, cid uint32, msgHandle impl.IMsgHandle) *Conn {
	c := &Conn{
		tcpServer: tcpS,
		cn:        cn,
		cid:       cid,
		isClose:   false, // 初始化就是开启中

		exitChan: make(chan bool, 1), // 建立缓存管道

		writerChan:     make(chan []byte),                         // 返回写数据通道
		writerBuffChan: make(chan []byte, utils.GO.MaxPacketSize), // 返回写数据通道

		MsgHandle: msgHandle,

		property: make(map[string]interface{}), // 属性初始化
	}

	// 创建链接的时候将链接加入到链接管理器
	c.tcpServer.GetConnMan().Add(c)
	return c
}

func (s *Conn) Start() {
	fmt.Println("开启链接服务cid=>", s.cid, "地址=》", s.GetConAddr())

	// 启动读数据
	go s.StartRead()

	// 启动写数据
	go s.StartWriter()

	// 启用钩子
	s.tcpServer.CallOnConnStart(s)

	// 监听关闭管道，如果关闭就返回
	for {
		select {
		case <-s.exitChan:
			// 通道关闭了
			return
		}

	}
}

func (s *Conn) Stop() {
	if s.isClose {
		// 已经关闭，不处理
		return
	}

	s.isClose = true

	// 启用钩子
	s.tcpServer.CallOnConnStop(s)
	// 关闭链接
	s.cn.Close()

	// 从链接管理器中删除链接
	s.tcpServer.GetConnMan().Remove(s)

	//通知从缓冲队列读数据的业务，该链接已经关闭
	s.exitChan <- true

	// 通知写工作 结束
	close(s.writerChan)

	// 关闭管道
	close(s.exitChan)

	fmt.Println("关闭链接服务,", s.cid)

}

func (s *Conn) GetTcpCon() *net.TCPConn {
	return s.cn

}

func (s *Conn) GetConAddr() net.Addr {
	return s.cn.RemoteAddr()

}

func (s *Conn) GetConId() uint32 {
	return s.cid

}
func (s *Conn) Send(buf []byte) error {

	// 写数据通道
	return nil
}

/*
开启读数据工作
*/
func (s *Conn) StartRead() {

	// 最终结束
	defer s.Stop()

	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head   能获取 长度 和id  8个字节
		headData := make([]byte, dp.GetHeadLen())

		//  io readfull 将byte 读满
		if _, err := io.ReadFull(s.GetTcpCon(), headData); err != nil {
			fmt.Println("获取链接的头部消息失败 ", err)
			s.exitChan <- true // 通道数据
			// 获取消息失败，con 就关闭了
			break
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("消息拆包失败 ", err)
			s.exitChan <- true // 通道数据
			break
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			// 读取消息内容 到data
			if _, err := io.ReadFull(s.GetTcpCon(), data); err != nil {
				fmt.Println("消息内容获取失败 ", err)
				s.exitChan <- true // 通道数据
				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			cn:  s,
			msg: msg,
		}

		// 使用协程池
		if s.MsgHandle.WorkPoolIsOn() {

			s.MsgHandle.SendMsgToTaskQueue(&req)

		} else {
			// 开启协程处理 路由
			go s.MsgHandle.DoMsgHandler(&req)
		}

	}
}

/*
开启写数据 工作
*/
func (s *Conn) StartWriter() {

	fmt.Println("连接：", s.cid, "，写工作开启")
	/*
		监听通道
	*/
	for {
		select {
		case data := <-s.writerChan:
			// 写 通道数据获取
			if _, err := s.cn.Write(data); err != nil {
				fmt.Println("写给客户端数据失败:, ", err, "连接退出")
				return
			}
		case data, ok := <-s.writerBuffChan:
			// 有缓冲通道
			if ok {
				//有数据要写给客户端
				if _, err := s.cn.Write(data); err != nil {
					fmt.Println("写给客户端数据失败:, ", err, "连接退出")
					return
				}
			} else {
				break
				fmt.Println("有缓冲写通道关闭了")
			}

		case <-s.exitChan:
			// 连接状态通道 获取数据，连接关闭了
			return

		}

	}
}

/*
将data 打包发出
*/
//直接将Message数据发送数据给远程的TCP客户端
func (s *Conn) SendMsg(msgId uint32, data []byte) error {

	if s.isClose {
		return errors.New("链接已关闭")
	}

	msg := NewMessage(msgId, data)
	// 将数据和id 打包
	dp := NewDataPack()

	by, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("消息打包错误 msg id = ", msgId, err)
		return errors.New("消息打包错误 ")
	}

	// 消息写入通道
	s.writerChan <- by

	return nil
}

// 添加带缓冲发送消息接口
func (s *Conn) SendBuffMsg(msgId uint32, data []byte) error {
	if s.isClose {
		return errors.New("链接已关闭")
	}

	msg := NewMessage(msgId, data)
	// 将数据和id 打包
	dp := NewDataPack()

	by, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("消息打包错误 msg id = ", msgId, err)
		return errors.New("消息打包错误 ")
	}

	// 消息写入有缓存通道
	s.writerBuffChan <- by

	return nil
}

// 设置链接属性
func (s *Conn) SetProperty(key string, value interface{}) {

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	s.property[key] = value
}

// 获取链接属性
func (s *Conn) GetProperty(key string) (interface{}, error) {

	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	v, ok := s.property[key]
	if ok {
		return v, nil
	}
	return nil, errors.New("没有此key")
}

// 移除链接属性
func (s *Conn) RemoveProperty(key string) {

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	delete(s.property, key)

}
