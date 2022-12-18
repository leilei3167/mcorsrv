package client_stub

import (
	"net/rpc"

	proto "rpc_grpc_learn/rpc/03_std_rpc_v2/handler"
)

type HelloServiceStub struct {
	*rpc.Client
}

func NewHelloServiceClient(protocol, add string) *HelloServiceStub {
	client, err := rpc.Dial(protocol, add)
	if err != nil {
		panic(err)
	}
	return &HelloServiceStub{client}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {
	err := c.Call(proto.HelloService+".Hello", request, reply)
	if err != nil {
		return err
	}
	return nil
}
