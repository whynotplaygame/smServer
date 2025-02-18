package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/global"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"smServer/utils"
	"sync"
)

var ArmyService = &armyService{
	passByPosArmys: make(map[int]map[int]*data.Army),
}

type armyService struct {
	passBy         sync.RWMutex
	passByPosArmys map[int]map[int]*data.Army
}

func (service *armyService) GetArmys(rid int) ([]model.Army, error) {
	ams := make([]data.Army, 0)
	am := &data.Army{}
	err := db.Engin.Table(am).Where("rid = ?", rid).Find(&ams)

	if err != nil {
		log.Println("查询玩家军队出错", err)
		return nil, common.New(constant.DBError, "查询玩家军队出错")
	}
	modelArmys := make([]model.Army, 0)
	for _, v := range ams {
		modelArmys = append(modelArmys, v.ToModel().(model.Army))
	}
	return modelArmys, nil
}

func (service *armyService) GetArmyByCity(rid, cid int) ([]model.Army, error) {
	ams := make([]data.Army, 0)
	am := &data.Army{}
	err := db.Engin.Table(am).Where("rid = ? and cityId = ?", rid, cid).Find(&ams)

	if err != nil {
		log.Println("查询玩家军队出错", err)
		return nil, common.New(constant.DBError, "查询玩家军队出错")
	}
	modelArmys := make([]model.Army, 0)
	for _, v := range ams {
		modelArmys = append(modelArmys, v.ToModel().(model.Army))
	}
	return modelArmys, nil
}

func (service *armyService) ScanBlock(req *model.ScanBlockReq, roleId int) ([]model.Army, error) {
	x := req.X
	y := req.Y
	length := req.Length
	out := make([]model.Army, 0)
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return out, nil
	}

	service.passBy.RLock()
	defer service.passBy.RUnlock()
	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {

			posId := global.ToPosition(i, j)
			armys, ok := service.passByPosArmys[posId]
			if ok {
				//是否在视野范围内
				is := armyIsInView(roleId, i, j)
				if is == false {
					continue
				}
				for _, army := range armys {
					out = append(out, army.ToModel().(model.Army))
				}
			}
		}
	}
	return out, nil

}

func armyIsInView(rid, x, y int) bool {
	//简单点 先设为true
	return true
}
