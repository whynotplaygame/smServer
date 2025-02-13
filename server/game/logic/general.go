package logic

import (
	"encoding/json"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/gameConfig/general"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"time"
)

var GeneralService = &gereralService{}

type gereralService struct {
}

func (service *gereralService) GetGenerals(rid int) ([]model.General, error) {
	grs := make([]*data.General, 0)
	gr := &data.General{}
	err := db.Engin.Table(gr).Where("rid = ?", rid).Find(&grs)

	if err != nil {
		log.Println("查询玩家武将出错", err)
		return nil, common.New(constant.DBError, "武将查询出错")
	}
	if len(grs) <= 0 { // 没有武将则随机3个
		var count = 0
		for {
			if count >= 3 {
				break
			}
			cfgId := general.General.Rand()
			gen, err := service.NewGeneral(cfgId, rid, 0)
			if err != nil {
				log.Println(err)
				continue
			}
			grs = append(grs, gen)
			count++
		}

	}
	modelGrs := make([]model.General, 0)
	for _, v := range grs {
		modelGrs = append(modelGrs, v.ToModel().(model.General))
	}
	return modelGrs, nil
}

const (
	GeneralNormal      = 0 //正常
	GeneralComposeStar = 1 //星级合成
	GeneralConvert     = 2 //转换
)

func (service *gereralService) NewGeneral(cfgId int, rid int, level int8) (*data.General, error) {
	cfg := general.General.GMap[cfgId]
	// 初始武将 无技能，但是有3个技能槽
	sa := make([]*model.GSkill, 3)
	ss, _ := json.Marshal(sa)
	gen := &data.General{
		PhysicalPower: gameConfig.Base.General.PhysicalPowerLimit,
		RId:           rid,
		CfgId:         cfg.CfgId,
		Order:         0,
		CityId:        0,
		Level:         level,
		CreatedAt:     time.Now(),
		CurArms:       cfg.Arms[0],
		HasPrPoint:    0,
		UsePrPoint:    0,
		AttackDis:     0,
		ForceAdded:    0,
		StrategyAdded: 0,
		DefenseAdded:  0,
		SpeedAdded:    0,
		DestroyAdded:  0,
		Star:          cfg.Star,
		StarLv:        0,
		ParentId:      0,
		SkillsArray:   sa,
		Skills:        string(ss),
		State:         GeneralNormal,
	}

	_, err := db.Engin.Table(gen).Insert(gen)
	if err != nil {
		log.Println("插入玩家武将出错", err)
		return nil, err
	}
	return gen, nil
}
