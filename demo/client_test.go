package main

import (
	"fmt"
	"github.com/guihai/ghtcpgs/server"
	"io"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	Client()
}

func Client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8991")
	if err != nil {
		fmt.Println("链接失败", err)
		return
	}
	defer conn.Close()

	fmt.Println("链接成功")

	back(conn)

}

// 使用消息格式  len(4)+ id(4)+data 发送
func back(conn net.Conn) {

	dp := server.NewDataPack()

	str1 := "客户端发送消息i"

	str2 := "客户端发送消息i+20"

	for i := 0; i < 20; i++ {

		msg1 := server.NewMessage(0, []byte(str1))
		msg2 := server.NewMessage(1, []byte(str2))

		msg3 := server.NewMessage(2, []byte(str2))

		by1, _ := dp.Pack(msg1)

		by2, _ := dp.Pack(msg2)

		by3, _ := dp.Pack(msg3)

		//by1 = append(by1, by2...)

		conn.Write(by1)
		conn.Write(by2)
		conn.Write(by3)

		// 获取收到的消息

		// 首先获取头部数据
		headData := make([]byte, 8)
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("获取链接的头部消息失败 ", err)
		}

		msg11, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("消息拆包失败 ", err)
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data11 []byte
		if msg11.GetDataLen() > 0 {
			data11 = make([]byte, msg11.GetDataLen())
			// 读取消息内容 到data
			if _, err := io.ReadFull(conn, data11); err != nil {
				fmt.Println("消息内容获取失败 ", err)
			}
		}
		msg11.SetData(data11)

		fmt.Println("收到服务端消息=>id:", msg11.GetDataId(), "=>长度:", msg11.GetDataLen(), "=>内容", string(msg11.GetData()))

		time.Sleep(time.Second)
	}

}
