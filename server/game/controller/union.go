package controller

import (
	"smServer/constant"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/logic"
	"smServer/server/game/middleware"
	"smServer/server/game/model"
)

var UnionController = &unionController{}

type unionController struct {
}

func (u *unionController) Router(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list", u.list, middleware.CheckRole())
}

func (u *unionController) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.ListRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	uns, err := logic.CoalitionService.List()
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = uns
	rsp.Body.Msg = rspObj
}
