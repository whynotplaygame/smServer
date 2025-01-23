package main

import (
	"smServer/config"
	"smServer/net"
	"smServer/server/game"
)

/*
	1,登录完了，创建角色
	2，需要根据用户查询此用户拥有的角色，没有 创建角色
	3，木材，铁，令牌，金钱，主城，武将，这些数据，要不要初始化，已经玩过游戏，这些数值是不是说需要查询
	4，地图，城池，资源土地，要塞，需要定义
	5，资源，军队，城池，武将等



*/

func main() {
	host := config.File.MustValue("game_server", "host", "127.0.0.1")
	port := config.File.MustValue("game_server", "port", "8001")
	//fmt.Println("host:", host)
	//fmt.Println("port:", port)
	s := net.NewServer(host + ":" + port)
	s.NeedSecret(false)
	game.Init()
	s.Router(game.Router)
	s.Start()
}
