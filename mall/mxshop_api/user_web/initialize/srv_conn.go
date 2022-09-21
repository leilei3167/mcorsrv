package initialize

import (
	"fmt"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/proto"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	dsn := fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=manual", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name)
	zap.S().Infof("[consul dsn]:%v", dsn)
	conn, err := grpc.Dial(dsn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //轮询模式插件
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn]: 连接用户服务失败:%v", err)
	}
	global.UserSrvClient = proto.NewUserClient(conn)
}

func InitSrvConn_1() {
	consulInfo := global.ServerConfig.ConsulInfo
	//0.从注册中心获取到用户服务的信息(服务发现,输入服务名,如"user-srv)
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,
		global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	for k, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		zap.S().Debugf("key:%v value:%v", k, v)
		break
	}
	if userSrvHost == "" {
		zap.S().Fatalf("InitSrvConn: 连接用户服务失败")
	}

	//连接用户gRPC服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户服务失败", "msg", err.Error())
		return
	}

	//生成客户端调用gRPC接口,赋值给全局变量
	//TODO:
	//1.这个服务下线了?
	//2.服务地址变更了?
	//3.并发安全?
	//4.此处相当于维护了一个可复用的连接,进一步可扩展为连接池 https://github.com/processout/grpc-go-pool
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
