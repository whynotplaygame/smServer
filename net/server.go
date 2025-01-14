package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type server struct {
	addr   string
	router *Router
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) Router(router *Router) {
	s.router = router
}

// 启动服务
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("login server started")
}

// http 升级到websocket
var wsUpgrader = websocket.Upgrader{
	// 允许所有cors跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 思考   websocket
	// 1, http 升级到 websocket
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("webocket 服务器出错", err)
	}
	log.Println("websocket服务连接成功")

	// websocket通道建立之后，不管是客户端，还是服务器，都可以收发消息
	// 发消息的时候，把消息当路由 来去处理，消息是有格式的，先定义消息格式
	// 客户端发消息的时候，{Name:"account.login"} 收到之后进行解析，认为要处理登录逻辑
	//sendErr := wsConn.WriteMessage(websocket.BinaryMessage, []byte("hello niyade"))
	//fmt.Println(sendErr)
	wsServer := NewWsServer(wsConn)
	wsServer.router = s.router
	wsServer.Star()
	wsServer.Handshake()
}
