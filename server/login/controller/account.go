package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"smServer/constant"
	"smServer/db"
	"smServer/net"
	"smServer/server/login/model"
	"smServer/server/login/proto"
	"smServer/utils"
	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	/*
		1, 用户名 密码，硬件id
		2, 根据用户名，查询user 表，得到数据
		3，进行密码对比，如果密码正确，登录成功
		4，保存用户登录记录
		5，爆粗男用户的最后一次登录信息
		6，客户端需要一个seesion, jwt 生成一个加密字符串算法
		7，客户端 在发起需要用户登录的行为时，判断用户是否合法
	*/
	loginReq := &proto.LoginReq{}
	loginRes := &proto.LoginRsp{}
	mapstructure.Decode(req.Body.Msg, loginReq) // 将map数据转成结构
	// username := loginReq.Username
	user := &model.User{}
	//                        表名
	ok, err := db.Engin.Table(user).Where("username=?", loginReq.Username).Get(user)
	if err != nil {
		log.Println("查询失败", err)
		return
	}
	if !ok { // 用户名存在
		// 没有查出来数据
		rsp.Body.Code = constant.UserNotExist
	}

	pwd := utils.Password(loginReq.Password, user.PassCode)

	if pwd != user.Passwd { // 密码不正确
		rsp.Body.Code = constant.PwdIncorrect
		return
	}

	// jwt a,b,c 三部分，a定义加密算法，b定义放入的数据 部分， 根据秘钥+a和b 生成加密字符串
	// 生成token,返回给前端
	token, _ := utils.Award(user.UId)

	rsp.Body.Code = constant.OK
	loginRes.UId = user.UId
	loginRes.Username = user.Username
	loginRes.Session = token
	loginRes.Password = ""
	rsp.Body.Msg = loginRes

	// 保存用户登录信息
	ul := &model.LoginHistory{
		UId:      user.UId,
		CTime:    time.Now(),
		Ip:       loginReq.Ip,
		Hardware: loginReq.Hardware,
		State:    model.Login,
	}

	_, err = db.Engin.Table(ul).Insert(ul)
	if err != nil {
		log.Println("插入用户登录信息报错：", err)
	}

	// 最后一次登录
	ll := &model.LoginLast{}
	ok, _ = db.Engin.Table(ll).Where("uid=?", user.UId).Get(ll)
	if ok {
		// 有数据，更新
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		db.Engin.Table(ll).Update(ll)
	} else {
		// 没数据，插入
		ll.IsLogout = 0
		ll.UId = user.UId
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		_, err = db.Engin.Table(ll).Insert(ll)
		if err != nil {
			log.Println("插入用户最后登录报错：", err)
		}
	}

	// 缓存一下，此用户和单签的ws连接,如果有重登陆，往客户端推送重登消
	net.Mgr.UserLogin(req.Conn, user.UId, token)
}
