package conf

type RunConfigObj struct {
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

var GO *RunConfigObj

// 初始化 全局变量，使用 init 方法 ，包被调用的时候会先调用这个方法
func NewRunConfigObj( // 服务名称
	Name string, IP string, Port string, Tcp string,
	MaxConn uint32, MaxPacketSize uint32, Version string,
	WorkPoolSize uint32, TaskQueueMaxSize uint32) {
	//初始化GlobalObject变量，设置一些默认值
	GO = &RunConfigObj{
		Name:          Name,
		IP:            IP,
		Port:          Port,
		Tcp:           Tcp,
		MaxConn:       MaxConn,
		MaxPacketSize: MaxPacketSize,
		Version:       Version,
		// 协程池
		WorkPoolSize: WorkPoolSize, // 和cpu 数量匹配合适
		// 协程池任务队列的最大容量
		TaskQueueMaxSize: TaskQueueMaxSize,
	}

}
