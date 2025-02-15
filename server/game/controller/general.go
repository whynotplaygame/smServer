package controller

import (
	"smServer/constant"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/logic"
	"smServer/server/game/middleware"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var DefaultGeneralController = &generalController{}

type generalController struct {
}

func (r *generalController) Router(router *net.Router) {
	g := router.Group("general")
	g.Use(middleware.Log()) // 为所有添加了日志的中间件
	g.AddRouter("myGenerals", r.myGenerals, middleware.CheckRole())
}

func (r *generalController) myGenerals(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 查询武将的时候，角色拥有的物件，查询出来即可
	// 如果初始化，进入游戏，武将没有，需要随机三个武将，很多游戏，初始化武将是一样的

	rspObj := &model.MyGeneralRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId

	gs, err := logic.GeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rspObj.Generals = gs
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
