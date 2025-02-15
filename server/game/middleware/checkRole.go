package middleware

import (
	"log"
	"smServer/constant"
	"smServer/net"
)

func CheckRole() net.MiddlewareFunc {
	return func(next net.HandlerFunc) net.HandlerFunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			log.Println("进入到角色检测中间件....")
			_, err := req.Conn.GetProperty("role")
			if err != nil {
				rsp.Body.Code = constant.RoleNotInConnect
				return
			}
			log.Println("中间件角色检测通过")
			next(req, rsp)
		}
	}
}
