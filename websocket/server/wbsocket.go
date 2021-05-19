package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	//Upgrader 用于升级 http 请求，把 http 请求升级为长连接的 WebSocket
	upgreder = websocket.Upgrader{
		// 指定升级 websocket 握手完成的超时时间
		HandshakeTimeout: 1024,
		// 指定 io 操作的缓存大小，如果不指定就会自动分配。
		ReadBufferSize: 1024,
		// 写数据操作的缓存池，如果没有设置值，write buffers 将会分配到链接生命周期里。
		//WriteBufferSize: 0,
		//WriteBufferPool: nil,
		//按顺序指定服务支持的协议，如值存在，则服务会从第一个开始匹配客户端的协议。
		//Subprotocols: nil,
		// 指定 http 的错误响应函数，如果没有设置 Error 则，会生成 http.Error 的错误响应。
		//Error: nil,
		// 请求检查函数，用于统一的链接检查，以防止跨站点请求伪造。如果不检查，就设置一个返回值为true的函数。
		// 如果请求Origin标头可以接受，CheckOrigin将返回true。 如果CheckOrigin为nil，则使用安全默认值：如果Origin请求头存在且原始主机不等于请求主机头，则返回false
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// EnableCompression 指定服务器是否应尝试协商每个邮件压缩（RFC 7692）。
		// 将此值设置为true并不能保证将支持压缩。
		// 目前仅支持“无上下文接管”模式
		//EnableCompression: false,
	}
)

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgreder.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}
