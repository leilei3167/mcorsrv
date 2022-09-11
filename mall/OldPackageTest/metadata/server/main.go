package main

import (
	"OldPackageTest/metadata/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
)

type serverHandler struct {
	proto.UnimplementedGreeterServer
}

func (s *serverHandler) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	//从ctx中获取全部元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("没有携带元数据")
		return &proto.HelloReply{Msg: "请附加元数据再发起请求"}, nil
	}

	return &proto.HelloReply{Msg: request.Name + fmt.Sprintf("%#v", md)}, nil

}

func main() {
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &serverHandler{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	panic(s.Serve(lis))

}
