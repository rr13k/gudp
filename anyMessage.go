package main

import (
	"github.com/rr13k/gudp/protos"

	"google.golang.org/protobuf/runtime/protoiface"
)

type AnyMessage struct {
	protoiface.MessageV1
	ReliableMessage   *protos.ReliableMessage
	UnreliableMessage *protos.UnreliableMessage
}
