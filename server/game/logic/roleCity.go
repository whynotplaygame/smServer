package logic

import (
	"log"
	"math/rand"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/global"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"time"
	"xorm.io/xorm"
)

var RoleCityService = &roleCityService{}

type roleCityService struct {
}

func (service *roleCityService) InitCity(rid int, name string, req *net.WsMsgReq) error {
	// 根据rid查询，如果没有就创建
	roleCity := &data.MapRoleCity{}
	ok, err := db.Engin.Table(roleCity).Where("rid = ?", rid).Get(roleCity)
	if err != nil {
		log.Println("查询角色城池出错", err)
		return common.New(constant.DBError, "数据库出错")
	}

	if ok {
		return nil
	} else {

		for {
			// 初始化
			roleCity.X = rand.Intn(global.MapWith)
			roleCity.Y = rand.Intn(global.MapHeight)
			// 这个城池能不能在这个坐标点创建，比如，需要判断 五格之内，不能有玩家的城池
			// done
			if isCanBuild(roleCity.X, roleCity.Y) {
				roleCity.RId = rid
				roleCity.Name = name
				roleCity.CurDurable = gameConfig.Base.City.Durable
				roleCity.CreatedAt = time.Now()
				roleCity.IsMain = 1
				if session := req.Context.Get("dbSession"); session != nil {
					_, err = session.(*xorm.Session).Table(roleCity).Insert(roleCity)
				} else {
					_, err = db.Engin.Table(roleCity).Insert(roleCity)
				}

				if err != nil {
					log.Println("插入角色城池出错", err)
					return common.New(constant.DBError, "数据库出错")
				}
				// 初始化城池的设施
				// done
				if err := CityFacilityService.TryCreate(roleCity.CityId, rid, req); err != nil {
					log.Println("插入角色城池设施出错", err)
					return common.New(err.(*common.MyError).Code(), err.Error())
				}
				break
			}

		}
	}
	return nil
}

func isCanBuild(x int, y int) bool {
	confs := gameConfig.MapRes.Confs
	pIndex := global.ToPosition(x, y)

	_, ok := confs[pIndex]
	if !ok {
		return false
	}
	sysbuild := gameConfig.MapRes.SysBuild
	// 系统城池5格内，不嫩创建玩家城池
	for _, v := range sysbuild {
		if v.Type == gameConfig.MapBuildSysCity {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				v.Y >= v.Y-5 &&
				y <= v.Y+5 { // 5格之内
				return false
			}
		}
	}
	return true
}

func (service *roleCityService) GetRoleCitys(id int) ([]model.MapRoleCity, error) {
	citys := make([]data.MapRoleCity, 0)
	city := &data.MapRoleCity{}
	err := db.Engin.Table(city).Where("rid = ?", id).Find(&citys)

	//这不能给容量 要不然 结果出错了
	modelCitys := make([]model.MapRoleCity, 0)
	if err != nil {
		log.Println("查询玩家城池出错", err)
		return modelCitys, err
	}
	for _, v := range citys {
		modelCitys = append(modelCitys, v.ToModel().(model.MapRoleCity))
	}
	return modelCitys, nil
}
