package data

import (
	"log"
	"smServer/db"
	"smServer/server/game/model"
	"time"
)

var RoleAttrDao = &roleAttrDao{
	raChan: make(chan *RoleAttribute, 100),
}

type roleAttrDao struct {
	raChan chan *RoleAttribute
}

func (r *roleAttrDao) running() {
	for {
		select {
		case rr := <-r.raChan:
			_, err := db.Engin.Table(new(RoleAttribute)).
				ID(rr.Id).
				Cols("parent_id", "collect_times", "last_collect_time", "pos_tags").
				Update(rr)
			if err != nil {
				log.Println("update role attribute err:", err)
			}
		}
	}
}

func init() {
	go RoleAttrDao.running() // 启动监听
}

type RoleAttribute struct {
	Id              int            `xorm:"id pk autoincr"`
	RId             int            `xorm:"rid"`
	UnionId         int            `xorm:"-"`                 //联盟id
	ParentId        int            `xorm:"parent_id"`         //上级id（被沦陷）
	CollectTimes    int8           `xorm:"collect_times"`     //征收次数
	LastCollectTime time.Time      `xorm:"last_collect_time"` //最后征收的时间
	PosTags         string         `xorm:"pos_tags"`          //位置标记
	PosTagArray     []model.PosTag `xorm:"-"`
}

func (r *RoleAttribute) TableName() string {
	return "role_attribute"
}

func (r *RoleAttribute) SyncExecute() {
	RoleAttrDao.raChan <- r
}
