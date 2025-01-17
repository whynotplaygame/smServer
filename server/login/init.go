package login

import (
	"smServer/db"
	"smServer/net"
	"smServer/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	//数据库初始化
	db.TestDb()
	// 还有别的初始化方法
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}

//func Init() {
//
//}
