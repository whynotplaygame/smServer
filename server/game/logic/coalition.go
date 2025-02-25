package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"sync"
	"time"
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

func (c *coalitionService) Get(id int) (model.Union, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		union := coa.ToModel().(model.Union)
		//盟主和副盟主信息
		main := make([]model.Major, 0)
		if role := RoleService.Get(coa.Chairman); role != nil {
			m := model.Major{Name: role.NickName, RId: role.RId, Title: model.UnionChairman}
			main = append(main, m)
		}
		if role := RoleService.Get(coa.ViceChairman); role != nil {
			m := model.Major{Name: role.NickName, RId: role.RId, Title: model.UnionViceChairman}
			main = append(main, m)
		}

		union.Major = main
		return union, nil
	}
	return model.Union{}, nil
}

func (c *coalitionService) GetCoalition(id int) *data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		return coa
	}
	return nil
}

func (c *coalitionService) GetListApply(unionId int, state int) ([]model.ApplyItem, error) {
	applys := make([]data.CoalitionApply, 0)
	err := db.Engin.Table(new(data.CoalitionApply)).
		Where("union_id = ? and state=?", unionId, state).
		Find(&applys)
	if err != nil {
		log.Println("coalitionService listApply error", err)
		return nil, common.New(constant.DBError, "数据库错误")
	}

	ais := make([]model.ApplyItem, 0)
	for _, v := range applys {
		var ai model.ApplyItem
		ai.RId = v.RId
		role := RoleService.Get(v.RId)
		ai.RId = role.RId
		ai.Id = v.Id
		ai.NickName = role.NickName
		ais = append(ais, ai)
	}
	return ais, nil

}

func (c *coalitionService) Create(name string, rid int) (*data.Coalition, bool) {
	m := &data.Coalition{Name: name,
		Ctime:       time.Now(),
		CreateId:    rid,
		Chairman:    rid,
		State:       data.UnionRunning,
		MemberArray: []int{rid}}

	_, err := db.Engin.Table(new(data.Coalition)).Insert(m)
	if err == nil {

		c.mutex.Lock()
		c.unions[m.Id] = m
		c.mutex.Unlock()

		return m, true
	} else {
		return nil, false
	}
}

func (c *coalitionService) MemberEnter(rid int, unionId int) {
	attr, ok := RoleAttrService.TryCreateRA(rid)
	if ok {
		attr.UnionId = unionId
		if attr.ParentId == unionId {
			c.DelChild(unionId, attr.RId)
		}
	}

	if rcs, ok := RoleCityService.GetByRId(rid); ok {
		for _, rc := range rcs {
			rc.SyncExecute()
		}
	}
}

func (c *coalitionService) DelChild(rid int, id2 int) {
	attr := RoleAttrService.Get(rid)
	if attr != nil {
		attr.ParentId = 0
		attr.SyncExecute()
	}
}

func (c *coalitionService) NewCreateLog(opNickName string, unionId int, opRId int) {
	ulog := &data.CoalitionLog{
		UnionId:  unionId,
		OPRId:    opRId,
		TargetId: 0,
		State:    data.UnionOpCreate,
		Des:      opNickName + " 创建了联盟",
		Ctime:    time.Now(),
	}

	db.Engin.InsertOne(ulog)
}

func (c *coalitionService) GetById(id int) *data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		return coa
	}
	return nil
}

func (c *coalitionService) NewJoin(targetNickName string, unionId int, opRId int, targetId int) {
	ulog := &data.CoalitionLog{
		UnionId:  unionId,
		OPRId:    opRId,
		TargetId: targetId,
		State:    data.UnionOpJoin,
		Des:      targetNickName + " 加入了联盟",
		Ctime:    time.Now(),
	}
	db.Engin.InsertOne(ulog)
}
