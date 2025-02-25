package model

// 联盟列表
type ListReq struct {
}

type ListRsp struct {
	List []Union `json:"list"`
}

const (
	UnionChairman     = 0 //盟主
	UnionViceChairman = 1 //副盟主
	UnionCommon       = 2 //普通成员
)

const (
	UnionUntreated = 0 //未处理
	UnionRefuse    = 1 //拒绝
	UnionAdopt     = 2 //通过
)

type Union struct {
	Id     int     `json:"id"`     //联盟id
	Name   string  `json:"name"`   //联盟名字
	Cnt    int     `json:"cnt"`    //联盟人数
	Notice string  `json:"notice"` //公告
	Major  []Major `json:"major"`  //联盟主要人物，盟主副盟主
}
type Major struct {
	RId   int    `json:"rid"`
	Name  string `json:"name"`
	Title int8   `json:"title"`
}
type Member struct {
	RId   int    `json:"rid"`
	Name  string `json:"name"`
	Title int8   `json:"title"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

// 联盟信息
type InfoReq struct {
	Id int `json:"id"` //联盟id
}

type InfoRsp struct {
	Id   int   `json:"id"` //联盟id
	Info Union `json:"info"`
}

// 获取申请列表
type ApplyReq struct {
	Id int `json:"id"`
}

type ApplyRsp struct {
	Id     int         `json:"id"`
	Applys []ApplyItem `json:"applys"`
}

type ApplyItem struct {
	Id       int    `json:"id"`
	RId      int    `json:"rid"`
	NickName string `json:"nick_name"`
}

// 创建联盟
type CreateReq struct {
	Name string `json:"name"`
}

type CreateRsp struct {
	Id   int    `json:"id"` //联盟id
	Name string `json:"name"`
}

// 申请加入联盟
type JoinReq struct {
	Id int `json:"id"` //联盟id
}

type JoinRsp struct {
}

// 审核
type VerifyReq struct {
	Id     int  `json:"id"`     //申请操作的id
	Decide int8 `json:"decide"` //1是拒绝，2是通过
}

type VerifyRsp struct {
	Id     int  `json:"id"`     //申请操作的id
	Decide int8 `json:"decide"` //1是拒绝，2是通过
}

// 联盟成员
type MemberReq struct {
	Id int `json:"id"` //联盟id
}

type MemberRsp struct {
	Id      int      `json:"id"` //联盟id
	Members []Member `json:"Members"`
}
