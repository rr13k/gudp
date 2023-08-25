package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"hash/crc32"
	protoc "pulseHero/src/protos"
	"pulseHero/src/udp"
	"pulseHero/src/udp/crypto"

	"google.golang.org/protobuf/proto"
)

func main() {

	host := "127.0.0.1"
	port := 12345

	udpConn := udp.StartUdpServer(host, port)

	pk, err := crypto.GenerateRSAKey(2048)
	if err != nil {
		fmt.Println(err)
	}

	var customAuthenticate = &CustomAuthenticate{}
	cfg := newServerConfig(customAuthenticate, pk)
	udpServer, err := udp.NewServer(udpConn, cfg...)
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

	udpServer.Serve()
}

// 自定义身份验证
type CustomAuthenticate struct{}

// 身份验证是传递token返回用户id
func (m *CustomAuthenticate) Authenticate(ctx context.Context, token []byte) (string, error) {
	return "", nil
}

// 收到消息后处理,业务消息可靠性无关, 这里只能接受到用户发送的数据
func received(id string, buffer []byte) {
	msg := buffer
	var mych protoc.ApiMessage
	err := proto.Unmarshal(msg, &mych)

	print("接收到了消息")

	// 这里只处理自身的pb,不理会udp的pb
	switch mych.NoticeWay.(type) {
	case *protoc.ApiMessage_ReliableMessage:
		fmt.Println("可靠的udp")
		_reliableMessage := mych.GetReliableMessage()
		// 验证校验和
		receivedChecksum := crc32.ChecksumIEEE(_reliableMessage.Data)
		if _reliableMessage.Checksum != receivedChecksum {
			fmt.Println("Checksum mismatch, message may be corrupted")
			return
		}

	case *protoc.ApiMessage_UnreliableMessage:
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

func newServerConfig(a udp.AuthClient, pk *rsa.PrivateKey) []udp.Option {
	return []udp.Option{
		udp.WithAuthClient(a),
		udp.WithSymmetricCrypto(crypto.NewAES(crypto.AES_CBC)),
		udp.WithAsymmetricCrypto(crypto.NewRSAFromPK(pk)),
		udp.WithReadBufferSize(2048),
		udp.WithMinimumPayloadSize(4),
		udp.WithProtocolVersion(1, 0),
	}
}
