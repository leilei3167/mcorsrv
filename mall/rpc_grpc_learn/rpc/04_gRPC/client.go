package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"rpc_grpc_learn/rpc/04_gRPC/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接服务
	conn, err := grpc.Dial("127.0.0.1:1234",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 实例化client
	c := proto.NewGreeterClient(conn)
	// 调用SayHello方法
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "leilei"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)

	// 服务端流模式使用
	c1 := proto.NewStreamGreeterClient(conn)
	n, err := c1.GetStream(context.Background(), &proto.StreamReqData{Data: "leilei"})
	if err != nil {
		panic(err)
	}

	for {
		d, err := n.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				panic(err)
			} else {
				fmt.Println("server send done!")
				break
			}
		}

		fmt.Println(d)
	}
}
