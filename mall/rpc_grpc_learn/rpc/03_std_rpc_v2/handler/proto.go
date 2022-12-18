package proto

// 提供一个公共的包,供客户端和服务端引用
const HelloService = "proto/HelloService"

type NewHelloService struct{}

func (s *NewHelloService) Hello(request string, reply *string) error {
	*reply = "hello " + request
	return nil
}
