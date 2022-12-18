package handler

import (
	"context"
	"strconv"
	"time"

	"rpc_grpc_learn/rpc/04_gRPC/proto"
)

// Server 实现生成的Greeter接口才能被注册,必须要嵌入一个匿名结构体
type Server struct {
	proto.UnimplementedGreeterServer // 必须要嵌入该结构体才能注册
}

// SayHello 就是具体的业务处理逻辑
func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) { // ctx和error是强制生成的
	return &proto.HelloReply{Message: "hello " + request.Name}, nil
}

// StreamServer 流模式的实现
type StreamServer struct {
	proto.UnimplementedStreamGreeterServer
}

func (s *StreamServer) GetStream(req *proto.StreamReqData, res proto.StreamGreeter_GetStreamServer) error {
	// 服务端流模式
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

EXIT:
	for {
		select {
		case <-ticker.C:
			if err := res.Send(&proto.StreamResData{
				Data: strconv.FormatInt(time.Now().UnixMilli(), 10),
			}); err != nil {
				break EXIT
			}
		case <-ctx.Done():
			if err := res.Send(&proto.StreamResData{
				Data: ctx.Err().Error(),
			}); err != nil {
				break EXIT
			}
			break EXIT
		}
	}

	return nil
}

func (s *StreamServer) PostStream(req proto.StreamGreeter_PostStreamServer) error {
	panic("not implemented") // TODO: Implement
}

func (s *StreamServer) AllStream(req proto.StreamGreeter_AllStreamServer) error {
	panic("not implemented") // TODO: Implement
}
