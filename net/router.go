package net

import (
	"log"
	"strings"
	"sync"
)

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc // 原有执行方法进行包装

// account   login || logou
type group struct {
	mutex         sync.RWMutex
	prefix        string // 前缀
	handlerMap    map[string]HandlerFunc
	middlewareMap map[string][]MiddlewareFunc // 中间件,一个路由有可能多个中间件
	middlewares   []MiddlewareFunc            // 所有中间件列表
}

// 添加处理函数的,即，websocket的处理函数
func (g *group) AddRouter(prefix string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	g.handlerMap[prefix] = handlerFunc
	g.middlewareMap[prefix] = middleware
}

func (g *group) Use(middleware ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:        prefix,
		handlerMap:    make(map[string]HandlerFunc),
		middlewareMap: make(map[string][]MiddlewareFunc),
	}
	r.group = append(r.group, g)
	return g
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h, ok := g.handlerMap[name] // 获取路由函数
	if !ok {                    // 没有匹配到
		h, ok = g.handlerMap["*"]
		if !ok {
			log.Println("路由未定义")
		}
	}
	if ok {
		//中间件， 执行路由之前，需要执行中间件
		for i := 0; i < len(g.middlewares); i++ {
			h = g.middlewares[i](h)
		}
		mm, ok := g.middlewareMap[name]
		if ok {
			for i := 0; i < len(mm); i++ {
				h = mm[i](h)
			}
		}

		h(req, rsp)
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
