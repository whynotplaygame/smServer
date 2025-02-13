package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
)

var WarService = &warService{}

type warService struct {
}

func (w *warService) GetWarReport(rid int) ([]model.WarReport, error) {
	wrs := make([]data.WarReport, 0)
	wr := &data.WarReport{}
	err := db.Engin.Table(wr).Where("a_rid = ? or d_rid", rid).
		Limit(30, 0).
		Desc("ctime").
		Find(&wrs)

	if err != nil {
		log.Println("查询战报军队出错", err)
		return nil, common.New(constant.DBError, "查询玩家战报出错")
	}
	modelWrs := make([]model.WarReport, 0)
	for _, v := range wrs {
		modelWrs = append(modelWrs, v.ToModel().(model.WarReport))
	}
	return modelWrs, nil
}
