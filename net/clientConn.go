package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"smServer/constant"
	"smServer/utils"
	"sync"
	"time"
)

type syncCtx struct {
	// groutine 的上下文，包含groutine的运行状态，环境，现场等信息
	ctx     context.Context
	cancel  context.CancelFunc
	outChan chan *RspBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	return &syncCtx{
		ctx:     ctx,
		cancel:  cancel,
		outChan: make(chan *RspBody),
	}
}

func (s *syncCtx) wait() *RspBody {
	select {
	case msg := <-s.outChan:
		return msg
	case <-s.ctx.Done(): // 超时
		log.Println("代理服务器相应消息超时")
		return nil
	}
}

type ClientConn struct {
	wsConn        *websocket.Conn
	handshake     bool      // 握手状态
	handshakeChan chan bool // 用于接收握手消息

	isClosed     bool // 监听连接是否关闭状态
	property     map[string]interface{}
	propertyLock sync.RWMutex
	Seq          int64

	onPush      func(conn *ClientConn, body *RspBody) // 通知代理服务器
	onClose     func(conn *ClientConn)
	syncCtxMap  map[int64]*syncCtx
	syncCtxLock sync.RWMutex
}

func (c *ClientConn) Start() bool {
	// 做的事情，就是 一直不停地接收消息
	// 等待握手的信息返回
	c.handshake = false      // 刚开始是false状态
	go c.wsReadloop()        // 一直读消息
	return c.waitHandHhake() // 等待握手的成功

}

func (c *ClientConn) waitHandHhake() bool {
	// 等待握手的成功，等待握手的消息
	// 万一程序出现问题，超时了一直响应不到
	if !c.isClosed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //维护goroutine中的信息，可以用作超时设置
		defer cancel()
		select {
		case _ = <-c.handshakeChan:
			log.Println("handshake done")
			return true
		case <-ctx.Done(): // 如果已经到了5秒钟
			log.Println("handshake timeout")
			return false
		}
	}
	return true
}

func (c *ClientConn) wsReadloop() {
	//for {
	//	_, data, err := c.wsConn.ReadMessage()
	//	fmt.Println(data, " ", err)
	//	// 收到握手消息了,放到管道
	//	// 读取到的信息，可能很多，握手，心跳，心跳，请求 accountt.login
	//	c.handshake = true
	//	c.handshakeChan <- true
	//
	//
	//}
	// 执行完，还是要走关闭通道逻辑
	defer func() {
		if err := recover(); err != nil {
			log.Println("捕捉到异常：", err)
			c.Close()
		}
	}()

	// 先读到客户端发过来的数据，然后处理，然后再回消息
	// 要经过路由 实际处理程序

	for {
		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		// fmt.Println("收到客户端发来的原始消息", data)
		// 收到消息，解析消息，前端发过来的消息，就是json格式
		// 1,data 解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错", err)
			continue
		}
		// 2,前端的消息，加密的消息，进行解密
		secretKey, err := c.GetProperty("secretKey")
		if err == nil {
			// 有加密
			key := secretKey.(string) // 将key转换成string
			// 客户端传过来的的数据是加密，需要解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误 in readloop。", err)
				// 出错后，不用发起握手
				// c.Handshake()
			} else {
				data = d
			}
		}
		// 3，data 转为body
		// fmt.Println("收到客户端发来的消息", data)
		body := &RspBody{}
		jsonErr := json.Unmarshal(data, body)
		if jsonErr != nil {
			log.Println("json格式错误 in readloop ", jsonErr)
		} else {
			// 判断握手 或者是别的请求,初次请求
			if body.Seq == 0 {
				if body.Name == HandshakeMsg {
					//获取秘钥
					hs := &Handshake{}
					mapstructure.Decode(body.Msg, hs)
					if hs.Key != "" {
						c.SetProperty("secretKey", hs.Key)
					} else {
						c.RemoveProperty("secretKey")
					}
					c.handshake = true
					c.handshakeChan <- true
				} else {
					//其他协议
					if c.onPush != nil {
						c.onPush(c, body)
					}
				}
			} else {
				// 非首次请求
				c.syncCtxLock.Lock()
				ctx, ok := c.syncCtxMap[body.Seq] // 通过序号
				c.syncCtxLock.Unlock()
				if ok {
					ctx.outChan <- body
				} else {
					log.Println("SEQ未发现：", body.Seq)
				}

			}
		}

	}
	c.Close() // 跳出循环，关闭通道

}

func (c *ClientConn) Close() {
	_ = c.wsConn.Close()
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn {
	return &ClientConn{
		wsConn:        wsConn,
		handshakeChan: make(chan bool), // 初始化管道
		Seq:           0,
		isClosed:      false,
		property:      make(map[string]interface{}),
		syncCtxMap:    map[int64]*syncCtx{},
	}
}

// 复制自wsserver ,改动
func (c *ClientConn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()         //先获取锁
	defer c.propertyLock.Unlock() //执行完解锁
	c.property[key] = value
}

// 获取属性
func (c *ClientConn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()         //先获取锁
	defer c.propertyLock.RUnlock() //执行完解锁
	value, ok := c.property[key]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("not found property")
	}
}

// 移除属性
func (c *ClientConn) RemoveProperty(key string) {
	c.propertyLock.Lock()         //先获取锁
	defer c.propertyLock.Unlock() //执行完解锁
	delete(c.property, key)
}

// 获取地址
func (c *ClientConn) Addr() string {
	return c.wsConn.RemoteAddr().String() // 通过连接中获取
}

// 往管道里放消息
func (c *ClientConn) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}} // 把消息内容封装起来
	// c.outChan <- rsp
	fmt.Println(rsp)
	c.write(rsp.Body)
}

func (c *ClientConn) write(body interface{}) error {
	fmt.Println("写消息", body)
	// 1,把数据转成json
	data, err := json.Marshal(body)
	if err != nil {
		log.Println("msg 2 json err:", err)
		return err
	}

	//// 2 加密
	//secretKey, err := c.GetProperty("secretKey")
	//if err == nil {
	//	key := secretKey.(string)
	//	// 数据加密
	//	data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	//	if err != nil {
	//		log.Println("加密失败", err)
	//		return err
	//	}
	//}
	// 3压缩 并写 回去
	if data, err := utils.Zip(data); err == nil {
		err := c.wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			println("写数据失败", err)
			return err
		}
	} else {
		log.Println("压缩失败", err)
		return err
	}
	return nil
}

func (c *ClientConn) setOnPush(hook func(conn *ClientConn, body *RspBody)) {
	c.onPush = hook
}

func (c *ClientConn) Send(name string, msg interface{}) *RspBody {
	// 把请求 发送给 代理服务器，登录服务器，等待返回
	c.Seq += 1
	seq := c.Seq
	sc := NewSyncCtx()
	c.syncCtxLock.Lock()
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()

	rsp := &RspBody{Name: name, Seq: seq, Code: constant.OK}

	//req 请求
	req := &ReqBody{Seq: seq, Name: name, Msg: msg}
	err := c.write(req)
	if err != nil {
		sc.cancel()
	} else {
		r := sc.wait()
		if r == nil {
			rsp.Code = constant.ProxyConnectError
		} else {
			rsp = r
		}
		return r
	}
	c.syncCtxLock.Lock()
	delete(c.syncCtxMap, seq)
	c.syncCtxLock.Unlock()
	return rsp
}
