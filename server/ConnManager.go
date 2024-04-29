package server

import (
	"errors"
	"fmt"
	"github.com/guihai/ghtcpgs/impl"
	"sync"
)

type ConnManager struct {
	// 链接map
	ConnMap map[uint32]impl.IConn
	// 对资源进行增删改查，需要上锁，单纯增不用锁
	lock sync.RWMutex // 读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		ConnMap: make(map[uint32]impl.IConn),
	}
}

// 添加链接
func (s *ConnManager) Add(conn impl.IConn) {

	s.lock.Lock()         // 上写锁
	defer s.lock.Unlock() // 解锁

	s.ConnMap[conn.GetConId()] = conn

	fmt.Println("链接 添加到  链接管理器中！ 链接数量  = ", s.Len())

}

// 删除连接
func (s *ConnManager) Remove(conn impl.IConn) {

	s.lock.Lock()         // 上写锁
	defer s.lock.Unlock() // 解锁

	delete(s.ConnMap, conn.GetConId())

	fmt.Println("链接 cid", conn.GetConId(), " 已删除！ 链接数量  = ", s.Len())

}

// 利用ConnID获取链接
func (s *ConnManager) Get(connID uint32) (impl.IConn, error) {

	s.lock.RLock()         // 上读锁
	defer s.lock.RUnlock() // 解锁

	c, ok := s.ConnMap[connID]
	if ok {
		return c, nil
	}
	return nil, errors.New("链接不存在")

}

// 获取当前连接 个数
func (s *ConnManager) Len() int {

	return len(s.ConnMap)

}

// 删除并停止所有链接
func (s *ConnManager) ClearConn() {

	s.lock.Lock()         // 上写锁
	defer s.lock.Unlock() // 解锁

	for u, conn := range s.ConnMap {

		// 链接停止
		conn.Stop()

		// map 中删除
		delete(s.ConnMap, u)
	}

	fmt.Println("链接管理器已清空！ 链接数量  = ", s.Len())
}
