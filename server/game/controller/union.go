package controller

import (
	"github.com/mitchellh/mapstructure"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/logic"
	"smServer/server/game/middleware"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"time"
)

var UnionController = &unionController{}

type unionController struct {
}

func (u *unionController) Router(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list", u.list, middleware.CheckRole())
	g.AddRouter("info", u.info, middleware.CheckRole())
	g.AddRouter("applyList", u.applyList, middleware.CheckRole())
	g.AddRouter("create", u.create, middleware.CheckRole())
	g.AddRouter("join", u.join, middleware.CheckRole())
	g.AddRouter("verify", u.verify, middleware.CheckRole())
	g.AddRouter("member", u.member, middleware.CheckRole())
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

func (u *unionController) info(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.InfoRsp{}
	rspObj := &model.InfoRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK

	un, err := logic.CoalitionService.Get(reqObj.Id)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
	}
	rspObj.Info = un
	rspObj.Id = un.Id
	rsp.Body.Msg = rspObj
}

func (u *unionController) applyList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 根据联盟id ，去查询申请列表，rid 申请人，你角色表，查询详情即可
	// state 0 正在申请 1 拒绝 2 同意
	// 什么能看到申请列表，只有盟主和 副盟主能看到
	reqObj := &model.ApplyReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.ApplyRsp{}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)

	// 查询联盟
	un := logic.CoalitionService.GetCoalition(reqObj.Id)
	if un == nil {
		rsp.Body.Code = constant.DBError
		return
	}
	if un.Chairman != role.RId && un.ViceChairman != role.RId {
		rspObj.Id = un.Id
		rspObj.Applys = make([]model.ApplyItem, 0)
		return
	}

	ais, err := logic.CoalitionService.GetListApply(reqObj.Id, 0)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	rspObj.Applys = ais
	rspObj.Id = reqObj.Id
}

func (u *unionController) create(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.CreateReq{}
	rspObj := &model.CreateRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)

	has := logic.RoleAttrService.IsHasUnion(role.RId)
	if has {
		rsp.Body.Code = constant.UnionAlreadyHas
		return
	}

	c, ok := logic.CoalitionService.Create(reqObj.Name, role.RId)
	if ok {
		rspObj.Id = c.Id
		logic.CoalitionService.MemberEnter(role.RId, c.Id)
		logic.CoalitionService.NewCreateLog(role.NickName, c.Id, role.RId)
	} else {
		rsp.Body.Code = constant.UnionCreateError
	}

}

func (u *unionController) join(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.JoinReq{}
	rspObj := &model.JoinRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)

	has := logic.RoleAttrService.IsHasUnion(role.RId)
	if has {
		rsp.Body.Code = constant.UnionAlreadyHas
		return
	}

	union := logic.CoalitionService.GetById(reqObj.Id)
	if union == nil {
		rsp.Body.Code = constant.UnionNotFound
		return
	}
	if len(union.MemberArray) >= gameConfig.Base.Union.MemberLimit {
		rsp.Body.Code = constant.PeopleIsFull
		return
	}

	//判断当前是否已经有申请
	has, _ = db.Engin.Table(data.CoalitionApply{}).Where(
		"union_id=? and state=? and rid=?",
		reqObj.Id, model.UnionUntreated, role.RId).Get(&data.CoalitionApply{})
	if has {
		rsp.Body.Code = constant.HasApply
		return
	}

	//写入申请列表
	apply := &data.CoalitionApply{
		RId:     role.RId,
		UnionId: reqObj.Id,
		Ctime:   time.Now(),
		State:   model.UnionUntreated}

	_, err := db.Engin.InsertOne(apply)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}

	//推送主、副盟主
	apply.SyncExecute() // todo
}

func (u *unionController) verify(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.VerifyReq{}
	rspObj := &model.VerifyRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	rsp.Body.Msg = rspObj
	rsp.Body.Seq = req.Body.Seq
	rspObj.Id = reqObj.Id
	rsp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)

	apply := &data.CoalitionApply{}
	ok, err := db.Engin.Table(new(data.CoalitionApply)).Where(
		"id=? and state=?", reqObj.Id, model.UnionUntreated).Get(apply)
	if ok && err == nil {
		targetRole := logic.RoleService.Get(apply.RId)
		if targetRole == nil {
			rsp.Body.Code = constant.RoleNotExist
			return
		}

		if u := logic.CoalitionService.GetById(apply.UnionId); u != nil {

			if u.Chairman != role.RId && u.ViceChairman != role.RId {
				rsp.Body.Code = constant.PermissionDenied
				return
			}

			if len(u.MemberArray) >= gameConfig.Base.Union.MemberLimit {
				rsp.Body.Code = constant.PeopleIsFull
				return
			}

			if ok := logic.RoleAttrService.IsHasUnion(apply.RId); ok {
				rsp.Body.Code = constant.UnionAlreadyHas
			} else {
				if reqObj.Decide == model.UnionAdopt {
					//同意
					c := logic.CoalitionService.GetById(apply.UnionId)
					if c != nil {
						c.MemberArray = append(c.MemberArray, apply.RId)
						logic.CoalitionService.MemberEnter(apply.RId, apply.UnionId)
						c.SyncExecute()
						logic.CoalitionService.NewJoin(targetRole.NickName, apply.UnionId, role.RId, apply.RId)
					}
				}
			}
			apply.State = reqObj.Decide
			db.Engin.Table(apply).ID(apply.Id).Cols("state").Update(apply)
		} else {
			rsp.Body.Code = constant.UnionNotFound
			return
		}

	} else {
		rsp.Body.Code = constant.InvalidParam
	}
}

func (u *unionController) member(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.MemberReq{}
	rspObj := &model.MemberRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	rsp.Body.Msg = rspObj
	rspObj.Id = reqObj.Id
	rsp.Body.Code = constant.OK

	union := logic.CoalitionService.GetById(reqObj.Id)
	if union == nil {
		rsp.Body.Code = constant.UnionNotFound
		return
	}

	rspObj.Members = make([]model.Member, 0)
	for _, rid := range union.MemberArray {
		if role := logic.RoleService.Get(rid); role != nil {
			m := model.Member{RId: role.RId, Name: role.NickName}
			if main := logic.RoleCityService.GetMainCity(role.RId); main != nil {
				m.X = main.X
				m.Y = main.Y
			}

			if rid == union.Chairman {
				m.Title = model.UnionChairman
			} else if rid == union.ViceChairman {
				m.Title = model.UnionViceChairman
			} else {
				m.Title = model.UnionCommon
			}
			rspObj.Members = append(rspObj.Members, m)
		}
	}
}
