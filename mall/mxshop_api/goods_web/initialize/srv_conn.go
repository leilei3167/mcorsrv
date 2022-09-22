package initialize

import (
	"fmt"

	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	dsn := fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=manual", consulInfo.Host, consulInfo.Port,
		global.ServerConfig.GoodsSrvInfo.Name)
	zap.S().Infof("[consul dsn]:%v", dsn)
	conn, err := grpc.Dial(dsn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询模式插件
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn]: 连接用户服务失败:%v", err)
	}
	global.GoodsSrvClient = proto.NewGoodsClient(conn)
}
