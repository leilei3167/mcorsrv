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
	//定义拦截器
	i := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Println("调用完成,耗时:", time.Since(start))
		return err
	}
	//生成带拦截器的选项
	optWithI := grpc.WithUnaryInterceptor(i)

	//先建立连接,此处设置客户端过滤器
	conn, err := grpc.Dial("127.0.0.1:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()), optWithI)
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
