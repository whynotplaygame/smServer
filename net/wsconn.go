package net

// 请求的格式
type ReqBody struct {
	Seq   int64       `json:"seq"`  // 序号
	Name  string      `json:"name"` // 用于标识路由
	Msg   interface{} `json:"msg"`
	Proxy string      `json:"proxy"`
}

// 响应的格式
type RspBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"` // 错误码
	Msg  interface{} `json:"msg"`
}

// 封装成request 和 respone
type WsMsgReq struct {
	Body *ReqBody `json:"body"`
	Conn WSConn
}

type WsMsgRsp struct {
	Body *RspBody `json:"body"`
}

// 理解为request请求，请求有参数，用于取参数
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

type Handshake struct {
	Key string `json:"key"`
}

type HeartBeat struct {
	CTime int64 `json:"ctime"`
	STime int64 `json:"stime"`
}
