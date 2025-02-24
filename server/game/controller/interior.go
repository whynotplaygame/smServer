package controller

import (
	"smServer/constant"
	"smServer/net"
	"smServer/server/game/gameConfig"
	"smServer/server/game/logic"
	"smServer/server/game/middleware"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"time"
)

var InteriorCotroller = &interiorController{}

type interiorController struct {
}

func (i *interiorController) Router(router *net.Router) {
	g := router.Group("interior")
	g.Use(middleware.Log())
	g.AddRouter("openCollect", i.openCollect, middleware.CheckRole())
	g.AddRouter("collect", i.collect, middleware.CheckRole())
}

func (i *interiorController) openCollect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// reqObj := &model.OpenCollectionReq{}
	rspObj := &model.OpenCollectionRsp{}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK

	rsp.Body.Msg = rspObj
	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)
	ra := logic.RoleAttrService.Get(role.RId)
	if ra != nil {
		rspObj.CurTimes = ra.CollectTimes
		rspObj.Limit = gameConfig.Base.Role.CollectTimesLimit
		// 征收间隔时间
		interval := gameConfig.Base.Role.CollectInterval
		if ra.LastCollectTime.IsZero() {
			rspObj.NextTime = 0
		} else {
			if rspObj.CurTimes >= rspObj.Limit {
				// 今天已经完成征收了，下次征收就是第二天,以最后一次征收时间为准
				// 第二天，从0点开始
				y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
				ti := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)) //东8区
				rspObj.NextTime = ti.UnixNano() / 1e6
			} else {
				ti := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
				rspObj.NextTime = ti.UnixNano() / 1e6
			}
		}

	}
}

func (i *interiorController) collect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 查询角色资源，得到当前的金币
	// 查询角色属性，获取征收的相关信息
	// 产讯获取当前的产量，征收的金币是多少
	// reqObj := &model.OpenCollectionReq{}
	rspObj := &model.CollectionRsp{}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK

	rsp.Body.Msg = rspObj

	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)
	// 属性
	ra := logic.RoleAttrService.Get(role.RId)
	if ra == nil {
		rsp.Body.Code = constant.DBError
		return
	}

	// 角色资源
	rs := logic.RoleResService.GetRoleRes(role.RId)
	if rs == nil {
		rsp.Body.Code = constant.DBError
		return
	}
	// 产量
	ry := logic.RoleResService.GetYield(role.RId)
	rs.Gold += ry.Gold
	// go channel  一旦需要更新，发一个需要更新的信号，接收方，消费方，接收小心，进行更新
	rs.SyncExecute()

	rspObj.Gold = ry.Gold

	// 计算征收
	curTime := time.Now()
	limit := gameConfig.Base.Role.CollectTimesLimit
	interval := gameConfig.Base.Role.CollectInterval
	lastTime := ra.LastCollectTime

	if curTime.YearDay() != lastTime.YearDay() || curTime.Year() != lastTime.Year() {
		ra.CollectTimes = 0
		ra.LastCollectTime = time.Time{}
	}

	ra.CollectTimes += 1
	ra.LastCollectTime = curTime
	// todo 进行数据库的更新操作
	ra.SyncExecute()

	rspObj.Limit = limit
	rspObj.CurTimes = ra.CollectTimes

	if rspObj.CurTimes >= rspObj.Limit {
		// 今天已经完成征收了，下次征收就是第二天,以最后一次征收时间为准
		// 第二天，从0点开始
		y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
		ti := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)) //东8区
		rspObj.NextTime = ti.UnixNano() / 1e6
	} else {
		ti := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
		rspObj.NextTime = ti.UnixNano() / 1e6
	}

}
