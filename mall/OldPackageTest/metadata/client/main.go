package main

import (
	"OldPackageTest/metadata/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

func main() {
	//先建立连接
	conn, err := grpc.Dial("127.0.0.1:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//初始化client
	client := proto.NewGreeterClient(conn)
	//设置元数据,写入ctx
	md := metadata.New(map[string]string{
		"name": "first_metaData",
		"time": time.Now().String(),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//调用
	res, err := client.SayHello(ctx, &proto.HelloRequest{Name: "雷磊"})
	if err != nil {
		return
	}
	fmt.Println("res:", res.Msg)
}
