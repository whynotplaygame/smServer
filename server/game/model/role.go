package model

type EnterServerReq struct {
	Session string `json:"session"`
}

type EnterServerRsp struct {
	Role    Role    `json:"role"`
	RoleRes RoleRes `json:"role_res"`
	Time    int64   `json:"time"`
	Token   string  `json:"token"`
}

// 数据库的字段，不一定是客户端需要的字段，做业务逻辑的时候，会将数据库的结果，映射到客户端需要的结果上
// 其中可能会做一些转换
// 所以 数据库操作 与 业务操作，分割，所以，数据库的role等需要一个toModel转换功能
// 业务对象
type Role struct {
	Rid      int    `json:"rid"`
	UId      int    `json:"uid"`
	NickName string `json:"nickName"`
	Sex      int8   `json:"sex"`
	Balance  int    `json:"balance"`
	HeadId   int    `json:"headId"`
	Profile  string `json:"profile"`
}

type RoleRes struct {
	Wood          int `json:"wood"`
	Iron          int `json:"iron"`
	Stone         int `json:"stone"`
	Grain         int `json:"grain"`
	Gold          int `json:"gold"`
	Decree        int `json:"decree"`
	WoodYield     int `json:"wood_yield"`
	IronYield     int `json:"iron_yield"`
	StoneYield    int `json:"stone_yield"`
	GrainYield    int `json:"grain_yield"`
	GoldYield     int `json:"gold_yield"`
	DepotCapacity int `json:"depot_capacity"` // 仓库容量
}
