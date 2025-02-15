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

var WarController = &warController{}

type warController struct {
}

func (r *warController) Router(router *net.Router) {
	g := router.Group("war")
	g.Use(middleware.Log()) // 为所有添加了日志的中间件
	g.AddRouter("report", r.report, middleware.CheckRole())
}

func (w *warController) report(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 查找战报表，得出数据
	rspObj := &model.WarReportRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId
	wrs, err := logic.WarService.GetWarReport(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = wrs
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
