package net

import (
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

type ProxyClient struct {
	proxy string
	conn  *ClientConn
}

func (c *ProxyClient) Connect() error {
	// 去连接 websocket 服务端
	//通过Dialer连接websocket服务器
	var dialer = websocket.Dialer{
		Subprotocols:     []string{"p1", "p2"},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	ws, _, err := dialer.Dial(c.proxy, nil) //  发起连接
	if err == nil {
		c.conn = NewClientConn(ws)
		if !c.conn.Start() {
			return errors.New("握手失败")
		}
	}
	return err
}

func NewProxyClient(proxy string) *ProxyClient {
	return &ProxyClient{
		proxy: proxy,
	}
}
