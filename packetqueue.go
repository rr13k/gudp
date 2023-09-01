package gudp

import "time"

// packet queue to store information about sent and received packets sorted in
// sequence order + we define ordering using the "sequenceMoreRecent" function,
// this works provided there is a large gap when sequence wrap occurs

type PacketData struct {
	sequence uint32        // 包seq序号
	time     time.Duration // 超时时间
	size     int           // 包长度
	re       int           // 重传次数
	bytes    []byte        // 包数据
}

// TODO: check if its better to have a slice of pointers to PacketData
type PacketQueue []PacketData

func (pq *PacketQueue) Exists(sequence uint32) bool {
	for i := 0; i < len(*pq); i++ {
		if (*pq)[i].sequence == sequence {
			return true
		}
	}
	return false
}

func (pq *PacketQueue) InsertSorted(p PacketData, maxSequence uint32) {
	if len(*pq) == 0 {
		*pq = append(*pq, p)
	} else {
		if !sequenceMoreRecent(p.sequence, (*pq)[0].sequence, maxSequence) {
			*pq = append(PacketQueue{p}, *pq...)
		} else if sequenceMoreRecent(p.sequence, (*pq)[len(*pq)-1].sequence, maxSequence) {
			*pq = append(*pq, p)
		} else {
			for i := 0; i < len(*pq); i++ {
				if sequenceMoreRecent((*pq)[i].sequence, p.sequence, maxSequence) {
					(*pq) = append((*pq)[:i], append(PacketQueue{p}, (*pq)[i:]...)...)
					i++
				}
			}
		}
	}
}

//
// utility functions
//

// 这个函数用于比较两个序列号的先后关系，判断哪个序列号更近（即更新）。在序列号循环的情况下，例如使用32位无符号整数表示的序列号，当序列号达到最大值后会从0重新开始。这时，判断序列号的先后关系需要考虑循环的情况。
func sequenceMoreRecent(s1, s2, maxSequence uint32) bool {
	return (s1 > s2) && (s1-s2 <= maxSequence/2) || (s2 > s1) && (s2-s1 > maxSequence/2)
}

func bitIndexForSequence(sequence, ack, maxSequence uint32) uint32 {
	// TODO: remove those asserts once done
	if sequence == ack {
		panic("assert(sequence != ack)")
	}
	if sequenceMoreRecent(sequence, ack, maxSequence) {
		panic("assert(!sequenceMoreRecent(sequence, ack, maxSequence))")
	}
	if sequence > ack {
		if ack >= 33 {
			panic("assert(ack < 33)")
		}
		if maxSequence < sequence {
			panic("assert(maxSequence >= sequence)")
		}
		return ack + (maxSequence - sequence)
	}
	if ack < 1 {
		panic("assert(ack >= 1)")
	}
	if sequence > ack-1 {
		panic("assert(sequence <= ack-1)")
	}
	return ack - 1 - sequence
}

// 确认位图是由接收端生成的，用于告诉发送端哪些数据包已经被接收。发送端根据这个确认位图来确定哪些数据包已经安全地到达了接收端
func generateAckBits(ack uint32, receivedQueue *PacketQueue, maxSequence uint32) uint32 {
	var ackBits uint32
	for itor := 0; itor < len(*receivedQueue); itor++ {
		iseq := (*receivedQueue)[itor].sequence

		if iseq == ack || sequenceMoreRecent(iseq, ack, maxSequence) {
			break
		}
		bitIndex := bitIndexForSequence(iseq, ack, maxSequence)
		if bitIndex <= 31 {
			ackBits |= 1 << bitIndex
		}
	}
	return ackBits
}

func processAck(ack, ackBits uint32,
	pendingAckQueue, ackedQueue *PacketQueue,
	acks *[]uint32, ackedPackets *uint32,
	rtt *time.Duration, maxSequence uint32) {
	if len(*pendingAckQueue) == 0 {
		return
	}

	i := 0
	for i < len(*pendingAckQueue) {
		var acked bool
		itor := &(*pendingAckQueue)[i]

		if itor.sequence == ack {
			acked = true
		} else if !sequenceMoreRecent(itor.sequence, ack, maxSequence) {
			bitIndex := bitIndexForSequence(itor.sequence, ack, maxSequence)
			if bitIndex <= 31 {
				acked = ((ackBits >> bitIndex) & 1) != 0
			}
		}

		if acked {
			(*rtt) += (itor.time - *rtt) / 10

			ackedQueue.InsertSorted(*itor, maxSequence)
			*acks = append(*acks, itor.sequence)
			*ackedPackets++
			//itor = pending_ack_queue.erase( itor );
			*pendingAckQueue = append((*pendingAckQueue)[:i], (*pendingAckQueue)[i+1:]...)
		} else {
			i++
		}
	}
}
