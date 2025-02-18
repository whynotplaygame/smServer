package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/middleware"
	"smServer/server/game/model/data"
	"time"

	"smServer/server/game/logic"
	"smServer/server/game/model"

	"smServer/utils"
)

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.Use(middleware.Log()) // 为所有添加了日志的中间件
	g.AddRouter("create", r.create)
	g.AddRouter("enterServer", r.enterServer)
	g.AddRouter("myProperty", r.myProperty, middleware.CheckRole())
	g.AddRouter("posTagList", r.posTagList)

}

func (r *RoleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 进去游戏
	// session 需要验证是否合法，合法情况喜爱，可以取出登录的用户id
	// 根据用户的id 去查询对应的游戏角色，如果有，继续，如果没有，提示无角色
	// 根据角色ID 查询角色拥有的资源 roleRes，如果有，返回，如果没有，初始化
	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		log.Println("参数错误：", err.Error())
		rsp.Body.Code = constant.InvalidParam
		return
	}
	session := reqObj.Session
	_, claim, err := utils.ParseToken(session)
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	uid := claim.Uid // 获取uid

	if err = logic.RoleService.EnterServer(uid, rspObj, req); err != nil {
		rspObj.Time = time.Now().Unix() / 1e6
		rsp.Body.Msg = rspObj
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

}

func (r *RoleController) myProperty(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 根据角色的id,去查询 军队，资源，建筑，城池，武将
	ro, err := req.Conn.GetProperty("role") // 这个属性是 enterGame时，setProperty
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role := ro.(*data.Role)
	rspObj := &model.MyRolePropertyRsp{}

	// 资源
	rspObj.RoleRes, err = logic.RoleService.GetRoleRes(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	// 城池
	rspObj.Citys, err = logic.RoleCityService.GetRoleCitys(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 建筑
	rspObj.MRBuilds, err = logic.RoleBuildService.GetBuilds(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 军队
	rspObj.Armys, err = logic.ArmyService.GetArmys(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	// 武将
	rspObj.Generals, err = logic.GeneralService.GetGenerals(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (r *RoleController) posTagList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.PosTagListRsp{}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	// 去角色属性表查询
	ro, err := req.Conn.GetProperty("role") // 这个属性是 enterGame时，setProperty
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := ro.(*data.Role).RId
	pts, err := logic.RoleAttrService.GetTagList(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.PosTags = pts
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (r *RoleController) create(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.CreateRoleReq{}
	rspObj := &model.CreateRoleRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		rsp.Body.Code = constant.InvalidParam
	}
	role := &data.Role{}
	ok, err := db.Engin.Where("uid=?", reqObj.UId).Get(role)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	if ok { // 已存在
		rsp.Body.Code = constant.RoleAlreadyCreate
		return
	}
	// 创建角色
	role.UId = reqObj.UId
	role.Sex = reqObj.Sex
	role.NickName = reqObj.NickName
	role.Balance = 0
	role.HeadId = reqObj.HeadId
	role.CreatedAt = time.Now()
	role.LoginTime = time.Now()

	_, err = db.Engin.Insert(role)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	rspObj.Role = role.ToModel().(model.Role)
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
