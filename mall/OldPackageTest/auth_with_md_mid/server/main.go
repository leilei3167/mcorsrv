package main

import (
	"OldPackageTest/metadata/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

type serverHandler struct {
	proto.UnimplementedGreeterServer
}

func (s *serverHandler) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	//从ctx中获取全部元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("没有携带元数据")
		return &proto.HelloReply{Msg: "请附加元数据再发起请求"}, status.Error(codes.Unauthenticated, "没有元数据")
	}
	tSlice, ok := md["token"]
	if !ok {
		if !ok {
			fmt.Println("没有auth元数据")
			return &proto.HelloReply{Msg: "请附加元数据再发起请求"}, status.Error(codes.Unauthenticated, "没有认证信息")
		}
	}

	return &proto.HelloReply{Msg: request.Name + fmt.Sprintf("%#v", tSlice[0])}, nil

}

func main() {
	//定义一个拦截器要处理的逻辑
	i := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到一个新请求")
		start := time.Now()
		res, err := handler(ctx, req)
		fmt.Println("请求完成,耗时:", time.Since(start))
		return res, err
	}

	//如何将拦截器加入到server中?
	optWithI := grpc.UnaryInterceptor(i)
	s := grpc.NewServer(optWithI) //注册到server中
	proto.RegisterGreeterServer(s, &serverHandler{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	panic(s.Serve(lis))
}
