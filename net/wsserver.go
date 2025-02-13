package net

import (
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"smServer/utils"
	"sync"
	"time"
)

type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp // 写队列，需要回复给客户端的的信息队列
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex // 为在写属性的时候，用的读写锁
	needSecret   bool
}

var cid int64

func NewWsServer(wsConn *websocket.Conn, needSecret bool) *wsServer {
	s := &wsServer{
		wsConn:     wsConn,
		outChan:    make(chan *WsMsgRsp, 1000),
		property:   make(map[string]interface{}),
		Seq:        0,
		needSecret: needSecret,
	}
	cid++
	s.SetProperty("cid", cid)
	return s
}

// 目前 已经完成对了服务的设置
func (ws *wsServer) Router(router *Router) {
	ws.router = router
}

// 对接口函数进行实现
// 对属性进行修改
func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()         //先获取锁
	defer w.propertyLock.Unlock() //执行完解锁
	w.property[key] = value
}

// 获取属性
func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()         //先获取锁
	defer w.propertyLock.RUnlock() //执行完解锁
	value, ok := w.property[key]
	if ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("not found property")
	}
}

// 移除属性
func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()         //先获取锁
	defer w.propertyLock.Unlock() //执行完解锁
	delete(w.property, key)
}

// 获取地址
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String() // 通过连接中获取
}

// 往管道里放消息
func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}} // 把消息内容封装起来
	w.outChan <- rsp
}

// 通道一旦建立，那么 收发消息 就要一直监听
func (w *wsServer) Start() {
	// 启动读写数据的处理逻辑
	go w.readMsLoop()
	go w.writeMsLoop()
}

// 写数据
func (w *wsServer) write(msg interface{}) {
	fmt.Println("写消息", msg)
	// 1,把数据转成json
	data, err := json.Marshal(msg)
	log.Println("服务器写数据：", string(data))
	if err != nil {
		log.Println("msg 2 json err:", err)
	}

	// 2 加密
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		key := secretKey.(string)
		// 数据加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	// 3压缩 并写 回去
	if data, err := utils.Zip(data); err == nil {
		w.wsConn.WriteMessage(websocket.BinaryMessage, data)
	}
}

// 写循环
func (w *wsServer) writeMsLoop() {
	for {
		select {
		case msg := <-w.outChan:
			w.write(msg.Body)
		}
	}
}

func (w *wsServer) readMsLoop() {

	// 执行完，还是要走关闭通道逻辑
	defer func() {
		if err := recover(); err != nil {
			//log.Fatal(err)
			log.Println("read msg err in readloop in defer:", err)
			w.Close()
		}
	}()

	// 先读到客户端发过来的数据，然后处理，然后再回消息
	// 要经过路由 实际处理程序

	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		// fmt.Println("收到客户端发来的原始消息", data)
		// 收到消息，解析消息，前端发过来的消息，就是json格式
		// 1,data 解压 unzip
		data, err = utils.UnZip(data)
		log.Println("解压出的数据 in readloop：", data)
		if err != nil {
			log.Println("解压数据出错", err)
			continue
		}
		// 2,前端的消息，加密的消息，进行解密

		if w.needSecret { // 如果需要加密
			secretKey, err := w.GetProperty("secretKey")
			if err == nil {
				// 有加密
				key := secretKey.(string) // 将key转换成string
				log.Println("加密的key为", key, "time", time.Now().Unix())
				// 客户端传过来的的数据是加密，需要解密
				d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
				if err != nil {
					log.Println("数据格式有误 解密失败，key:", key, "in readloop。d:", d, " error:", err)
					// 出错后，发起握手
					w.Handshake()
				} else {
					data = d
				}
			}
		}

		// 3，data 转为body
		// fmt.Println("收到客户端发来的消息", data)
		body := &ReqBody{}
		jsonErr := json.Unmarshal(data, body)
		if jsonErr != nil {
			log.Println("json格式错误 in readloop ", jsonErr)
		} else {
			log.Println("收到的前端的数据：", string(data))

			context := &WsContext{
				property: make(map[string]interface{}),
			}
			// 获取到前端传递的数据了，拿上数据，去具体业务进行处理
			req := &WsMsgReq{Conn: w, Body: body, Context: context}

			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}

			if req.Body.Name == "heartbeat" {
				// 回心跳消息
				h := &HeartBeat{}
				mapstructure.Decode(req.Body.Msg, h)
				h.STime = time.Now().UnixNano() / 1e6 // 纳秒
				rsp.Body.Msg = h

			} else { // 非心跳，走路由
				if w.router != nil {
					log.Println("路由执行:", req)
					w.router.Run(req, rsp) // 对rsp进行赋值，赋完值，放到respone队列
				}
			}

			w.outChan <- rsp
		}

	}
	w.Close() // 跳出循环，关闭通道
}

// 既然是服务，还是要有关闭逻辑
func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}

const HandshakeMsg = "handshake"

// 当游戏客户端，发送请求的时候，会先记性握手协议
// 后端会发送对应的加密key给客户端
// 客户端在发送数据的时候，就会用此key进行加密处理

func (w *wsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string) // 如果有值
	} else {
		secretKey = utils.RandSeq(16) // 没值，生成一个
	}

	handshake := &Handshake{Key: secretKey}

	body := &RspBody{Name: HandshakeMsg, Msg: handshake}
	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data, err := utils.Zip(data); err == nil {
			w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
