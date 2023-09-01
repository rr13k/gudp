package gudp

import (
	"fmt"
	"testing"
)

const MaxSeqNum = 32

type AckBitmap struct {
	bitmap uint32
}

func (ab *AckBitmap) Set(packetNum int) {
	ab.bitmap |= (1 << packetNum)
}

func parseAckBitmap(bitmap uint32) []int {
	var receivedPackets []int

	for i := 0; i < MaxSeqNum; i++ {
		if (bitmap & (1 << i)) != 0 {
			receivedPackets = append(receivedPackets, i)
		}
	}

	return receivedPackets
}

func Test_ackBits(t *testing.T) {

	ackBitmap := AckBitmap{}

	// Simulate receiving packets and updating the ACK bitmap
	receivedPackets := []int{31}

	// for i := 0; i < 30; i++ {
	// 	receivedPackets = append(receivedPackets, i)
	// }

	for _, packetNum := range receivedPackets {
		ackBitmap.Set(packetNum)
		fmt.Printf("Received packet %d, updated ACK bitmap: %032b\n", packetNum, ackBitmap.bitmap)
	}

	bitmap := ackBitmap.bitmap

	// 模拟客户端解析
	fmt.Println("客户端解析:", parseAckBitmap(bitmap))

	fmt.Println("测试完成", bitmap)

}
