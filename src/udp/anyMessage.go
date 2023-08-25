package udp

import (
	protoc "pulseHero/src/protos"

	"google.golang.org/protobuf/runtime/protoiface"
)

type AnyMessage struct {
	protoiface.MessageV1
	ReliableMessage   *protoc.ReliableMessage
	UnreliableMessage *protoc.UnreliableMessage
}
