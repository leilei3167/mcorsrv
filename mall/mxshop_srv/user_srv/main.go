package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"user_srv/user_srv/handler"
	"user_srv/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("p", 50051, "端口")
	flag.Parse()
	fmt.Printf("serve on: %s:%d", *IP, *Port)

	s := grpc.NewServer()
	proto.RegisterUserServer(s, &handler.UserServer{})
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}
	panic(s.Serve(listener))
}
