package logic

import (
	"log"
	"smServer/db"
	"smServer/server/game/gameConfig"
	"smServer/server/game/model/data"
)

var RoleResService = &roleResService{}

type roleResService struct {
}

func (service *roleResService) GetRoleRes(rid int) *data.RoleRes {
	roleRes := &data.RoleRes{}
	ok, err := db.Engin.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return nil
	}
	if ok {
		return roleRes
	}
	return nil
}

func (service *roleResService) GetYield(rid int) data.Yield {
	// 基础产量 + 城池设施的产量 + 建筑产量
	rbYield := RoleBuildService.getYield(rid)
	cfYield := CityFacilityService.GetYield(rid)
	var y data.Yield
	y.Gold = rbYield.Gold + cfYield.Gold + gameConfig.Base.Role.GoldYield
	y.Stone = rbYield.Stone + cfYield.Stone + gameConfig.Base.Role.StoneYield
	y.Iron = rbYield.Iron + cfYield.Iron + gameConfig.Base.Role.IronYield
	y.Grain = rbYield.Grain + cfYield.Grain + gameConfig.Base.Role.GrainYield
	y.Wood = rbYield.Wood + cfYield.Wood + gameConfig.Base.Role.WoodYield

	return y

}
