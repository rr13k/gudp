## gudp

## 简介

`gudp` 为`godot game` 所实现，提供轻量快速的可靠`udp`实现, 用户通过简单的函数即可快速搭建可靠的`udp`服务。

[godot-client](https://github.com/rr13k/gudp-godot)点击访问

使用`mit`协议开源，支持通过`fork`进行贡献。

## 特性
1. 支持可靠消息、不可靠消息
2. 支持身份认证
3. 内置支持重发和心跳检测
4. 使用`proto`作为通用协议
5. 支持`rpc`协议(`udp`实现)

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

## 使用rpc


``` go 
GudpServer, err = gudp.NewServerManager(udpConn)
if err != nil {
	fmt.Println(err.Error())
}
// NewServerManager连接后,配置rpc结构体
GudpServer.RegisterRpc(new(rpcs.UdpRpc))


// rpcs.go 文件中
type UdpRpc struct {
}

// 算数运算请求结构体
type UdpRpcRequest struct {
	A int
	B int
}

// 算数运算响应结构体
type UdpRpcResponse struct {
	Pro int // 乘积
	Quo int // 商
	Rem int // 余数
}

// 基础测试
func (u *UdpRpc) Multiply(req UdpRpcRequest, res *UdpRpcResponse) error {
	res.Pro = req.A * req.B
	fmt.Println("支持了rpc Multiply:", res.Pro)
	return nil
}

```

在`gdscript`中使用如下步骤调用并获取返回值
```js
var res = await gudpClient.RpcCall("Multiply",{"A":2,"B":3})
print(res)
```