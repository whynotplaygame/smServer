package logic

import (
	"encoding/json"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/gameConfig"
	"smServer/server/game/model/data"
	"xorm.io/xorm"
)

var CityFacilityService = &cityFacilityService{}

type cityFacilityService struct {
}

func (c *cityFacilityService) TryCreate(cid, rid int, req *net.WsMsgReq) error {
	cf := &data.CityFacility{}
	ok, err := db.Engin.Table(cf).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println("查询城市设施出错", err)
		return common.New(constant.DBError, "数据库错误")
	}
	if ok {
		return nil
	}
	cf.RId = rid
	cf.CityId = cid
	list := gameConfig.FaiclityConfig.List
	facs := make([]data.Facility, len(list))
	for index, v := range list {
		fac := data.Facility{
			Name:         v.Name,
			Type:         v.Type,
			PrivateLevel: 0,
			UpTime:       0,
		}
		facs[index] = fac
	}
	dataJson, _ := json.Marshal(facs)
	cf.Facilities = string(dataJson)

	if session := req.Context.Get("dbSession"); session != nil {
		_, err = session.(*xorm.Session).Table(cf).Insert(cf)
	} else {
		_, err = db.Engin.Table(cf).Insert(cf)
	}

	if err != nil {
		log.Println("插入城市设施出错", err)
		return common.New(constant.DBError, "数据库错误")
	}
	return nil
}
