package controller

import (
	"fmt"
	"log"
	"smServer/config"
	"smServer/constant"
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
	} else {
		proxyStr = h.gameProxy
	}
	if proxyStr == "" {
		rsp.Body.Code = constant.ProxyNotInConnect
		return
	}
	h.proxyMutex.Lock()

	_, ok := h.proxyMap[proxyStr]
	if !ok {
		h.proxyMap[proxyStr] = map[int64]*net.ProxyClient{}
	}
	h.proxyMutex.Unlock()
	//客户端id
	c, err := req.Conn.GetProperty("cid")
	if err != nil {
		log.Println("cid 没有取到", err)
		rsp.Body.Code = constant.InvalidParam
		return
	}
	cid := c.(int64) // 转成 int64
	proxy, ok := h.proxyMap[proxyStr][cid]

	if !ok {
		proxy = net.NewProxyClient(proxyStr) // 创建了登录的代理对的
		h.proxyMutex.Lock()
		h.proxyMap[proxyStr][cid] = proxy
		h.proxyMutex.Unlock()

		err := proxy.Connect() // 创建出来需要进行连接
		if err != nil {
			h.proxyMutex.Lock()
			delete(h.proxyMap[proxyStr], cid) // 删除
			h.proxyMutex.Unlock()
			rsp.Body.Code = constant.ProxyConnectError
			return
		}
		h.proxyMap[proxyStr][cid] = proxy
		proxy.SetProperty("cid", cid)
		proxy.SetProperty("proxy", proxyStr)
		proxy.SetProperty("gateConn", req.Conn)
		proxy.SetOnPush(h.onPush)

	}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	r, err := proxy.Send(req.Body.Name, req.Body.Msg) // 发送后，应该有响应
	if err == nil {
		rsp.Body.Code = r.Code
		rsp.Body.Msg = r.Msg
	} else {
		rsp.Body.Code = constant.ProxyConnectError
		rsp.Body.Msg = nil
	}
}

func (h *Handler) onPush(conn *net.ClientConn, body *net.RspBody) {
	gc, err := conn.GetProperty("gateConn")
	if err != nil {
		log.Println("on push gateConn:", err)
		return
	}
	gateConn := gc.(net.WSConn)
	gateConn.Push(body.Name, body.Msg)
}

func isAccount(name string) bool {
	return strings.HasPrefix(name, "account")
}
