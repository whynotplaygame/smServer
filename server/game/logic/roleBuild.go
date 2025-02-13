package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var RoleBuildService = &roleBuildService{}

type roleBuildService struct {
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
