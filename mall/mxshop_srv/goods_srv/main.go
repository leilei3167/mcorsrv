package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/handler"
	"mxshop_srv/goods_srv/initialize"
	"mxshop_srv/goods_srv/proto"
	"mxshop_srv/goods_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var host = "172.30.90.215"

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("p", 0, "端口")
	flag.Parse()
	if *Port == 0 { //如果未指定监听端口,则随机选一个
		p, err := utils.GetFreePort()
		if err != nil {
			panic(err)
		}
		*Port = p
	}

	//1.
	initialize.InitLogger()
	//2.
	initialize.InitConfig()
	initialize.InitDB()

	s := grpc.NewServer()
	proto.RegisterUserServer(s, &handler.UserServer{})
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}

	//注册健康检查的服务(默认的接口实现)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	//注册到consul
	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", host, *Port),
		Timeout:                        "3s",
		Interval:                       "3s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	serviceID := uuid.New().String()
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name //注册的服务名字
	registration.ID = serviceID                  //注册到consul时要保证id不同,同名又同id的服务将会被覆盖
	registration.Port = *Port
	registration.Tags = []string{"leilei", "user", "srv"}
	registration.Address = host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("RPC服务开启,addr:%s:%d", host, *Port)

	go func() {
		panic(s.Serve(listener))
	}()
	//优雅退出,收到退出通知后,将自己的服务立刻注销(否则consul会很久之后才会注销)
	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Errorf("注销服务错误:%v", err)
	} else {
		zap.S().Info("注销服务成功")
	}

}
