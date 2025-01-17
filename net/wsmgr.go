package net

import (
	"log"
	"sync"
)

var Mgr = &WsMgr{
	userCache: make(map[int]WSConn),
}

type WsMgr struct {
	uc        sync.RWMutex
	userCache map[int]WSConn
}

func (m *WsMgr) UserLogin(conn WSConn, uid int, token string) {
	m.uc.Lock()
	defer m.uc.Unlock()

	oldConn := m.userCache[uid]
	if oldConn != nil {
		// 有用户登录着呢
		if conn != oldConn {
			// 有用户抢登录了
			log.Println("有用户抢登录了")
			oldConn.Push("relogin", nil)
		}
	}
	m.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}
