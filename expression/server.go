package expression

import (
	"context"
	"crypto/rsa"
	"fmt"
	"hash/crc32"

	gudp "github.com/rr13k/gudp"

	"github.com/rr13k/gudp/crypto"
	"github.com/rr13k/gudp/gudp_protos"

	"google.golang.org/protobuf/proto"
)

func main2() {

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

	udpServer.Serve()
}

// 自定义身份验证
type CustomAuthenticate struct{}

// 身份验证是传递token返回用户id
func (m *CustomAuthenticate) Authenticate(ctx context.Context, token []byte) (string, error) {
	return "", nil
}

// 收到消息后处理,业务消息可靠性无关, 这里只能接受到用户发送的数据
func received(cli *gudp.Client, buffer []byte) {
	msg := buffer
	var mych gudp_protos.ApiMessage
	err := proto.Unmarshal(msg, &mych)

	print("接收到了消息")

	// 这里只处理自身的pb,不理会udp的pb
	switch mych.NoticeWay.(type) {
	case *gudp_protos.ApiMessage_ReliableMessage:
		fmt.Println("可靠的udp")
		_reliableMessage := mych.GetReliableMessage()
		// 验证校验和
		receivedChecksum := crc32.ChecksumIEEE(_reliableMessage.Data)
		if _reliableMessage.Checksum != receivedChecksum {
			fmt.Println("Checksum mismatch, message may be corrupted")
			return
		}
		fmt.Println(receivedChecksum)

	case *gudp_protos.ApiMessage_UnreliableMessage:
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

func newServerConfig(a gudp.AuthClient, pk *rsa.PrivateKey) []gudp.Option {
	return []gudp.Option{
		gudp.WithAuthClient(a),
		gudp.WithSymmetricCrypto(crypto.NewAES(crypto.AES_CBC)),
		gudp.WithAsymmetricCrypto(crypto.NewRSAFromPK(pk)),
		gudp.WithReadBufferSize(2048),
		gudp.WithMinimumPayloadSize(4),
	}
}
