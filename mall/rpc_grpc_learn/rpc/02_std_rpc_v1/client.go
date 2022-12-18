package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 1.建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}

	// 2.调用哪一个函数
	var reply string
	err = client.Call("HelloService.Hello", "leilei", &reply)
	if err != nil {
		panic(err)
	}
	/*
		调用函数大部分都是net包的内容过于冗余,使用起来不方便,能不能再封装,使得rpc调用就像调用本地函数一样?
		我想这样调用:
		client.Hello()
	*/

	fmt.Println(reply)
}

/*
	标准库的rpc,客户端和服务端可以将数据协议换为json,也可以把网络协议更换为http协议

*/
