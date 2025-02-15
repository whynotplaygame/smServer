package middleware

import (
	"fmt"
	"log"
	"smServer/net"
)

func Log() net.MiddlewareFunc {
	return func(next net.HandlerFunc) net.HandlerFunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			log.Println("请求的路由:", req.Body.Name)
			log.Println("请求的参数:", fmt.Sprintf("%v", req.Body.Msg))
			log.Printf("请求的参数:%v", req.Body.Msg)
			next(req, rsp)
			log.Println("响应的数据:", rsp.Body.Msg)
		}
	}
}
