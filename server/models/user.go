package models

import "time"

type User struct {
	UId      int       `xorm:"uid pk autoincr"` // 数据库的映射关系
	Username string    `xorm:"username" valitate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$]"`
	Passwd   string    `xorm:"Passwd" valitate:"min=1,max=30"`
	PassCode string    `xorm:"passcode"`
	HardWare string    `xorm:"hardware"`
	Status   int       `xorm:"status"`
	Ctime    time.Time `xorm:"ctime"`
	Mtime    time.Time `xorm:"mtime"`
	IsOnline bool      `xorm:"-"`
}

// xorm 自行指定表名

func (u *User) TableName() string {
	return "user"
}
