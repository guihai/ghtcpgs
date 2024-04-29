package utils

import (
	"encoding/json"
	"io/ioutil"
)

/*
全局配置对象
*/

type GlobalObj struct {
	// 服务名称
	Name string
	// 绑定ip
	IP string
	// 绑定端口
	Port string
	// 传输协议
	Tcp string
	// 最大连接数
	MaxConn uint32
	// 数据包最大值
	MaxPacketSize uint32
	// 版本号
	Version string

	// 协程池
	WorkPoolSize uint32
	// 协程池任务队列的最大容量
	TaskQueueMaxSize uint32
}

// 定义全局使用的变量
var GO *GlobalObj

// 初始化 全局变量，使用 init 方法 ，包被调用的时候会先调用这个方法
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GO = &GlobalObj{
		Name:          "游戏服务",
		IP:            "0.0.0.0",
		Port:          "8999",
		Tcp:           "tcp4",
		MaxConn:       100,
		MaxPacketSize: 2048,
		Version:       "v9",

		// 协程池
		WorkPoolSize: 10, // 和cpu 数量匹配合适
		// 协程池任务队列的最大容量
		TaskQueueMaxSize: 1024,
	}

	// 加载配置参数
	reloadConfig()

}

func reloadConfig() {
	buf, err := ioutil.ReadFile("etc/etc.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, GO)
	if err != nil {
		panic(err)
	}
}
