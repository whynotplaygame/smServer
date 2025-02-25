package data

import (
	"log"
	"smServer/db"
	"smServer/server/game/model"
	"sync"
	"time"
)

var RoleCityDao = &mapRoleCityDao{
	rcChan: make(chan *MapRoleCity, 100),
}

type mapRoleCityDao struct {
	rcChan chan *MapRoleCity
}

func (m *mapRoleCityDao) running() {
	for {
		select {
		case rc := <-m.rcChan:
			if rc.CityId > 0 {
				//where  city_id = ?
				_, err := db.Engin.Table(rc).ID(rc.CityId).Update(rc)
				if err != nil {
					log.Println("mapRoleCityDao running error", err)
				}
			}
		}
	}
}

type MapRoleCity struct {
	mutex      sync.Mutex `xorm:"-"`
	CityId     int        `xorm:"cityId pk autoincr"`
	RId        int        `xorm:"rid"`
	Name       string     `xorm:"name" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	X          int        `xorm:"x"`
	Y          int        `xorm:"y"`
	IsMain     int8       `xorm:"is_main"`
	CurDurable int        `xorm:"cur_durable"`
	CreatedAt  time.Time  `xorm:"created_at"`
	OccupyTime time.Time  `xorm:"occupy_time"`
}

func (m *MapRoleCity) TableName() string {
	return "map_role_city"
}

func (m *MapRoleCity) ToModel() interface{} {
	p := model.MapRoleCity{}
	p.X = m.X
	p.Y = m.Y
	p.CityId = m.CityId
	p.UnionId = GetUnion(m.RId)
	p.UnionName = ""
	p.ParentId = 0
	p.MaxDurable = 1000
	p.CurDurable = m.CurDurable
	p.Level = 1
	p.RId = m.RId
	p.Name = m.Name
	p.IsMain = m.IsMain == 1
	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6
	return p
}

func (m *MapRoleCity) SyncExecute() {
	RoleCityDao.rcChan <- m
}
