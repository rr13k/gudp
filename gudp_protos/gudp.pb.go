// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.8
// source: gudp.proto

package gudp_protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 因godot不支持自定义类型的oneof, 因此消息增加一位通过枚举区分
type GudpMessageType int32

const (
	GudpMessageType_GUDP_TYPE         GudpMessageType = 0
	GudpMessageType_PING              GudpMessageType = 1
	GudpMessageType_PONG              GudpMessageType = 2
	GudpMessageType_HANDSHAKEMESSAGE  GudpMessageType = 3
	GudpMessageType_UNRELIABLEMESSAGE GudpMessageType = 4
	GudpMessageType_RELIABLEMESSAGE   GudpMessageType = 5
	GudpMessageType_RPCMESSAGE        GudpMessageType = 6
)

// Enum value maps for GudpMessageType.
var (
	GudpMessageType_name = map[int32]string{
		0: "GUDP_TYPE",
		1: "PING",
		2: "PONG",
		3: "HANDSHAKEMESSAGE",
		4: "UNRELIABLEMESSAGE",
		5: "RELIABLEMESSAGE",
		6: "RPCMESSAGE",
	}
	GudpMessageType_value = map[string]int32{
		"GUDP_TYPE":         0,
		"PING":              1,
		"PONG":              2,
		"HANDSHAKEMESSAGE":  3,
		"UNRELIABLEMESSAGE": 4,
		"RELIABLEMESSAGE":   5,
		"RPCMESSAGE":        6,
	}
)

func (x GudpMessageType) Enum() *GudpMessageType {
	p := new(GudpMessageType)
	*p = x
	return p
}

func (x GudpMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GudpMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_gudp_proto_enumTypes[0].Descriptor()
}

func (GudpMessageType) Type() protoreflect.EnumType {
	return &file_gudp_proto_enumTypes[0]
}

func (x GudpMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GudpMessageType.Descriptor instead.
func (GudpMessageType) EnumDescriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{0}
}

type ApiCharacter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId   string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Scene      string `protobuf:"bytes,2,opt,name=scene,proto3" json:"scene,omitempty"`
	Name       string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Career     string `protobuf:"bytes,4,opt,name=career,proto3" json:"career,omitempty"`
	Level      int32  `protobuf:"varint,5,opt,name=level,proto3" json:"level,omitempty"`
	Experience int64  `protobuf:"varint,6,opt,name=experience,proto3" json:"experience,omitempty"` // Position position = 7;
}

func (x *ApiCharacter) Reset() {
	*x = ApiCharacter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiCharacter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiCharacter) ProtoMessage() {}

func (x *ApiCharacter) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiCharacter.ProtoReflect.Descriptor instead.
func (*ApiCharacter) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{0}
}

func (x *ApiCharacter) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *ApiCharacter) GetScene() string {
	if x != nil {
		return x.Scene
	}
	return ""
}

func (x *ApiCharacter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ApiCharacter) GetCareer() string {
	if x != nil {
		return x.Career
	}
	return ""
}

func (x *ApiCharacter) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *ApiCharacter) GetExperience() int64 {
	if x != nil {
		return x.Experience
	}
	return 0
}

// 通过proto实现可靠udp
type ReliableMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SequenceNumber int32  `protobuf:"varint,1,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"` // 序列号
	Data           []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Checksum       uint32 `protobuf:"varint,3,opt,name=checksum,proto3" json:"checksum,omitempty"` // 校验码
	SessionId      string `protobuf:"bytes,4,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Seq            uint32 `protobuf:"varint,5,opt,name=seq,proto3" json:"seq,omitempty"`         // 本地序列号
	Ack            uint32 `protobuf:"varint,6,opt,name=ack,proto3" json:"ack,omitempty"`         // 远程序列号
	AckBits        uint32 `protobuf:"varint,7,opt,name=ackBits,proto3" json:"ackBits,omitempty"` //确认位图
}

func (x *ReliableMessage) Reset() {
	*x = ReliableMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReliableMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReliableMessage) ProtoMessage() {}

func (x *ReliableMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReliableMessage.ProtoReflect.Descriptor instead.
func (*ReliableMessage) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{1}
}

func (x *ReliableMessage) GetSequenceNumber() int32 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *ReliableMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ReliableMessage) GetChecksum() uint32 {
	if x != nil {
		return x.Checksum
	}
	return 0
}

func (x *ReliableMessage) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *ReliableMessage) GetSeq() uint32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *ReliableMessage) GetAck() uint32 {
	if x != nil {
		return x.Ack
	}
	return 0
}

func (x *ReliableMessage) GetAckBits() uint32 {
	if x != nil {
		return x.AckBits
	}
	return 0
}

// 处理握手进行身份验证,保障连接是有效的
type HandshakeMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cookie        []byte `protobuf:"bytes,1,opt,name=cookie,proto3" json:"cookie,omitempty"`
	Random        []byte `protobuf:"bytes,2,opt,name=random,proto3" json:"random,omitempty"`
	Key           []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Timestamp     int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	SessionId     []byte `protobuf:"bytes,5,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	ClientVersion string `protobuf:"bytes,6,opt,name=clientVersion,proto3" json:"clientVersion,omitempty"`
	Extra         []byte `protobuf:"bytes,7,opt,name=extra,proto3" json:"extra,omitempty"`
}

func (x *HandshakeMessage) Reset() {
	*x = HandshakeMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeMessage) ProtoMessage() {}

func (x *HandshakeMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeMessage.ProtoReflect.Descriptor instead.
func (*HandshakeMessage) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{2}
}

func (x *HandshakeMessage) GetCookie() []byte {
	if x != nil {
		return x.Cookie
	}
	return nil
}

func (x *HandshakeMessage) GetRandom() []byte {
	if x != nil {
		return x.Random
	}
	return nil
}

func (x *HandshakeMessage) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *HandshakeMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *HandshakeMessage) GetSessionId() []byte {
	if x != nil {
		return x.SessionId
	}
	return nil
}

func (x *HandshakeMessage) GetClientVersion() string {
	if x != nil {
		return x.ClientVersion
	}
	return ""
}

func (x *HandshakeMessage) GetExtra() []byte {
	if x != nil {
		return x.Extra
	}
	return nil
}

// 不可靠的常规udp
type UnreliableMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data      []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	SessionId string `protobuf:"bytes,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *UnreliableMessage) Reset() {
	*x = UnreliableMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnreliableMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnreliableMessage) ProtoMessage() {}

func (x *UnreliableMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnreliableMessage.ProtoReflect.Descriptor instead.
func (*UnreliableMessage) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{3}
}

func (x *UnreliableMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *UnreliableMessage) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type RpcMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RpcId  int32  `protobuf:"varint,1,opt,name=rpc_id,json=rpcId,proto3" json:"rpc_id,omitempty"`
	Method string `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Data   []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RpcMessage) Reset() {
	*x = RpcMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcMessage) ProtoMessage() {}

func (x *RpcMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcMessage.ProtoReflect.Descriptor instead.
func (*RpcMessage) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{4}
}

func (x *RpcMessage) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

func (x *RpcMessage) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *RpcMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ApiMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to NoticeWay:
	//
	//	*ApiMessage_Ping
	//	*ApiMessage_Pong
	//	*ApiMessage_HandshakeMessage
	//	*ApiMessage_UnreliableMessage
	//	*ApiMessage_ReliableMessage
	//	*ApiMessage_Mytest
	NoticeWay  isApiMessage_NoticeWay `protobuf_oneof:"notice_way"`
	Namename22 string                 `protobuf:"bytes,7,opt,name=namename22,proto3" json:"namename22,omitempty"`
}

func (x *ApiMessage) Reset() {
	*x = ApiMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiMessage) ProtoMessage() {}

func (x *ApiMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiMessage.ProtoReflect.Descriptor instead.
func (*ApiMessage) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{5}
}

func (m *ApiMessage) GetNoticeWay() isApiMessage_NoticeWay {
	if m != nil {
		return m.NoticeWay
	}
	return nil
}

func (x *ApiMessage) GetPing() *Ping {
	if x, ok := x.GetNoticeWay().(*ApiMessage_Ping); ok {
		return x.Ping
	}
	return nil
}

func (x *ApiMessage) GetPong() *Pong {
	if x, ok := x.GetNoticeWay().(*ApiMessage_Pong); ok {
		return x.Pong
	}
	return nil
}

func (x *ApiMessage) GetHandshakeMessage() *HandshakeMessage {
	if x, ok := x.GetNoticeWay().(*ApiMessage_HandshakeMessage); ok {
		return x.HandshakeMessage
	}
	return nil
}

func (x *ApiMessage) GetUnreliableMessage() *UnreliableMessage {
	if x, ok := x.GetNoticeWay().(*ApiMessage_UnreliableMessage); ok {
		return x.UnreliableMessage
	}
	return nil
}

func (x *ApiMessage) GetReliableMessage() *ReliableMessage {
	if x, ok := x.GetNoticeWay().(*ApiMessage_ReliableMessage); ok {
		return x.ReliableMessage
	}
	return nil
}

func (x *ApiMessage) GetMytest() *Mytest {
	if x, ok := x.GetNoticeWay().(*ApiMessage_Mytest); ok {
		return x.Mytest
	}
	return nil
}

func (x *ApiMessage) GetNamename22() string {
	if x != nil {
		return x.Namename22
	}
	return ""
}

type isApiMessage_NoticeWay interface {
	isApiMessage_NoticeWay()
}

type ApiMessage_Ping struct {
	Ping *Ping `protobuf:"bytes,1,opt,name=ping,proto3,oneof"`
}

type ApiMessage_Pong struct {
	Pong *Pong `protobuf:"bytes,2,opt,name=pong,proto3,oneof"`
}

type ApiMessage_HandshakeMessage struct {
	HandshakeMessage *HandshakeMessage `protobuf:"bytes,3,opt,name=handshakeMessage,proto3,oneof"`
}

type ApiMessage_UnreliableMessage struct {
	UnreliableMessage *UnreliableMessage `protobuf:"bytes,4,opt,name=unreliableMessage,proto3,oneof"`
}

type ApiMessage_ReliableMessage struct {
	ReliableMessage *ReliableMessage `protobuf:"bytes,5,opt,name=reliableMessage,proto3,oneof"`
}

type ApiMessage_Mytest struct {
	Mytest *Mytest `protobuf:"bytes,6,opt,name=mytest,proto3,oneof"`
}

func (*ApiMessage_Ping) isApiMessage_NoticeWay() {}

func (*ApiMessage_Pong) isApiMessage_NoticeWay() {}

func (*ApiMessage_HandshakeMessage) isApiMessage_NoticeWay() {}

func (*ApiMessage_UnreliableMessage) isApiMessage_NoticeWay() {}

func (*ApiMessage_ReliableMessage) isApiMessage_NoticeWay() {}

func (*ApiMessage_Mytest) isApiMessage_NoticeWay() {}

type Mytest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Mytest) Reset() {
	*x = Mytest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Mytest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Mytest) ProtoMessage() {}

func (x *Mytest) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Mytest.ProtoReflect.Descriptor instead.
func (*Mytest) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{6}
}

func (x *Mytest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SentAt int64 `protobuf:"varint,1,opt,name=sent_at,json=sentAt,proto3" json:"sent_at,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{7}
}

func (x *Ping) GetSentAt() int64 {
	if x != nil {
		return x.SentAt
	}
	return 0
}

type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PingSentAt int64 `protobuf:"varint,1,opt,name=ping_sent_at,json=pingSentAt,proto3" json:"ping_sent_at,omitempty"`
	ReceivedAt int64 `protobuf:"varint,2,opt,name=received_at,json=receivedAt,proto3" json:"received_at,omitempty"`
	SentAt     int64 `protobuf:"varint,3,opt,name=sent_at,json=sentAt,proto3" json:"sent_at,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{8}
}

func (x *Pong) GetPingSentAt() int64 {
	if x != nil {
		return x.PingSentAt
	}
	return 0
}

func (x *Pong) GetReceivedAt() int64 {
	if x != nil {
		return x.ReceivedAt
	}
	return 0
}

func (x *Pong) GetSentAt() int64 {
	if x != nil {
		return x.SentAt
	}
	return 0
}

type ApiChat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Msg    string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	// Types that are assignable to NoticeWay:
	//
	//	*ApiChat_Email
	//	*ApiChat_Phone
	NoticeWay isApiChat_NoticeWay `protobuf_oneof:"notice_way"`
}

func (x *ApiChat) Reset() {
	*x = ApiChat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gudp_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiChat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiChat) ProtoMessage() {}

func (x *ApiChat) ProtoReflect() protoreflect.Message {
	mi := &file_gudp_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiChat.ProtoReflect.Descriptor instead.
func (*ApiChat) Descriptor() ([]byte, []int) {
	return file_gudp_proto_rawDescGZIP(), []int{9}
}

func (x *ApiChat) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ApiChat) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (m *ApiChat) GetNoticeWay() isApiChat_NoticeWay {
	if m != nil {
		return m.NoticeWay
	}
	return nil
}

func (x *ApiChat) GetEmail() string {
	if x, ok := x.GetNoticeWay().(*ApiChat_Email); ok {
		return x.Email
	}
	return ""
}

func (x *ApiChat) GetPhone() string {
	if x, ok := x.GetNoticeWay().(*ApiChat_Phone); ok {
		return x.Phone
	}
	return ""
}

type isApiChat_NoticeWay interface {
	isApiChat_NoticeWay()
}

type ApiChat_Email struct {
	Email string `protobuf:"bytes,3,opt,name=email,proto3,oneof"`
}

type ApiChat_Phone struct {
	Phone string `protobuf:"bytes,4,opt,name=phone,proto3,oneof"`
}

func (*ApiChat_Email) isApiChat_NoticeWay() {}

func (*ApiChat_Phone) isApiChat_NoticeWay() {}

var File_gudp_proto protoreflect.FileDescriptor

var file_gudp_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x75, 0x64, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x01, 0x0a,
	0x0d, 0x61, 0x70, 0x69, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x12, 0x1b,
	0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x63, 0x65, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x65, 0x6e,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x72, 0x65, 0x65, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x72, 0x65, 0x65, 0x72, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65,
	0x6e, 0x63, 0x65, 0x22, 0xc7, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x73, 0x65,
	0x71, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03,
	0x61, 0x63, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x6b, 0x42, 0x69, 0x74, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x63, 0x6b, 0x42, 0x69, 0x74, 0x73, 0x22, 0xcd, 0x01,
	0x0a, 0x10, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6f, 0x6b, 0x69, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x63, 0x6f, 0x6f, 0x6b, 0x69, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x61, 0x6e, 0x64,
	0x6f, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x22, 0x46, 0x0a,
	0x11, 0x55, 0x6e, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x50, 0x0a, 0x0b, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x70, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xdb, 0x02, 0x0a, 0x0b, 0x61, 0x70, 0x69, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x04,
	0x70, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x04, 0x70, 0x6f, 0x6e,
	0x67, 0x12, 0x3f, 0x0a, 0x10, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x48, 0x61,
	0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00,
	0x52, 0x10, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x42, 0x0a, 0x11, 0x75, 0x6e, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x55, 0x6e, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x48, 0x00, 0x52, 0x11, 0x75, 0x6e, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3c, 0x0a, 0x0f, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x52, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x48, 0x00, 0x52, 0x0f, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x06, 0x6d, 0x79, 0x74, 0x65, 0x73, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x79, 0x74, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52,
	0x06, 0x6d, 0x79, 0x74, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x32, 0x32, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x61, 0x6d,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x32, 0x42, 0x0c, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x5f, 0x77, 0x61, 0x79, 0x22, 0x1c, 0x0a, 0x06, 0x4d, 0x79, 0x74, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0x1f, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x73,
	0x65, 0x6e, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x65,
	0x6e, 0x74, 0x41, 0x74, 0x22, 0x62, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x20, 0x0a, 0x0c,
	0x70, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x73, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x73, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x22, 0x73, 0x0a, 0x08, 0x61, 0x70, 0x69, 0x5f,
	0x63, 0x68, 0x61, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12,
	0x16, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x42,
	0x0c, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x5f, 0x77, 0x61, 0x79, 0x2a, 0x88, 0x01,
	0x0a, 0x11, 0x67, 0x75, 0x64, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x47, 0x55, 0x44, 0x50, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04,
	0x50, 0x4f, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x48, 0x41, 0x4e, 0x44, 0x53, 0x48,
	0x41, 0x4b, 0x45, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11,
	0x55, 0x4e, 0x52, 0x45, 0x4c, 0x49, 0x41, 0x42, 0x4c, 0x45, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x10, 0x04, 0x12, 0x13, 0x0a, 0x0f, 0x52, 0x45, 0x4c, 0x49, 0x41, 0x42, 0x4c, 0x45, 0x4d,
	0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x50, 0x43, 0x4d,
	0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x06, 0x42, 0x12, 0x5a, 0x10, 0x67, 0x75, 0x64, 0x70,
	0x2f, 0x67, 0x75, 0x64, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gudp_proto_rawDescOnce sync.Once
	file_gudp_proto_rawDescData = file_gudp_proto_rawDesc
)

func file_gudp_proto_rawDescGZIP() []byte {
	file_gudp_proto_rawDescOnce.Do(func() {
		file_gudp_proto_rawDescData = protoimpl.X.CompressGZIP(file_gudp_proto_rawDescData)
	})
	return file_gudp_proto_rawDescData
}

var file_gudp_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gudp_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_gudp_proto_goTypes = []interface{}{
	(GudpMessageType)(0),      // 0: gudp_message_type
	(*ApiCharacter)(nil),      // 1: api_character
	(*ReliableMessage)(nil),   // 2: ReliableMessage
	(*HandshakeMessage)(nil),  // 3: HandshakeMessage
	(*UnreliableMessage)(nil), // 4: UnreliableMessage
	(*RpcMessage)(nil),        // 5: rpc_message
	(*ApiMessage)(nil),        // 6: api_message
	(*Mytest)(nil),            // 7: Mytest
	(*Ping)(nil),              // 8: Ping
	(*Pong)(nil),              // 9: Pong
	(*ApiChat)(nil),           // 10: api_chat
}
var file_gudp_proto_depIdxs = []int32{
	8, // 0: api_message.ping:type_name -> Ping
	9, // 1: api_message.pong:type_name -> Pong
	3, // 2: api_message.handshakeMessage:type_name -> HandshakeMessage
	4, // 3: api_message.unreliableMessage:type_name -> UnreliableMessage
	2, // 4: api_message.reliableMessage:type_name -> ReliableMessage
	7, // 5: api_message.mytest:type_name -> Mytest
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_gudp_proto_init() }
func file_gudp_proto_init() {
	if File_gudp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gudp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiCharacter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReliableMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnreliableMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Mytest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pong); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gudp_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiChat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_gudp_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*ApiMessage_Ping)(nil),
		(*ApiMessage_Pong)(nil),
		(*ApiMessage_HandshakeMessage)(nil),
		(*ApiMessage_UnreliableMessage)(nil),
		(*ApiMessage_ReliableMessage)(nil),
		(*ApiMessage_Mytest)(nil),
	}
	file_gudp_proto_msgTypes[9].OneofWrappers = []interface{}{
		(*ApiChat_Email)(nil),
		(*ApiChat_Phone)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gudp_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gudp_proto_goTypes,
		DependencyIndexes: file_gudp_proto_depIdxs,
		EnumInfos:         file_gudp_proto_enumTypes,
		MessageInfos:      file_gudp_proto_msgTypes,
	}.Build()
	File_gudp_proto = out.File
	file_gudp_proto_rawDesc = nil
	file_gudp_proto_goTypes = nil
	file_gudp_proto_depIdxs = nil
}
