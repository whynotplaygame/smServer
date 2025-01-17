package model

import "time"

const (
	Login = iota
	Logout
)

type LoginHistory struct {
	Id       int       `xorm:"id pk autoincr"` // 数据库的映射关系
	UId      int       `xorm:"uid"`
	CTime    time.Time `xorm:"Ctime"`
	Ip       string    `xorm:"ip"`
	State    int8      `xorm:"state"`
	Hardware string    `xorm:"hardware"`
}

type LoginLast struct {
	Id         int       `xorm:"id pk autoincr"` // 数据库的映射关系
	UId        int       `xorm:"uid"`
	LoginTime  time.Time `xorm:"login_time"`
	logoutTime time.Time `xorm:"logout_time"`
	Session    string    `xorm:"session"`
	Ip         string    `xorm:"ip"`
	IsLogout   int8      `xorm:"is_logout"`
	Hardware   string    `xorm:"hardware"`
}

// xorm 自行指定表名

func (u *LoginHistory) TableName() string {
	return "login_history"
}

// xorm 自行指定表名

func (*LoginLast) TableName() string {
	return "login_last"
}
