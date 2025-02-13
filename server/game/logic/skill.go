package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var SkillService = &skillService{}

type skillService struct {
}

func (w *skillService) GetSkills(rid int) ([]model.Skill, error) {
	sks := make([]data.Skill, 0)
	sk := &data.Skill{}
	err := db.Engin.Table(sk).Where("rid = ?", rid).Find(&sks)

	if err != nil {
		log.Println("查询技能出错", err)
		return nil, common.New(constant.DBError, "查询玩家技能出错")
	}
	modelSkills := make([]model.Skill, 0)
	for _, v := range sks {
		modelSkills = append(modelSkills, v.ToModel().(model.Skill))
	}
	return modelSkills, nil
}
