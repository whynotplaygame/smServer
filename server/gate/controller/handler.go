package controller

import (
	"fmt"
	"smServer/config"
	"smServer/net"
	"strings"
	"sync"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int64]*net.ProxyClient),
}

type Handler struct {
	proxyMutex sync.Mutex
	// 代理地址 --》 客户端连接 （游戏客客户端连接id --> 连接）
	proxyMap   map[string]map[int64]*net.ProxyClient
	loginProxy string
	gameProxy  string
}

func (h *Handler) Router(r *net.Router) {
	h.loginProxy = config.File.MustValue("gate_server", "login_proxy", "ws://127.0.0.1:8003")
	h.gameProxy = config.File.MustValue("gate_server", "game_proxy", "ws://127.0.0.1:8001")
	g := r.Group("*")
	g.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	fmt.Println("网关的处理器。。。。。。")
	// account 转发
	name := req.Body.Name // 就是请求的标识，路由的名称
	proxyStr := ""
	if isAccount(name) { //通过标识判断是否是账号相关
		proxyStr = h.loginProxy
	}
	proxy := net.NewProxyClient(proxyStr) // 创建了登录的代理对的
	proxy.Connect()                       // 创建出来需要进行连接
}

func isAccount(name string) bool {
	return strings.HasPrefix(name, "account")
}
