package logic

import (
	"encoding/json"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/common"
	"smServer/server/game/model"
	"smServer/server/game/model/data"
	"sync"
	"xorm.io/xorm"
)

var RoleAttrService = &roleAttrService{
	attrs: make(map[int]*data.RoleAttribute), // 初始化
}

type roleAttrService struct {
	mutex sync.RWMutex
	attrs map[int]*data.RoleAttribute
}

func (r *roleAttrService) Load() {
	ras := make([]*data.RoleAttribute, 0)
	err := db.Engin.Table(new(data.RoleAttribute)).Find(&ras)
	if err != nil {
		log.Println("load role attributes err:", err)
	}
	//
	for _, v := range ras {
		r.attrs[v.RId] = v
	}
	//查询所有的联盟，进行匹配
	uns := CoalitionService.ListCoalition()
	for _, un := range uns {
		for _, ma := range un.MemberArray {
			ra, ok := r.attrs[ma]
			if ok {
				ra.UnionId = un.Id
				r.attrs[ma] = ra
			}
		}
	}
	log.Println("load attrs:", r.attrs[2])
	log.Println("load role attributes success")
}

func (service *roleAttrService) TryCreate(rid int, req *net.WsMsgReq) error {
	// 根据rid查询，如果没有就创建
	role := &data.RoleAttribute{}
	ok, err := db.Engin.Table(role).Where("rid = ?", rid).Get(role)
	if err != nil {
		log.Println("查询角色属性出错", err)
		return common.New(constant.DBError, "数据库出错")
	}

	if ok {
		//service.mutex.Lock()
		//defer service.mutex.Unlock()
		//service.attrs[rid] = role // 由于getProperty同样也是查询，就把这个数据缓存下来。方便
		return nil
	} else {
		// 初始化
		role.RId = rid
		role.UnionId = 0
		role.ParentId = 0
		role.PosTags = ""
		if session := req.Context.Get("dbSession"); session != nil {
			_, err = session.(*xorm.Session).Table(role).Insert(role)
		} else {
			_, err = db.Engin.Table(role).Insert(role)
		}

		if err != nil {
			log.Println("插入角色属性出错", err)
			return common.New(constant.DBError, "数据库出错")
		}

		service.mutex.Lock()
		defer service.mutex.Unlock()
		service.attrs[rid] = role // 由于getProperty同样也是查询，就把这个数据缓存下来。减少一次数据库查询

	}
	return nil
}

func (service *roleAttrService) GetTagList(rid int) ([]model.PosTag, error) {
	ra, ok := service.attrs[rid] // 查看缓存的数据
	if !ok {
		ra := &data.RoleAttribute{}
		var err error
		ok, err = db.Engin.Table(ra).Where("rid = ?", rid).Get(ra)
		if err != nil {
			log.Println("查询角色属性出错", err)
			return nil, common.New(constant.DBError, "数据库错误")
		}
	}

	posTags := make([]model.PosTag, 0)

	if ok {
		tags := ra.PosTags
		if tags != "" {
			err := json.Unmarshal([]byte(tags), &posTags)
			if err != nil {
				return nil, common.New(constant.DBError, "数据库错误")
			}

		}
	}
	return posTags, nil
}

func (r *roleAttrService) Get(rid int) *data.RoleAttribute {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ra, OK := r.attrs[rid]
	if OK {
		return ra
	}
	return nil

}

func (r *roleAttrService) GetUnion(rid int) int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ram, ok := r.attrs[rid]
	log.Println("rid unionID", rid, ram.UnionId)
	if ok {
		return ram.UnionId
	}
	return 0
}

func (r *roleAttrService) IsHasUnion(rid int) bool {
	r.mutex.RLock()
	ra, ok := r.attrs[rid]
	r.mutex.RUnlock()
	if ok {
		return ra.UnionId != 0
	} else {
		return false
	}
}

func (ra *roleAttrService) TryCreateRA(rid int) (*data.RoleAttribute, bool) {
	attr := ra.Get(rid)
	if attr != nil {
		return attr, true
	} else {
		ra.mutex.Lock()
		defer ra.mutex.Unlock()
		attr := ra.create(rid)
		return attr, attr != nil
	}
}

func (ra *roleAttrService) create(rid int) *data.RoleAttribute {
	roleAttr := &data.RoleAttribute{RId: rid, ParentId: 0, UnionId: 0}
	if _, err := db.Engin.Insert(roleAttr); err != nil {
		return nil
	} else {
		ra.attrs[rid] = roleAttr
		return roleAttr
	}
}
