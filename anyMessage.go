package gudp

import (
	"github.com/rr13k/gudp/gudp_protos"

	"google.golang.org/protobuf/runtime/protoiface"
)

type AnyMessage struct {
	protoiface.MessageV1
	ReliableMessage   *gudp_protos.ReliableMessage
	UnreliableMessage *gudp_protos.UnreliableMessage
}
