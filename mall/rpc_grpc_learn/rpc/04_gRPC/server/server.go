package main

import (
	"net"

	"rpc_grpc_learn/rpc/04_gRPC/proto"
	"rpc_grpc_learn/rpc/04_gRPC/server/handler"

	"google.golang.org/grpc"
)

func main() {
	// 初始化server
	g := grpc.NewServer()
	// 注册服务
	proto.RegisterGreeterServer(g, &handler.Server{})
	proto.RegisterStreamGreeterServer(g, &handler.StreamServer{})
	// 启动监听服务
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	err = g.Serve(lis)
	if err != nil {
		panic(err)
	}
}
