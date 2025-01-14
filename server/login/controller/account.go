package controller

import (
	"smServer/net"
	"smServer/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{}
	loginRes.UId = 1
	loginRes.Username = "admin"
	loginRes.Session = "123"
	loginRes.Password = ""
	rsp.Body.Msg = loginRes

}
