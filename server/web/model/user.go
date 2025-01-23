package model

import "time"

type User struct {
	UId      int       `xorm:"uid pk autoincr"` // 数据库的映射关系
	Username string    `xorm:"username" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$]"`
	Passwd   string    `xorm:"passwd" validate:"min=1,max=30"`
	Passcode string    `xorm:"passcode"`
	Hardware string    `xorm:"hardware"`
	Status   int       `xorm:"status"`
	Ctime    time.Time `xorm:"ctime"`
	Mtime    time.Time `xorm:"mtime"`
	IsOnline bool      `xorm:"-"`
}

// xorm 自行指定表名

func (u *User) TableName() string {
	return "user"
}
