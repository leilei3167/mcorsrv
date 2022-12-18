package main

import (
	"net"
	"net/rpc"

	proto "rpc_grpc_learn/rpc/03_std_rpc_v2/handler"
	"rpc_grpc_learn/rpc/03_std_rpc_v2/server_stub"
)

type HelloService struct{}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello " + request
	return nil
}

func main() {
	// 1.初始化Server
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	// 2.使用公共的协议名,注册服务,通过结构体绑定解决callID的问题
	// _ = rpc.RegisterName(proto.HelloService, &HelloService{})

	// 对于server,我也只想关注业务逻辑,而不是他的名字之类的,该如何封装?
	_ = server_stub.RegisterHelloService(&proto.NewHelloService{})

	// 3. 监听端口
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)
	}
}
