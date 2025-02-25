package logic

import (
	"fmt"
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
	"smServer/utils"
	"sync"
	"time"
	"xorm.io/xorm"
)

var RoleCityService = &roleCityService{
	dbRB:   make(map[int]*data.MapRoleCity),
	posRC:  make(map[int]*data.MapRoleCity),
	roleRC: make(map[int][]*data.MapRoleCity),
}

type roleCityService struct {
	mutex sync.RWMutex
	dbRB  map[int]*data.MapRoleCity //缓存
	// 位置 key posId
	posRC map[int]*data.MapRoleCity
	// key 角色id
	roleRC map[int][]*data.MapRoleCity //RB 是RoleBuild的意思，key是角色ID
}

func (service *roleCityService) Load() {
	// 查询所有的角色城池
	db.Engin.Find(service.dbRB)

	for _, v := range service.dbRB {
		posId := global.ToPosition(v.X, v.Y)
		service.posRC[posId] = v
		_, ok := service.roleRC[v.RId]
		if !ok {
			service.roleRC[v.RId] = make([]*data.MapRoleCity, 0)
		}
		service.roleRC[v.RId] = append(service.roleRC[v.RId], v)

	}
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
			if service.isCanBuild(roleCity.X, roleCity.Y) {
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
				// 创建新城池后，更新内存中的城池列表，
				posId := global.ToPosition(roleCity.X, roleCity.Y)
				service.posRC[posId] = roleCity
				_, ok := service.roleRC[rid]
				if !ok {
					service.roleRC[rid] = make([]*data.MapRoleCity, 0)
				} else {
					service.roleRC[rid] = append(service.roleRC[rid], roleCity)
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

func (service *roleCityService) isCanBuild(x int, y int) bool {
	confs := gameConfig.MapRes.Confs
	pIndex := global.ToPosition(x, y)

	_, ok := confs[pIndex]
	if !ok {
		return false
	}
	sysbuild := gameConfig.MapRes.SysBuild

	// 城池 1 范围内不能超过边界
	if x+1 >= global.MapWith || y+1 >= global.MapHeight || x-1 < 0 || y-1 < 0 {
		return false
	}

	// 系统城池5格内，不嫩创建玩家城池
	for _, v := range sysbuild {
		if v.Type == gameConfig.MapBuildSysCity {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				y >= v.Y-5 &&
				y <= v.Y+5 { // 5格之内
				return false
			}
		}
	}

	// 玩家城池5格内
	for i := x - 5; i <= x+5; i++ {
		for j := y - 5; j <= y+5; j++ {
			posId := global.ToPosition(i, j)
			_, ok := service.posRC[posId]
			if ok {
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

func (service *roleCityService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleCity, error) {
	x := req.X
	y := req.Y
	lenght := req.Length
	mrcs := make([]model.MapRoleCity, 0)
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return mrcs, nil
	}
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	maxX := utils.MinInt(global.MapWith, x+lenght-1)
	maxY := utils.MinInt(global.MapHeight, y+lenght-1)

	// 范围 x-lenght, x + length  y-length  y + length
	for i := x - lenght; i <= maxX; i++ {
		for j := y - lenght; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			// fmt.Println("扫描建筑posId:", posId)
			mrc, ok := service.posRC[posId]
			if ok {
				fmt.Println("扫描城池posId:", posId)
				mrcs = append(mrcs, mrc.ToModel().(model.MapRoleCity))
			}
		}
	}

	return mrcs, nil
}

func (service *roleCityService) GetByRId(rid int) ([]*data.MapRoleCity, bool) {
	service.mutex.RLock()
	r, ok := service.roleRC[rid]
	service.mutex.RUnlock()
	return r, ok
}

func (service *roleCityService) GetMainCity(rid int) *data.MapRoleCity {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	rcs, ok := service.roleRC[rid]
	if ok {
		for _, v := range rcs {
			if v.IsMain == 1 {
				return v
			}
		}
	}
	return nil
}
