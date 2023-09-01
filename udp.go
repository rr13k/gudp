package gudp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/rr13k/gudp/crypto"
	"github.com/rr13k/gudp/encoding"
	"github.com/rr13k/gudp/gudp_protos"

	"google.golang.org/protobuf/proto"
)

var conn *net.UDPConn

// rawRecord is sent to the rawRecords channel when a new payload has been received
type rawRecord struct {
	payload []byte
	addr    *net.UDPAddr
}

type Server struct {
	protosolVersion     [2]byte       // 协议版本(预留字段)
	readBufferSize      int           // 用于设置读取缓冲区的大小，用于接收传入的字节流。
	minimumPayloadSize  int           //最小有效载荷大小，用于忽略太短的传入字节流，以防止处理空的或格式不正确的记录
	heartbeatExpiration time.Duration // Expiration time of last heartbeat to delete client
	conn                *net.UDPConn

	authClient AuthClient // 对用户令牌进行认证

	transcoder encoding.Transcoder // 解码器部分，定义接口，默认按pb实现，用户可自定义

	asymmCrypto crypto.Asymmetric // 用于解密客户端握手 hello 记录的主体部分
	symmCrypto  crypto.Symmetric  //用于在成功握手后加密和解密客户端记录的主体部分
	handler     HandlerFunc       // 用户自定义消息接收处理

	clients map[string]*Client // 管理已连接的客户端

	Error chan error // 当服务出现错误

	garbageCollectionTicker *time.Ticker       // 用于定期检查客户端的垃圾收集定时器，以清理不活跃的客户端
	garbageCollectionStop   chan bool          // 用于停止垃圾收集定时器的通道
	sessionManager          *sessionManager    // 用于生成 cookie 和 session ID
	sessions                map[string]*Client // 一个映射（Map），将客户端的 IP 地址和端口映射到客户端对象。用于快速查找客户端
	rawRecords              chan rawRecord     // 用于接收原始记录的通道，记录将被解析和处理
	logger                  *log.Logger        // 用于记录日志信息的日志记录器
	stop                    chan bool          // 用于停止服务器监听的通道
	wg                      *sync.WaitGroup    // 用于等待所有并发处理完成

	reliabilitySystem ReliabilitySystem // 用于可靠消息
}

// custom error types
var (
	ErrInvalidRecordType            = errors.New("invalid record type")
	ErrInsecureEncryptionKeySize    = errors.New("insecure encryption key size")
	ErrClientSessionNotFound        = errors.New("client session not found")
	ErrClientAddressIsNotRegistered = errors.New("client address is not registered")
	ErrClientNotFound               = errors.New("client not found")
	ErrMinimumPayloadSizeLimit      = errors.New("minimum payload size limit")
	ErrClientCookieIsInvalid        = errors.New("client cookie is invalid")
	ErrInvalidPayloadBodySize       = errors.New("payload body size is invalid")
)

const (
	HandshakeClientHelloRecordType byte = 1 << iota
	HandshakeHelloVerifyRecordType
	HandshakeServerHelloRecordType
	PingRecordType
	PongRecordType
	UnAuthenticated

	defaultRSAKeySize         int = 2048 // default RSA key size, this options is used to initiate new RSA implementation if no asymmetric encryption is passed
	defaultMinimumPayloadSize int = 3
	defaultReadBufferSize     int = 2048

	insecureSymmKeySize int = 32 // A symmetric key smaller than 256 bits is not secure. 256-bits = 32 bytes in size
)

type Option func(*Server)

func NewServerManager(conn *net.UDPConn, options ...Option) (*Server, error) {
	s := Server{
		conn: conn,

		clients:  make(map[string]*Client),
		sessions: make(map[string]*Client),

		garbageCollectionStop: make(chan bool, 1),
		stop:                  make(chan bool, 1),
		Error:                 make(chan error, 1),

		wg: &sync.WaitGroup{},

		rawRecords:        make(chan rawRecord),
		reliabilitySystem: *NewReliabilitySystem(1000000),
	}

	for _, opt := range options {
		opt(&s)
	}

	if s.readBufferSize == 0 {
		s.readBufferSize = defaultReadBufferSize
	}

	if s.minimumPayloadSize == 0 {
		s.minimumPayloadSize = defaultMinimumPayloadSize
	}

	if s.symmCrypto == nil {
		s.symmCrypto = crypto.NewAES(crypto.AES_CBC)
	}

	var err error
	if s.asymmCrypto == nil {
		s.asymmCrypto, err = crypto.NewRSA(defaultRSAKeySize)
		if err != nil {
			return nil, err
		}
	}

	if s.authClient == nil {
		s.authClient = &DefaultAuthClient{}
	}

	s.sessionManager, err = newSessionManager()
	if err != nil {
		return nil, err
	}

	if s.logger == nil { // discard logging if no logger is set
		// s.logger = log.New(io.Discard, "", 0)
		s.logger = log.New(os.Stderr, "", 0) // 输出到控制台
	}

	return &s, nil
}

func CreateUdpAddr(ip string, port int) (*net.UDPAddr, error) {
	IP := net.ParseIP(ip)

	if IP == nil {
		return nil, errors.New("miss ip addr by:" + ip)
	}

	udp_addr := &net.UDPAddr{
		IP:   IP,
		Port: port,
	}

	return udp_addr, nil
}

// 发送可靠消息
func (c *Server) sendReliablePacket(addr *net.UDPAddr, data []byte) error {
	var reliableMsg = gudp_protos.ReliableMessage{}
	reliableMsg.Seq = c.reliabilitySystem.LocalSequence()
	reliableMsg.Ack = c.reliabilitySystem.RemoteSequence()
	reliableMsg.AckBits = c.reliabilitySystem.GenerateAckBits()
	reliableMsg.Data = data

	// 发送可靠数据包
	reliablePayload, err := proto.Marshal(&reliableMsg)
	if err != nil {
		c.logger.Printf("serverHandshakeVerify proto.Marshal error: %s", err.Error())
	}

	err = c.sendReliabilityMessgae(addr, reliablePayload)
	// 处理已发送的数据包，更新相关的状态和数据结构，以及执行一些调试和测试操作
	c.reliabilitySystem.PacketSent(len(data))
	fmt.Println("发送回用户的token内容", data)
	return err
}

// HandlerFunc is called when a custom message type is received from the client
type HandlerFunc func(client *Client, p []byte)

// SetHandler sets the handler function as a callback to call when a custom record type is received from the client
func (s *Server) SetHandler(f HandlerFunc) {
	s.handler = f
}

// parseRecord 解析消息返回对应的类型
func ParseRecord(rec *[]byte) (gudp_protos.GudpMessageType, error) {
	msg_type := gudp_protos.GudpMessageType(int32((*rec)[0]))
	fmt.Println("msg_type", msg_type)
	*rec = (*rec)[1:]
	return msg_type, nil
}

// handlerRecord validate & parse incoming bytes to a record instance, then process it depends on the record type
// all incoming bytes will ignore if it doesn't meet minimum payload size (to prevent process empty or wrong formatted records)
func (s *Server) handleRecord(record []byte, addr *net.UDPAddr) {
	fmt.Println("接收到udp消息...", record)
	if len(record) < s.minimumPayloadSize {
		s.logger.Println(ErrMinimumPayloadSizeLimit)
		return
	}

	msg_type, err := ParseRecord(&record)

	if err != nil {
		s.logger.Printf("error while parsing record: %s", err)
		return
	}

	switch msg_type {
	case gudp_protos.GudpMessageType_PING:
		fmt.Println("接收ping消息", record)
		var ping = gudp_protos.Ping{}
		err := proto.Unmarshal(record, &ping)
		if err != nil {
			fmt.Println(err)
		}
		s.handlePingRecord(context.Background(), addr, &ping)

	case gudp_protos.GudpMessageType_PONG:
		fmt.Println("asdad")

	case gudp_protos.GudpMessageType_HANDSHAKEMESSAGE:
		fmt.Println("处理握手")
		var hand = gudp_protos.HandshakeMessage{}
		err := proto.Unmarshal(record, &hand)
		if err != nil {
			fmt.Println(err)
		}
		s.handleHandshakeRecord(context.Background(), addr, &hand)

	case gudp_protos.GudpMessageType_RELIABLEMESSAGE:
		// _reliableMessage := apiMsg.GetReliableMessage()
		var msg = gudp_protos.ReliableMessage{}
		fmt.Println("RELIABLEMESSAGE raw:", record)
		err := proto.Unmarshal(record, &msg)
		size := len(record)
		if err != nil {
			fmt.Println("error:", err)
		}
		// todo: 这里存在问题，如果确认消息丢包，服务端可能接收到多次内容，因为这里没处理消息已被接受情况
		// 正确的流程应该是，处理前先检查消息是否在本地被确认，以免引传输时确认消息丢失重复处理
		s.sendRawToParent(addr, msg.GetData())
		// 返回确认接受消息
		s.reliabilityPacketReceived(&msg, size, addr)

	case gudp_protos.GudpMessageType_UNRELIABLEMESSAGE:
		var msg = gudp_protos.UnreliableMessage{}
		_ = proto.Unmarshal(record, &msg)
		fmt.Println("接收不可靠")
		s.sendRawToParent(addr, msg.GetData())
	default:
		fmt.Println("用户自定义消息")
		// 处理用户发送的消息，将原始数据发送给用户自己解析
		// 底层不处理session, 没必要，上层自己构建session
		// s.handleCustomRecord(context.Background(), addr, nil)
	}
}

// 接收到可靠数据包，确认消息为单向的不用反复确认
func (s *Server) reliabilityPacketReceived(msg *gudp_protos.ReliableMessage, size int, addr *net.UDPAddr) int {
	fmt.Println("已发送:", len(s.reliabilitySystem.sentQueue), "待确认包:", len(s.reliabilitySystem.pendingAckQueue), "已接收包:", s.reliabilitySystem.recvPackets, "已确认包:", len(s.reliabilitySystem.ackedQueue))
	// 如果已接收队列中已存在当前序列号的数据包，则不进行处理
	s.reliabilitySystem.PacketReceived(msg.Seq, size)

	// 这里，只有接收到确认消息时才有用
	s.reliabilitySystem.ProcessAck(msg.Ack, msg.AckBits)
	// 当size 为 0 表示收到的客户端确认消息
	if size != 0 { // 收到客户端消息，返回确认消息
		msg.Data = []byte{}
		msg.Ack = msg.Seq
		fmt.Println("back ack ->", msg.Seq)
		data, _ := proto.Marshal(msg)
		// 发送确认消息给客户端
		s.sendReliabilityMessgae(addr, data)
	}
	// 在这里对客户端的传递的消息进行确认。在接收时可以根据ack来进行区分是否为确认消息。
	return size
}

// 将消息传递给用户响应函数
func (s *Server) sendRawToParent(addr *net.UDPAddr, bytes []byte) {
	cl, ok := s.findClientByAddr(addr)
	if !ok {
		s.logger.Printf("error while authenticating other type record: %s", ErrClientAddressIsNotRegistered)
		s.unAuthenticated(addr)
		return
	}

	if s.handler != nil {
		s.handler(cl, bytes)
	}
}

func (s *Server) unAuthenticated(addr *net.UDPAddr) {
	payload := composeRecordBytes(UnAuthenticated, s.protosolVersion, []byte{})
	err := s.sendSockt(addr, payload)
	if err != nil {
		s.logger.Printf("error while sending UnAuthenticated record to the client: %s", err)
		return
	}
}

// 解析出sessionid， 这里不该从这里解，因为全部使用pb
// parseSessionID parses the session ID from the record decrypted body, the session ID must prepend to the body before encryption in the client
func parseSessionID(p []byte, sLen int) ([]byte, []byte, error) {
	if len(p) < sLen {
		return nil, nil, ErrInvalidPayloadBodySize
	}
	return p[:sLen], p[sLen:], nil
}

// handleCustomRecord handle custom record with authorizing the record and call the handler func if is set
func (s *Server) handleCustomRecord(ctx context.Context, addr *net.UDPAddr, r *record) {
	cl, ok := s.findClientByAddr(addr)
	if !ok {
		s.logger.Printf("error while authenticating other type record: %s", ErrClientAddressIsNotRegistered)
		s.unAuthenticated(addr)
		return
	}

	// 解密
	payload, err := s.symmCrypto.Decrypt(r.Body, cl.eKey)
	if err != nil {
		s.logger.Printf("error while decrypting other type record: %s", err)
		return
	}

	// 从发送内容中解出session
	var sessionID, body []byte
	sessionID, body, err = parseSessionID(payload, len(cl.sessionID))
	if err != nil {
		s.logger.Printf("error while parsing session id for ping: %s", err)
		return
	}

	// 校验session是否有效
	if !cl.ValidateSessionID(sessionID) {
		s.logger.Printf("error while validating client session for other type record: %s", ErrClientSessionNotFound)
		s.unAuthenticated(addr)
		return
	}

	// 将消息传递给用户
	if s.handler != nil {
		s.handler(cl, body)
	}

	now := time.Now()
	cl.Lock()
	cl.lastHeartbeat = &now
	cl.Unlock()
}

// registerClient generates a new session ID & registers an address with token ID & encryption key as a Client
func (s *Server) registerClient(addr *net.UDPAddr, ID string, eKey []byte) (*Client, error) {
	sessionID, err := s.sessionManager.GenerateSessionID(addr, ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cl := &Client{
		ID:            ID,
		sessionID:     sessionID,
		addr:          addr,
		eKey:          eKey,
		lastHeartbeat: &now,
	}
	s.clients[ID] = cl
	s.sessions[fmt.Sprintf("%s_%d", addr.IP.String(), addr.Port)] = cl

	return cl, nil
}

// 处理握手的，保证可靠
// handleHandshakeRecord handles handshake process
// HandshakeClientHello record is encrypted by the server public key & contains the client encryption key
// if the ClientHello was valid, the server generates a unique cookie for the client address, encrypt it with the client key & then send it
// client must send the HandshakeClientHelloVerify request (same as Hello) with the generated cookie & the token to prove that the sender address is valid
// server validate the HelloVerify record, then authenticate the client token & if they're valid, generate a session ID, encrypt it & send it back as ServerHello record
// after client registration, the client must prepend the Session ID before the record body unencrypted bytes, then encrypt them & compose the record
func (s *Server) handleHandshakeRecord(ctx context.Context, addr *net.UDPAddr, r *gudp_protos.HandshakeMessage) {
	// var payload []byte

	// 先不加密，后续在说，开发效率太慢
	// payload, err := s.asymmCrypto.Decrypt(r.Body)
	// if err != nil {
	// 	s.logger.Printf("error while decrypting record body: %s", err)
	// 	return
	// }

	//TODO validate the client random
	// 增加随机数防止重放攻击

	if len(r.GetCookie()) == 0 {
		cookie := s.sessionManager.GetAddrCookieHMAC(addr, []byte(r.GetClientVersion()), []byte(r.GetSessionId()), r.GetRandom()) //TODO session id is empty

		if len(r.GetKey()) < insecureSymmKeySize {
			s.logger.Printf("error while parsing ClientHello record: %s", ErrInsecureEncryptionKeySize)
			return
		}

		r.Cookie = cookie
		r.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
		handshakePayload, err := proto.Marshal(r)
		if err != nil {
			s.logger.Printf("serverHandshakeVerify proto.Marshal error: %s", err.Error())
		}

		// 加密
		// handshakePayload, err = s.symmCrypto.Encrypt(handshakePayload, r.GetKey())
		// if err != nil {
		// 	s.logger.Printf("error while encrypting HelloVerify record: %s", err)
		// 	return
		// }

		err = s.sendHandMessgaeToAddr(addr, handshakePayload)
		if err != nil {
			s.logger.Printf("error while sending HelloVerify record to the client: %s", err)
			return
		}
	} else {
		fmt.Println("收到携带cookie握手")
		cookie := s.sessionManager.GetAddrCookieHMAC(addr, []byte(r.GetClientVersion()), []byte(r.GetSessionId()), r.GetRandom()) //TODO session id is empty
		if !crypto.HMACEqual(r.GetCookie(), cookie) {
			s.logger.Printf("error while validation HelloVerify record cookie: %s", ErrClientCookieIsInvalid)
			return
		}

		if len(r.GetKey()) < insecureSymmKeySize {
			s.logger.Printf("error while validating HelloVerify record key: %s", ErrInsecureEncryptionKeySize)
			return
		}

		var token []byte
		var err error
		if len(r.Extra) > 0 {
			token, err = s.symmCrypto.Decrypt(r.Extra, r.GetKey())
			if err != nil {
				s.logger.Printf("error while decrypting HelloVerify record token: %s", err)
				return
			}
		}

		var ID string
		ID, err = s.authClient.Authenticate(ctx, token)
		if err != nil {
			s.logger.Printf("error while authenticating client token: %s", err)
			return
		}

		var cl *Client
		cl, err = s.registerClient(addr, ID, r.GetKey())
		if err != nil {
			s.logger.Printf("error while registering client: %s", err)
			return
		}

		serverHandshakeVerify := &gudp_protos.HandshakeMessage{}
		serverHandshakeVerify.SessionId = cl.sessionID
		serverHandshakeVerify.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
		serverHandshakeHelloBytes, _ := proto.Marshal(serverHandshakeVerify)

		err = s.sendHandMessgaeToAddr(addr, serverHandshakeHelloBytes)
		if err != nil {
			s.logger.Printf("error while sending server hello record: %s", err)
			return
		}
	}
}

// returns the Client by IP & Port
func (s *Server) findClientByAddr(addr *net.UDPAddr) (*Client, bool) {
	cl, ok := s.sessions[fmt.Sprintf("%s_%d", addr.IP.String(), addr.Port)]
	if !ok {
		return nil, ok
	}

	return cl, true
}

// handlePingRecord handles ping record and sends pong response
func (s *Server) handlePingRecord(ctx context.Context, addr *net.UDPAddr, r *gudp_protos.Ping) {
	cl, ok := s.findClientByAddr(addr)
	if !ok {
		s.logger.Printf("error while authenticating ping record: %s", ErrClientAddressIsNotRegistered)
		return
	}

	pong := &gudp_protos.Pong{}
	pong.ReceivedAt = time.Now().UnixNano()
	pong.SentAt = time.Now().UnixNano()
	pong.PingSentAt = r.SentAt

	pingBytes, err := proto.Marshal(pong)
	if err != nil {
		fmt.Println("ping proto.Marshal error:", err)
	}
	// 如有必要后续可校验session
	err = s.sendToClient(cl, pingBytes)
	if err != nil {
		s.logger.Printf("error while sending pong record: %s", err)
		return
	}

	now := time.Now()
	cl.Lock()
	cl.lastHeartbeat = &now
	cl.Unlock()
}

// incoming bytes are parsed to the record struct
type record struct {
	Type       byte
	ProtoMajor uint8
	ProtoMinor uint8
	Body       []byte
	Extra      []byte
}

func (s *Server) clientGarbageCollection() {
	for {
		select {
		case <-s.garbageCollectionStop:
			if s.garbageCollectionTicker != nil {
				s.garbageCollectionTicker.Stop()
			}
		case <-s.garbageCollectionTicker.C:
			for _, c := range s.clients {
				if c.lastHeartbeat != nil && time.Now().After(c.lastHeartbeat.Add(s.heartbeatExpiration)) {
					delete(s.clients, c.ID)
					delete(s.sessions, fmt.Sprintf("%s_%d", c.addr.IP.String(), c.addr.Port))
				}
			}
		}
	}
}

func (s *Server) handleRawRecords() {
	for r := range s.rawRecords {
		s.handleRecord(r.payload, r.addr)
	}
}

// Serve starts listening to the UDP port for incoming bytes & then sends payload and sender address into the rawRecords channel if no error is found
func (s *Server) Serve() {
	// 如果设置了心跳，则按照心跳时间用定时器检查客户端连接
	if s.heartbeatExpiration > 0 {
		if s.garbageCollectionTicker != nil {
			s.garbageCollectionTicker.Stop()
		}
		s.garbageCollectionTicker = time.NewTicker(s.heartbeatExpiration)
		s.garbageCollectionStop = make(chan bool, 1)
		go s.clientGarbageCollection()
	}

	// 通过信号处理接受消息
	s.rawRecords = make(chan rawRecord)
	go s.handleRawRecords()

	// 设置读取超时为空，一直等待数据读取
	err := s.conn.SetReadDeadline(time.Time{}) // reset read deadline
	if err != nil {
		fmt.Println("Read timeout error:", err)
	}
	s.stop = make(chan bool, 1) // reset the stop channel
	for {
		select {
		case <-s.stop:
			return
		default:
			buf := make([]byte, s.readBufferSize)
			n, addr, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					continue
				}

				s.logger.Printf("error while reading from udp: %s", err)
				continue
			}
			s.rawRecords <- rawRecord{
				payload: buf[0:n],
				addr:    addr,
			}
		}
	}
}

// 发送握手消息
func (s *Server) sendHandMessgaeToAddr(addr *net.UDPAddr, record []byte) error {
	_type := byte(gudp_protos.GudpMessageType_HANDSHAKEMESSAGE)
	recordWithPrefix := append([]byte{_type}, record...)
	return s.sendSockt(addr, recordWithPrefix)
}

// 将消息发送到指定地址
func (s *Server) SendMessage(addr *net.UDPAddr, bytes []byte, reliability bool) error {
	if reliability {
		return s.sendReliablePacket(addr, bytes)
	} else {
		return s.sendUnreliabilityMessgae(addr, bytes)
	}
}

// 将消息发送到指定地址
func (s *Server) SendClientMessage(cli *Client, bytes []byte, reliability bool) error {
	return s.SendMessage(cli.addr, bytes, reliability)
}

// 将消息发送到指定地址
func (s *Server) sendSockt(addr *net.UDPAddr, bytes []byte) error {
	_, err := s.conn.WriteToUDP(bytes, addr)
	return err
}

// 发送可靠消息
func (s *Server) sendReliabilityMessgae(addr *net.UDPAddr, record []byte) error {
	_type := byte(gudp_protos.GudpMessageType_RELIABLEMESSAGE)
	recordWithPrefix := append([]byte{_type}, record...)
	return s.sendSockt(addr, recordWithPrefix)
}

// 发送不可靠消息
func (s *Server) sendUnreliabilityMessgae(addr *net.UDPAddr, record []byte) error {
	_type := byte(gudp_protos.GudpMessageType_UNRELIABLEMESSAGE)
	recordWithPrefix := append([]byte{_type}, record...)
	return s.sendSockt(addr, recordWithPrefix)
}

// sends a record byte array to the Client. the record type is prepended to the record body as a byte
func (s *Server) sendToClient(client *Client, payload []byte) error {
	payload, err := s.symmCrypto.Encrypt(payload, client.eKey)
	if err != nil {
		return err
	}
	return s.sendSockt(client.addr, payload)
}

// composeRecordBytes composes record bytes, prepend the record header (type & protosol version) to the body
func composeRecordBytes(typ byte, version [2]byte, payload []byte) []byte {
	return append([]byte{typ, version[0], version[1]}, payload...)
}

func CreateUdpConn(host string, port int) *net.UDPConn {
	// 设置服务器监听地址和端口
	// 创建一个 UDP 地址
	serverAddr, err := CreateUdpAddr(host, port)
	if err != nil {
		fmt.Println(err)
	}

	// 创建 UDP 连接
	conn, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}

	fmt.Printf("UDP server listening on %s success~\n", serverAddr)

	return conn
}

// WithProtosolVersion sets the server protosol version
func WithProtosolVersion(major, minor uint8) Option {
	return func(s *Server) {
		s.protosolVersion = [2]byte{major, minor}
	}
}

// WithHeartbeatExpiration sets the server heartbeat expiration option
func WithHeartbeatExpiration(t time.Duration) Option {
	return func(s *Server) {
		s.heartbeatExpiration = t
	}
}

// WithMinimumPayloadSize sets the minimum payload size option
func WithMinimumPayloadSize(i int) Option {
	return func(s *Server) {
		s.minimumPayloadSize = i
	}
}

// WithReadBufferSize sets the read buffer size option
func WithReadBufferSize(i int) Option {
	return func(s *Server) {
		s.readBufferSize = i
	}
}

// WithSymmetricCrypto sets the symmetric cryptography implementation
func WithSymmetricCrypto(sc crypto.Symmetric) Option {
	return func(s *Server) {
		s.symmCrypto = sc
	}
}

// WithAsymmetricCrypto sets the asymmetric cryptography implementation
func WithAsymmetricCrypto(ac crypto.Asymmetric) Option {
	return func(s *Server) {
		s.asymmCrypto = ac
	}
}

// WithTranscoder sets the transcoder implementation
func WithTranscoder(t encoding.Transcoder) Option {
	return func(s *Server) {
		s.transcoder = t
	}
}

// WithAuthClient sets the auth client implementation
func WithAuthClient(ac AuthClient) Option {
	return func(s *Server) {
		s.authClient = ac
	}
}

// WithLogger sets the logger
func WithLogger(l *log.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}
