package main

import (
	"net"
	"net/rpc"
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

	// 2.注册服务,通过结构体绑定解决callID的问题
	_ = rpc.RegisterName("HelloService", &HelloService{})

	// 3. 监听端口
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)
	}
}
