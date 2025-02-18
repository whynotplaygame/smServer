package game

import (
	"smServer/db"
	"smServer/net"
	"smServer/server/game/controller"
	"smServer/server/game/gameConfig"
	"smServer/server/game/gameConfig/general"
	"smServer/server/game/logic"
)

var Router = &net.Router{}

func Init() {
	db.TestDb() // 初始化数据库
	//加载基础配置
	gameConfig.Base.Load()
	// 加载地图的资源配置
	gameConfig.MapBuildConf.Load()
	// 加载地图单元格配置
	gameConfig.MapRes.Load()
	// 加载城池设施配置
	gameConfig.FaiclityConfig.Load()
	// 加载武将配置
	general.General.Load()
	// 加载技能配置信息
	gameConfig.Skill.Load()
	// 加载所有的建筑信息
	logic.RoleBuildService.Load()
	// 加载所有的城池信息
	logic.RoleCityService.Load()

	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
	controller.DefaultGeneralController.Router(Router)
	controller.DefaultArmyController.Router(Router)
	controller.WarController.Router(Router)
	controller.SkillController.Router(Router)
}
