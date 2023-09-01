package gudp

import (
	"fmt"
	"time"
)

// 进行可靠性保障
// maxSequence: 最大的序列号值，用于在序列号循环时进行比较。当序列号超过该值时，会重新从 0 开始。这用于处理序列号的循环问题。

// localSequence: 本地的序列号，用于标识发送的数据包的序列号。每次发送数据包后，此值会递增，当达到 maxSequence 时会重置为 0。

// remoteSequence: 远程的序列号，表示最近接收到的对方发送的数据包的序列号。用于追踪已接收的数据包。

// sentPackets: 总共发送的数据包数量。

// recvPackets: 总共接收的数据包数量。

// lostPackets: 丢失的数据包数量，通过比较已发送队列和已确认队列，可以计算出未确认的丢失数据包数量。

// ackedPackets: 已确认的数据包数量，即已收到确认的数据包数量。

// sentBandwidth: 近似的每秒发送带宽，以比特/秒为单位。

// ackedBandwidth: 近似的每秒确认带宽，以比特/秒为单位。

// rtt: 估计的往返时延（Round-Trip Time），表示从发送数据包到接收确认的时间间隔。

// rttMax: 估计的最大往返时延，用于设置丢失数据包的时间阈值。

// acks: 已确认的数据包的序列号列表，用于发送确认给对方。

// sentQueue: 已发送数据包队列，包括等待确认的数据包。

// pendingAckQueue: 待确认的数据包队列，包括等待对方确认的数据包。

// receivedQueue: 已接收数据包队列，用于追踪已接收的数据包，以便生成确认信息。

// ackedQueue: 已确认的数据包队列，用于跟踪已经确认的数据包，以计算已确认带宽等统计信息。

// 这些参数共同构成了可靠传输系统的核心部分，实现了数据包的确认、重传、丢失恢复和统计信息收集等功能。通过综合运用这些参数，系统能够保障数据包的可靠传输和相关的传输状态跟踪。
// reliability system to support reliable connection
//   - manages sent, received, pending ack and acked packet queues
//   - separated out from reliable connection so they can be unit-tested
type ReliabilitySystem struct {
	maxSequence    uint32 // maximum sequence value before wrap around (used to test sequence wrap at low # values)
	localSequence  uint32 // local sequence number for most recently sent packet
	remoteSequence uint32 // remote sequence number for most recently received packet

	sentPackets  uint32 // total number of packets sent
	recvPackets  uint32 // total number of packets received
	lostPackets  uint32 // total number of packets lost
	ackedPackets uint32 // total number of packets acked

	sentBandwidth  float64       // approximate sent bandwidth over the last second
	ackedBandwidth float64       // approximate acked bandwidth over the last second
	rtt            time.Duration // estimated round trip time
	rttMax         time.Duration // maximum expected round trip time (hard coded to one second for the moment)

	acks []uint32 // acked packets from last set of packet receives. cleared each update!

	sentQueue       PacketQueue // 已发送数据包队列，包括等待确认的数据包。
	pendingAckQueue PacketQueue // 待确认的数据包队列，包括等待对方确认的数据包。
	receivedQueue   PacketQueue // 已接收数据包队列，用于追踪已接收的数据包，以便生成确认信息。
	ackedQueue      PacketQueue // 已确认的数据包队列，用于跟踪已经确认的数据包，以计算已确认带宽等统计信息。
}

func NewReliabilitySystem(maxSequence uint32) *ReliabilitySystem {
	rs := &ReliabilitySystem{
		rttMax:      1 * time.Second,
		maxSequence: maxSequence,
	}
	rs.Reset()
	return rs
}

func (rs *ReliabilitySystem) Reset() {
	rs.localSequence = 0
	rs.remoteSequence = 0
	rs.sentQueue = PacketQueue{}
	rs.receivedQueue = PacketQueue{}
	rs.pendingAckQueue = PacketQueue{}
	rs.ackedQueue = PacketQueue{}
	rs.sentPackets = 0
	rs.recvPackets = 0
	rs.lostPackets = 0
	rs.ackedPackets = 0
	rs.sentBandwidth = 0.0
	rs.ackedBandwidth = 0.0
	rs.rtt = 0 * time.Second
	rs.rttMax = 1 * time.Second
}

// 验证收到的包
func (rs *ReliabilitySystem) PacketSent(size int) {
	// todo: 判断 接受到的包是否已存在，如果存在再接收到则报错。需要注意客户端每次序号从0开始，如果重启1次就很容易导致这个问题
	print("rs.localSequence:", rs.localSequence)
	if rs.sentQueue.Exists(rs.localSequence) {
		fmt.Printf("local sequence %d exists\n", rs.localSequence)
		for i := 0; i < len(rs.sentQueue); i++ {
			fmt.Printf(" + %d\n", rs.sentQueue[i].sequence)
		}
		panic("assert( !sentQueue.exists( localSequence ) )")
	}
	if rs.pendingAckQueue.Exists(rs.localSequence) {
		panic("assert( !pendingAckQueue.exists( localSequence ) )")
	}

	var data PacketData
	data.sequence = rs.localSequence
	data.time = 0
	data.size = size
	rs.sentQueue = append(rs.sentQueue, data)
	rs.pendingAckQueue = append(rs.pendingAckQueue, data)
	rs.sentPackets++
	rs.localSequence++
	if rs.localSequence > rs.maxSequence {
		rs.localSequence = 0
	}
}

func (rs *ReliabilitySystem) PacketReceived(sequence uint32, size int) {
	rs.recvPackets++
	if rs.receivedQueue.Exists(sequence) {
		return
	}
	var data PacketData
	data.sequence = sequence
	data.time = 0
	data.size = size

	rs.receivedQueue = append(rs.receivedQueue, data)
	if sequenceMoreRecent(sequence, rs.remoteSequence, rs.maxSequence) {
		rs.remoteSequence = sequence
	}
}

func (rs *ReliabilitySystem) GenerateAckBits() uint32 {
	return generateAckBits(rs.remoteSequence, &rs.receivedQueue, rs.maxSequence)
}

func (rs *ReliabilitySystem) ProcessAck(ack, ackBits uint32) {
	processAck(ack, ackBits, &rs.pendingAckQueue, &rs.ackedQueue, &rs.acks, &rs.ackedPackets, &rs.rtt, rs.maxSequence)
}

func (rs *ReliabilitySystem) Update(deltaTime time.Duration) {
	rs.acks = []uint32{}
	rs.advanceQueueTime(deltaTime)
	rs.updateQueues()
	rs.updateStats()
}

// data accessors

func (rs *ReliabilitySystem) LocalSequence() uint32 {
	return rs.localSequence
}

func (rs *ReliabilitySystem) RemoteSequence() uint32 {
	return rs.remoteSequence
}

func (rs *ReliabilitySystem) MaxSequence() uint32 {
	return rs.maxSequence
}

func (rs *ReliabilitySystem) Acks() []uint32 {
	return rs.acks
}

func (rs *ReliabilitySystem) SentPackets() uint32 {
	return rs.sentPackets
}

func (rs *ReliabilitySystem) ReceivedPackets() uint32 {
	return rs.recvPackets
}

func (rs *ReliabilitySystem) LostPackets() uint32 {
	return rs.lostPackets
}

func (rs *ReliabilitySystem) AckedPackets() uint32 {
	return rs.ackedPackets
}

func (rs *ReliabilitySystem) SentBandwidth() float64 {
	return rs.sentBandwidth
}

func (rs *ReliabilitySystem) AckedBandwidth() float64 {
	return rs.ackedBandwidth
}

func (rs *ReliabilitySystem) RoundTripTime() time.Duration {
	return rs.rtt
}

func (rs *ReliabilitySystem) HeaderSize() int {
	return 12
}

func (rs *ReliabilitySystem) advanceQueueTime(deltaTime time.Duration) {
	for i := 0; i < len(rs.sentQueue); i++ {
		rs.sentQueue[i].time += deltaTime
	}
	for i := 0; i < len(rs.receivedQueue); i++ {
		rs.receivedQueue[i].time += deltaTime
	}
	for i := 0; i < len(rs.pendingAckQueue); i++ {
		rs.pendingAckQueue[i].time += deltaTime
	}
	for i := 0; i < len(rs.ackedQueue); i++ {
		rs.ackedQueue[i].time += deltaTime
	}
}

func (rs *ReliabilitySystem) updateQueues() {
	const epsilon = 1 * time.Millisecond

	for len(rs.sentQueue) > 0 && rs.sentQueue[0].time > rs.rttMax+epsilon {
		// pop front
		rs.sentQueue = rs.sentQueue[1:]
	}

	if len(rs.receivedQueue) > 0 {
		latestSequence := rs.receivedQueue[len(rs.receivedQueue)-1].sequence

		var minSequence uint32
		if latestSequence >= 34 {
			minSequence = (latestSequence - 34)
		} else {
			minSequence = rs.maxSequence - (34 - latestSequence)
		}

		for len(rs.receivedQueue) > 0 && !sequenceMoreRecent(rs.receivedQueue[0].sequence, minSequence, rs.maxSequence) {
			// pop front
			rs.receivedQueue = rs.receivedQueue[1:]
		}
	}

	for len(rs.ackedQueue) > 0 && rs.ackedQueue[0].time > rs.rttMax*2-epsilon {
		// pop front
		rs.ackedQueue = rs.ackedQueue[1:]
	}

	for len(rs.pendingAckQueue) > 0 && rs.pendingAckQueue[0].time > rs.rttMax+epsilon {
		// pop front
		rs.pendingAckQueue = rs.pendingAckQueue[1:]
		rs.lostPackets++
	}
}

func (rs *ReliabilitySystem) updateStats() {
	var sentBytesPerSecond float64
	for i := 0; i < len(rs.sentQueue); i++ {
		sentBytesPerSecond += float64(rs.sentQueue[i].size)
	}
	var (
		ackedPacketsPerSecond float64
		ackedBytesPerSecond   float64
	)
	for i := 0; i < len(rs.ackedQueue); i++ {
		if rs.ackedQueue[i].time >= rs.rttMax {
			ackedPacketsPerSecond++
			ackedBytesPerSecond += float64(rs.ackedQueue[i].size)
		}
	}
	sentBytesPerSecond /= float64(rs.rttMax)
	ackedBytesPerSecond /= float64(rs.rttMax)
	rs.sentBandwidth = sentBytesPerSecond * (8.0 / 1000.0)
	rs.ackedBandwidth = ackedBytesPerSecond * (8.0 / 1000.0)
}
