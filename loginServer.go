package main

import (
	"smServer/config"
	"smServer/net"
	"smServer/server/login"
)

// http://localhost:8080/api/login
// localhost:8080 服务器   /api/login 路由
// webscoket: ws://localhost:8080  服务器   发消息  封装为路由

func main() {
	host := config.File.MustValue("login_server", "host", "localhost")
	port := config.File.MustValue("login_server", "port", "8003")
	//fmt.Println("host:", host)
	//fmt.Println("port:", port)
	s := net.NewServer(host + ":" + port)
	s.NeedSecret(false)
	login.Init()
	s.Router(login.Router)
	s.Start()

}
