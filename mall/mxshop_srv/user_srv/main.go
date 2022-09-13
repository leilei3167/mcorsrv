package main

import (
	"flag"
	"fmt"
	"net"
	"user_srv/user_srv/global"
	"user_srv/user_srv/handler"
	"user_srv/user_srv/initialize"
	"user_srv/user_srv/proto"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("p", 50051, "端口")
	flag.Parse()
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
		GRPC:                           fmt.Sprintf("172.18.208.68:%d", *Port),
		Timeout:                        "3s",
		Interval:                       "3s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"leilei", "user", "srv"}
	registration.Address = "172.18.208.68"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	panic(s.Serve(listener))
}
