package main

import (
	"smServer/config"
	"smServer/net"
	"smServer/server/gate"
)

/*
	1, 登录功能account.login, 需要通过网关，转发，登录服务器
	2，网关（websocket的客户端） 如何和 登录服务器 交互
	3，网关又和游戏的客户端 进行交互，网络是websocket的服务器
	4，websocket的服务器 已经实现了
	5，需要实现websocket的客户端
	6，网关：代理服务器 （代理地址，代理的连接通道） 客户端连接（websocket连接）
	7，路由： 路由接受所有的请求，网关的websocket的服务器端的功能
	8，握手协议， 检测第一次建立连接的时候 授信
*/

func main() {
	host := config.File.MustValue("gate_server", "host", "localhost")
	port := config.File.MustValue("gate_server", "port", "8004")
	//fmt.Println("host:", host)
	//fmt.Println("port:", port)
	s := net.NewServer(host + ":" + port)
	s.NeedSecret(true)
	gate.Init()
	s.Router(gate.Router)
	s.Start()

}
