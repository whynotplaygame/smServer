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
	list := gameConfig.FacilityConfig.List
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

func (c *cityFacilityService) GetByRId(rid int) ([]*data.CityFacility, error) {
	cf := make([]*data.CityFacility, 0)
	err := db.Engin.Table(new(data.CityFacility)).Where("rid=?", rid).Find(&cf)
	if err != nil {
		log.Println(err)
		return cf, common.New(constant.DBError, "数据库错误")
	}
	return cf, nil
}

func (c *cityFacilityService) GetYield(rid int) data.Yield {
	// 查询 靶标中的设置 获取到
	// 设施的不同乐行，去设施配置总查询匹配，匹配到增加产量的设施，木头，金钱， 产量的计算
	// 设施的等级不同，产量也不同

	cfs, err := c.GetByRId(rid)
	var y data.Yield
	if err == nil {
		for _, cf := range cfs {
			facilies := cf.Facility()
			for _, fa := range facilies {
				// 计算等级，资源产出是不同的
				if fa.GetLevel() > 0 {
					values := gameConfig.FacilityConfig.GetValues(fa.Type, fa.GetLevel())
					adds := gameConfig.FacilityConfig.GetAdditions(fa.Type)
					for i, t := range adds {
						if t == gameConfig.TypeWood {
							y.Wood += values[i]
						} else if t == gameConfig.TypeGrain {
							y.Grain += values[i]
						} else if t == gameConfig.TypeIron {
							y.Iron += values[i]
						} else if t == gameConfig.TypeStone {
							y.Stone += values[i]
						} else if t == gameConfig.TypeTax {
							y.Gold += values[i]
						}
					}

				}
			}
		}
	}

	return y
}
