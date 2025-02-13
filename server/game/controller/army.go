package controller

import (
	"github.com/mitchellh/mapstructure"
	"smServer/constant"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/logic"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var DefaultArmyController = &armyController{}

type armyController struct {
}

func (r *armyController) Router(router *net.Router) {
	g := router.Group("army")
	g.AddRouter("myList", r.myList)
}

func (a *armyController) myList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ArmyListReq{}
	rspObj := &model.ArmyListRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId
	ams, err := logic.ArmyService.GetArmyByCity(rid, reqObj.CityId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Armys = ams
	rspObj.CityId = reqObj.CityId
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
