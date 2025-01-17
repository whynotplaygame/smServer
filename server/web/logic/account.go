package logic

import (
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/server/common"
	"smServer/utils"
	"time"

	//loginmodel "smServer/server/models"
	"smServer/server/models"
	"smServer/server/web/model"
)

var DefaultAccount = &AccountLogic{}

type AccountLogic struct {
}

func (l *AccountLogic) Register(rq *model.RegisterReq) error {
	username := rq.Username
	user := &models.User{}

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
		user.PassCode = utils.RandSeq(6)
		user.Passwd = utils.Password(rq.Password, user.PassCode)
		user.HardWare = rq.Hardware
		_, err = db.Engin.Table(user).Insert(user)
		if err != nil {
			log.Println("插入数据失败", err)
			return common.New(constant.DBError, "数据库异常")
		}
		return nil
	}
}
