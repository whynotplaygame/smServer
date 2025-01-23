package game

import (
	"smServer/db"
	"smServer/net"
	"smServer/server/game/controller"
	"smServer/server/game/gameConfig"
)

var Router = &net.Router{}

func Init() {
	db.TestDb() // 初始化数据库
	//加载基础配置
	gameConfig.Base.Load()
	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
}
