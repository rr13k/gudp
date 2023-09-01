package main2

import (
	"fmt"
	"hash/crc32"

	"github.com/rr13k/gudp"
	"github.com/rr13k/gudp/protos"
	"go.starlark.net/lib/proto"
)

// m *testing.M

func Test_main() {
	host := "127.0.0.1"
	port := 12345

	udpConn := gudp.CreateUdpConn(host, port)
	udpServer, err := gudp.NewServerManager(udpConn)
	if err != nil {
		fmt.Println(err.Error())
	}

	udpServer.SetHandler(received)

	go func() { // handling the server errors
		for {
			uerr := <-udpServer.Error
			if uerr != nil {
				fmt.Println("Errors on udp server: ", uerr.Error())
			}
		}
	}()

	print("hero test main start success!")
	udpServer.Serve()
}

// 收到消息后处理,业务消息可靠性无关, 这里只能接受到用户发送的数据
func received(id string, buffer []byte) {
	msg := buffer
	var mych protos.ApiMessage
	err := proto.Unmarshal(msg, &mych)

	print("接收到了消息")

	// 这里只处理自身的pb,不理会udp的pb
	switch mych.NoticeWay.(type) {
	case *protos.ApiMessage_ReliableMessage:
		fmt.Println("可靠的udp")
		_reliableMessage := mych.GetReliableMessage()
		// 验证校验和
		receivedChecksum := crc32.ChecksumIEEE(_reliableMessage.Data)
		if _reliableMessage.Checksum != receivedChecksum {
			fmt.Println("Checksum mismatch, message may be corrupted")
			return
		}
		fmt.Println(receivedChecksum)

	case *protos.ApiMessage_UnreliableMessage:
		fmt.Println("不可靠")
	}

	if err != nil {
		fmt.Println("非pb")
	}

	// 回复客户端
	reply := []byte("Hello from server!")
	fmt.Println(reply)
	if err != nil {
		fmt.Println("Error sending reply:", err)
	}
}
