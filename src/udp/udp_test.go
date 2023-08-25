package udp

import (
	"fmt"
	"hash/crc32"
	"testing"
)

func Test_crc32(t *testing.T) {

	receivedChecksum := crc32.ChecksumIEEE([]byte("asda你哈da阿斯达sdasd"))
	fmt.Println("receivedChecksum:", receivedChecksum)
}
