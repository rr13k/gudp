package gudp

import (
	"fmt"
	"hash/crc32"
	"testing"
)

func Test_crc32(t *testing.T) {

	receivedChecksum := crc32.ChecksumIEEE([]byte("asda你哈da阿斯达sdasd"))
	fmt.Println("receivedChecksum:", receivedChecksum)
}

// 测试当值不在枚举范围内情况
func Test_parseEnum(t *testing.T) {
	// aa, _ := ParseRecord([]byte{12})

	// fmt.Println(aa.String())
}
