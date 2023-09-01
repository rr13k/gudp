## gudp

## 简介

gudp 为godot game 所实现，提供轻量快速的可靠udp实现, 用户通过简单的函数即可快速搭建可靠的udp服务。

使用mit协议开源，支持通过fork进行贡献。

## 特性
1. 支持可靠消息、不可靠消息
2. 支持身份认证
3. 内置支持重发和心跳检测
4. 使用proto作为通用协议

## 使用方式

```go

package main

import (
	"errors"
	"fmt"
	"github.com/rr13k/gudp"
	"google.golang.org/protobuf/proto"
)

var GudpServer *gudp.Server

func main() {
	host := "127.0.0.1"
	port := 12345

	var err error

	udpConn := gudp.CreateUdpConn(host, port)
	GudpServer, err = gudp.NewServerManager(udpConn)
	if err != nil {
		fmt.Println(err.Error())
	}
    
	GudpServer.SetHandler(onReceived)
	go func() { // handling the server errors
		for {
			uerr := <-GudpServer.Error
			if uerr != nil {
				fmt.Println("Errors on udp server: ", uerr.Error())
			}
		}
	}()

	GudpServer.Serve()
}

// 当接受到消息触发
func onReceived(client *gudp.Client, buffer []byte) {
	fmt.Println("on msg:", buffer)
	var hi = []byte("hello world~")
	// 发送可靠消息
	GudpServer.SendClientMessage(client, hi, true)
	// 发送不可靠消息
	GudpServer.SendClientMessage(client, hi, false)
}

```