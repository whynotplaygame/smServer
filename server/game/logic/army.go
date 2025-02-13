package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var ArmyService = &armyService{}

type armyService struct {
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
