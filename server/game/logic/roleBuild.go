package logic

import (
	"fmt"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/global"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"smServer/utils"
	"sync"
)

var RoleBuildService = &roleBuildService{
	posRB:  make(map[int]*data.MapRoleBuild),
	roleRB: make(map[int][]*data.MapRoleBuild),
}

type roleBuildService struct {
	mutex sync.RWMutex

	// 位置 key posId
	posRB map[int]*data.MapRoleBuild
	// key 角色id
	roleRB map[int][]*data.MapRoleBuild //RB 是RoleBuild的意思，key是角色ID
}

func (service *roleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild, error) {
	mrs := make([]data.MapRoleBuild, 0)
	mr := &data.MapRoleBuild{}
	err := db.Engin.Table(mr).Where("rid = ?", rid).Find(&mrs)

	if err != nil {
		log.Println("查询玩家建筑出错", err)
		return nil, common.New(constant.DBError, "建筑查询出错")
	}
	modelMrs := make([]model.MapRoleBuild, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.MapRoleBuild))
	}
	return modelMrs, nil
}

func (service *roleBuildService) Load() {
	// 加载系统的建筑以及玩家的建筑
	// 首先需要判断数据库，是否保存了系统饿建筑，没有，进行一个保存
	total, err := db.Engin.
		Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).
		Count(new(data.MapRoleBuild))
	if err != nil {
		log.Println("加载系统建筑")
		panic(err) //
		// return 初始的加载，必须挺掉
	}
	if int64(len(gameConfig.MapRes.SysBuild)) != total {
		// 证明数据存储的系统建筑有问题
		db.Engin.
			Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).
			Delete(new(data.MapRoleBuild))

		for _, v := range gameConfig.MapRes.SysBuild {
			build := data.MapRoleBuild{
				RId:   0,
				Type:  v.Type,
				Level: v.Level,
				X:     v.X,
				Y:     v.Y,
			}
			build.Init()
			db.Engin.InsertOne(&build)
		}
	}

	// 查询所有的角色建筑
	dbPB := make(map[int]*data.MapRoleBuild)
	db.Engin.Find(dbPB)

	for _, v := range dbPB {
		posId := global.ToPosition(v.X, v.Y)
		service.posRB[posId] = v
		_, ok := service.roleRB[v.RId]
		if !ok {
			service.roleRB[v.RId] = make([]*data.MapRoleBuild, 0)
		} else {
			service.roleRB[v.RId] = append(service.roleRB[v.RId], v)
		}
	}
}

func (service *roleBuildService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleBuild, error) {

	x := req.X
	y := req.Y
	length := req.Length
	mrbs := make([]model.MapRoleBuild, 0)
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return mrbs, nil
	}
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	// 范围 x-lenght, x + length  y-length  y + length
	for i := x - length; i <= maxX; i++ {
		for j := y - length; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			// fmt.Println("扫描建筑posId:", posId)
			mrb, ok := service.posRB[posId]
			if ok {
				fmt.Println("扫描建筑posId:", posId)
				mrbs = append(mrbs, mrb.ToModel().(model.MapRoleBuild))
			}
		}
	}

	return mrbs, nil
}

func (service *roleBuildService) getYield(rid int) data.Yield {
	service.mutex.RLock()
	defer service.mutex.RUnlock()

	rbs, ok := service.roleRB[rid]
	var y data.Yield
	if ok {
		for _, v := range rbs {
			y.Iron += v.Iron
			y.Wood += v.Wood
			y.Grain += v.Grain
			y.Stone += v.Stone
		}
	}
	return y
}
