package net

import (
	"log"
	"strings"
)

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

// account   login || logou
type group struct {
	prefix     string // 前缀
	handlerMap map[string]HandlerFunc
}

// 添加处理函数的,即，websocket的处理函数
func (g *group) AddRouter(prefix string, handlerFunc HandlerFunc) {
	g.handlerMap[prefix] = handlerFunc
}

func (r *Router) Group(prefix string) *group {
	g := &group{prefix: prefix, handlerMap: make(map[string]HandlerFunc)}
	r.group = append(r.group, g)
	return g
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name] // 获取路由函数
	if h != nil {
		h(req, rsp)
	} else { // 没有匹配到
		h = g.handlerMap["*"]
		if h != nil {
			h(req, rsp)
		} else {
			log.Println("路由未定义")
		}
	}
}

type Router struct {
	group []*group
}

// 由于router是小写，为了可以暴露被用，使用NewRouter，方便外部调用
func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	// req.Body.Name  路径   登录业务  account.login,     account 是组，login 路由标识
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0] // 前缀
		name = strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		} else if g.prefix == "*" { // 为网关的*放行
			g.exec(name, req, rsp)
		}
	}
}
