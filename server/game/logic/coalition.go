package logic

import (
	"log"
	"smServer/db"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"sync"
)

var CoalitionService = &coalitionService{
	unions: make(map[int]*data.Coalition),
}

type coalitionService struct {
	mutex  sync.RWMutex
	unions map[int]*data.Coalition
}

func (c *coalitionService) Load() {
	rr := make([]*data.Coalition, 0)
	err := db.Engin.Table(new(data.Coalition)).Where("state=?", data.UnionRunning).Find(&rr)
	if err != nil {
		log.Println("coalitionService load error", err)
	}
	for _, v := range rr {
		c.unions[v.Id] = v
	}
}

func (c *coalitionService) List() ([]model.Union, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]model.Union, 0)
	for _, v := range c.unions {
		// 盟主
		mas := make([]model.Major, 0)
		if role := RoleService.Get(v.Chairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}

		// 副盟主
		if role := RoleService.Get(v.ViceChairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionViceChairman,
			}
			mas = append(mas, ma)
		}
		union := v.ToModel().(model.Union)
		union.Major = mas
		uns = append(uns, union)
	}
	return uns, nil
}

func (c *coalitionService) ListCoalition() []*data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]*data.Coalition, 0)
	for _, v := range c.unions {
		uns = append(uns, v)
	}
	return uns
}
