package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"smServer/utils"
	"time"
)

var RoleService = &roleService{}

type roleService struct {
}

func (service *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, req *net.WsMsgReq) error {

	role := &data.Role{}

	session := db.Engin.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Println("事务开启失败", err)
		return err
	}
	req.Context.Set("dbSession", session)

	ok, err := db.Engin.Table(role).Where("uid = ?", uid).Get(role)
	if err != nil {
		log.Println("查询用户出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	if ok {
		rid := role.RId

		roleRes := &data.RoleRes{}
		ok, err = db.Engin.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
		if err != nil {
			log.Println("查询角色资源出错", err)
			return common.New(constant.DBError, "数据库出错")
		}
		if !ok {
			roleRes.RId = rid
			roleRes.Gold = gameConfig.Base.Role.Gold
			roleRes.Decree = gameConfig.Base.Role.Decree
			roleRes.Grain = gameConfig.Base.Role.Grain
			roleRes.Iron = gameConfig.Base.Role.Iron
			roleRes.Stone = gameConfig.Base.Role.Stone
			roleRes.Wood = gameConfig.Base.Role.Wood
			_, err = session.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("插入角色资源出错", err)
				return common.New(constant.DBError, "数据库出错")
			}
		}

		// 数据库操作 与 业务操作，分割，所以，数据库的role等需要一个toModel转换功能
		rsp.RoleRes = roleRes.ToModel().(model.RoleRes) //
		rsp.Role = role.ToModel().(model.Role)
		rsp.Time = time.Now().UnixNano() / 1e6
		token, _ := utils.Award(rid) // 利用rid 针对于角色生成token
		rsp.Token = token
		req.Conn.SetProperty("role", role) // 方便使用
		// 初始化玩家属性
		if err := RoleAttrService.TryCreate(rid, req); err != nil {
			session.Rollback()
			return common.New(constant.DBError, "数据库出错")
		}

		// 初始化城池
		if err := RoleCityService.InitCity(rid, role.NickName, req); err != nil {
			session.Rollback()
			return common.New(constant.DBError, "数据库出错")
		}
	} else {
		log.Println("角色不存在", err)
		return common.New(constant.RoleNotExist, "数据库出错")
	}

	if err := session.Commit(); err != nil {
		log.Println("事务提交失败", err)
		return common.New(constant.RoleNotExist, "数据库出错")
	}
	return nil
}

func (service *roleService) GetRoleRes(rid int) (model.RoleRes, error) {
	roleRes := &data.RoleRes{}
	ok, err := db.Engin.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return model.RoleRes{}, common.New(constant.DBError, "数据库出错")
	}
	if ok {
		return roleRes.ToModel().(model.RoleRes), nil
	}
	return model.RoleRes{}, common.New(constant.DBError, "角色资源不存在")
}

func (service *roleService) Get(rid int) *data.Role {
	role := &data.Role{}
	ok, err := db.Engin.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return nil
	}
	if ok {
		return role
	}
	return nil
}

func (service *roleService) GetRoleNickName(rid int) string {
	role := &data.Role{}
	ok, err := db.Engin.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return ""
	}
	if ok {
		return role.NickName
	}
	return ""
}
