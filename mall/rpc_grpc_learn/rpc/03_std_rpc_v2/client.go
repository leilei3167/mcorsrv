package main

import (
	"fmt"

	"rpc_grpc_learn/rpc/03_std_rpc_v2/client_stub"
)

func main() {
	/* 	client, err := rpc.Dial("tcp", "localhost:1234")
	   	if err != nil {
	   		panic(err)
	   	}

	   	var reply string
	   	err = client.Call(proto.HelloService+".Hello", "leilei", &reply)
	   	if err != nil {
	   		panic(err)
	   	}

	   	fmt.Println(reply)
	*/

	// 我只想用client.Hello()这样使用,怎么办?
	// 封将rpc连接等逻辑单独封装至client_stub中
	client := client_stub.NewHelloServiceClient("tcp", "localhost:1234")
	var reply string
	err := client.Hello("leilei", &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}

/*
此时,client_stub和server_stub的代码其实可以自动生成了,client_stub封装调用方法的过程,server_stub通过接口实现
不同业务逻辑的解耦和注册
业务逻辑可以实现在handler文件夹中
这个就是最基本的,gRPC的实现原理
*/
