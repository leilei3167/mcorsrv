package server_stub

import (
	"net/rpc"

	proto "rpc_grpc_learn/rpc/03_std_rpc_v2/handler"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}

// 为了避免耦合,使用接口作为参数
func RegisterHelloService(srv HelloServicer) error {
	return rpc.RegisterName(proto.HelloService, srv)
}
