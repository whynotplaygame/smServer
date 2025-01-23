package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/game/gameConfig"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"smServer/utils"
	"time"
)

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.AddRouter("enterServer", r.enterServer)
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
	role := &data.Role{}
	ok, err := db.Engin.Table(role).Where("uid = ?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		rsp.Body.Code = constant.DBError
		return
	}

	if ok {
		rsp.Body.Code = constant.OK
		rsp.Body.Msg = rspObj

		rid := role.RId

		roleRes := &data.RoleRes{}
		ok, err = db.Engin.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
		if err != nil {
			log.Println("查询角色资源出错", err)
			rsp.Body.Code = constant.DBError
			return
		}
		if !ok {
			roleRes.RId = rid
			roleRes.Gold = gameConfig.Base.Role.Gold
			roleRes.Decree = gameConfig.Base.Role.Decree
			roleRes.Grain = gameConfig.Base.Role.Grain
			roleRes.Iron = gameConfig.Base.Role.Iron
			roleRes.Stone = gameConfig.Base.Role.Stone
			roleRes.Wood = gameConfig.Base.Role.Wood
			_, err = db.Engin.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("插入角色资源出错", err)
				rsp.Body.Code = constant.DBError
				return
			}
		} else {

		}
		// 数据库操作 与 业务操作，分割，所以，数据库的role等需要一个toModel转换功能
		rspObj.RoleRes = roleRes.ToModel().(model.RoleRes) //
		rspObj.Role = role.ToModel().(model.Role)
		rspObj.Time = time.Now().UnixNano() / 1e6
		token, _ := utils.Award(rid) // 利用rid 针对于角色生成token
		rspObj.Token = token

	} else {
		rsp.Body.Code = constant.RoleNotExist
		return
	}
}
