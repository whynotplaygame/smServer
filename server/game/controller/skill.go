package controller

import (
	"smServer/constant"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/logic"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var SkillController = &skillController{}

type skillController struct {
}

func (r *skillController) Router(router *net.Router) {
	g := router.Group("skill")
	g.AddRouter("list", r.list)
}

func (w *skillController) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 查找战报表，得出数据
	rspObj := &model.SkillListRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId
	skills, err := logic.SkillService.GetSkills(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = skills
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
