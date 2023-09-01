package pb

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	errInvalidProtobufMessage = errors.New("invalid protobuf message")
)

type Protobuf struct{}

func (p *Protobuf) Marshal(msg interface{}) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, errInvalidProtobufMessage
	}
	return proto.Marshal(m)
}

func (p *Protobuf) Unmarshal(raw []byte, msg interface{}) error {
	m, ok := msg.(proto.Message)
	if !ok {
		return errInvalidProtobufMessage
	}
	return proto.Unmarshal(raw, m)
}
