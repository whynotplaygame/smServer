package data

import (
	"log"
	"smServer/db"
	"smServer/server/game/model"
)

var RoleResDao = &roleResDao{
	rrChan: make(chan *RoleRes, 100),
}

type roleResDao struct {
	rrChan chan *RoleRes
}

func (d *roleResDao) running() {
	for {
		select {
		case rr := <-d.rrChan:
			_, err := db.Engin.Table(new(RoleRes)).
				ID(rr.Id).
				Cols("wood", "iron", "stone", "grain", "gold").
				Update(rr)
			if err != nil {
				log.Println("update role res err:", err)
			}
		}
	}
}

func init() {
	go RoleResDao.running()
}

// 产量
type Yield struct {
	Wood  int
	Iron  int
	Stone int
	Grain int
	Gold  int
}

type RoleRes struct {
	Id     int `xorm:"id pk autoincr"`
	RId    int `xorm:"rid"`
	Wood   int `xorm:"wood"`
	Iron   int `xorm:"iron"`
	Stone  int `xorm:"stone"`
	Grain  int `xorm:"grain"`
	Gold   int `xorm:"gold"`
	Decree int `xorm:"decree"`
}

func (r *RoleRes) TableName() string {
	return "role_res"
}

func (r *RoleRes) ToModel() interface{} {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Iron = r.Iron
	p.Stone = r.Stone
	p.Grain = r.Grain
	p.Wood = r.Wood
	p.Decree = r.Decree

	yield := GetYield(r.RId)

	p.GoldYield = yield.Gold
	p.IronYield = yield.Iron
	p.StoneYield = yield.Stone
	p.GrainYield = yield.Grain
	p.WoodYield = yield.Wood
	p.DepotCapacity = 10000
	return p
}

func (r *RoleRes) SyncExecute() {
	RoleResDao.rrChan <- r
}
