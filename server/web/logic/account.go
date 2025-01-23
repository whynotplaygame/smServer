package logic

// 很神奇。提出来作为公共的user 模型，一重新打开工程，就引用不到。只能先在不同地方分别创建
// 同样的user model. 目前是web 和logic 都会有这个

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/utils"
	"time"

	"smServer/server/web/model"
)

var DefaultAccount = &AccountLogic{}

type AccountLogic struct {
}

func (l *AccountLogic) Register(rq *model.RegisterReq) error {
	username := rq.Username
	user := &model.User{}

	ok, err := db.Engin.Table(user).Where("username=?", username).Get(user)
	if err != nil {
		return common.New(constant.DBError, "数据库异常")
	}
	if ok {
		// 有数据，提示用户已存在
		return common.New(constant.UserExist, "用户已存在")
	} else {
		user.Mtime = time.Now()
		user.Ctime = time.Now()
		user.Username = rq.Username
		user.Passcode = utils.RandSeq(6)
		user.Passwd = utils.Password(rq.Password, user.Passcode)
		user.Hardware = rq.Hardware
		_, err = db.Engin.Table(user).Insert(user)
		if err != nil {
			log.Println("插入数据失败", err)
			return common.New(constant.DBError, "数据库异常")
		}
		return nil
	}
}
