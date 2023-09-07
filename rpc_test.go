package gudp

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/rr13k/gudp/gudp_protos"
)

type UdpRpc struct {
}

// 算数运算请求结构体
type UdpRpcRequest struct {
	A int
	B int
}

// 算数运算响应结构体
type UdpRpcResponse struct {
	Pro int // 乘积
	Quo int // 商
	Rem int // 余数
}

func (u *UdpRpc) Multiply(req UdpRpcRequest, res *UdpRpcResponse) error {
	res.Pro = req.A * req.B
	return nil
}

func TestXxx(t *testing.T) {

	// mock 请求类型
	var udpRpcRequest = UdpRpcRequest{
		A: 1,
		B: 2,
	}
	udpRpcRequestByte, _ := json.Marshal(udpRpcRequest)

	rpc, _ := NewRpc(new(UdpRpc))

	var msg = gudp_protos.RpcMessage{
		Method: "Multiply",
		Data:   udpRpcRequestByte,
	}

	var mtype = rpc.method[msg.GetMethod()]

	var ctx = context.Background()

	var data = msg.GetData()
	var argv reflect.Value
	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if mtype.ArgType.Kind() == reflect.Pointer {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	json.Unmarshal(data, argv.Interface())
	if argIsValue {
		argv = argv.Elem()
	}
	//解析转换完成请求参数
	// 传递响应内容
	var replyv = reflect.New(mtype.ReplyType.Elem())
	switch mtype.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
	}

	rpc.Call(ctx, mtype, argv, replyv)
}
