package net

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type ClientConn struct {
	wsConn        *websocket.Conn
	handshake     bool      // 握手状态
	handshakeChan chan bool // 用于接收握手消息
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

func (c *ClientConn) wsReadloop() {
	for {
		_, data, err := c.wsConn.ReadMessage()
		fmt.Println(data, " ", err)
		// 收到握手消息了,放到管道
		c.handshake = true
		c.handshakeChan <- true
	}

}

func NewClientConn(wsConn *websocket.Conn) *ClientConn {
	return &ClientConn{
		wsConn:        wsConn,
		handshakeChan: make(chan bool), // 初始化管道
	}
}
