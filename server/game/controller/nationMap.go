package controller

import (
	"github.com/mitchellh/mapstructure"
	"smServer/constant"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/logic"
	"smServer/server/game/middleware"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var DefaultNationMapController = &nationMapController{}

type nationMapController struct {
}

func (r *nationMapController) Router(router *net.Router) {
	g := router.Group("nationMap")
	g.Use(middleware.Log()) // 为所有添加了日志的中间件
	g.AddRouter("config", r.config)
	g.AddRouter("scanBlock", r.scanBlock, middleware.CheckRole())
}

func (r *nationMapController) config(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// reqObj := &model.ConfigReq{}
	rspObj := &model.ConfigRsp{}

	cfgs := gameConfig.MapBuildConf.Cfg
	rspObj.Confs = make([]model.Conf, len(cfgs))

	for index, v := range cfgs {
		rspObj.Confs[index].Type = v.Type
		rspObj.Confs[index].Name = v.Name
		rspObj.Confs[index].Level = v.Level
		rspObj.Confs[index].Defender = v.Defender
		rspObj.Confs[index].Durable = v.Durable
		rspObj.Confs[index].Grain = v.Grain
		rspObj.Confs[index].Iron = v.Iron
		rspObj.Confs[index].Stone = v.Stone
		rspObj.Confs[index].Wood = v.Wood
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

}

func (r *nationMapController) scanBlock(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ScanBlockReq{}
	rspObj := &model.ScanRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK

	// 扫描角色建筑
	mrb, err := logic.RoleBuildService.ScanBlock(reqObj)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.MRBuilds = mrb

	// 扫描角色城池
	mrc, err := logic.RoleCityService.ScanBlock(reqObj)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.MCBuilds = mrc

	role, _ := req.Conn.GetProperty("role")
	rl := role.(*data.Role)
	// 扫描角色军队
	mra, err := logic.ArmyService.ScanBlock(reqObj, rl.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Armys = mra
	rsp.Body.Msg = rspObj
}
