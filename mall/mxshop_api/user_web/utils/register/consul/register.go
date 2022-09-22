package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// 将web服务注册到consul,要考虑解耦,因为可能consul会被更换.

type Registry struct {
	Host string
	Port int
}

func NewRegistry(host string, port int) RegistryClient {
	return &Registry{Host: host, Port: port}
}

// 注册中心应该有注册和注销的方法.
type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceID string) error
}

func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port) // consul地址

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	check := &api.AgentServiceCheck{ // 注册健康检查服务
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	// 注册服务
	rgs := &api.AgentServiceRegistration{Name: name, ID: id, Port: port, Tags: tags, Address: address, Check: check}
	if err := client.Agent().ServiceRegister(rgs); err != nil {
		return err
	}
	return nil
}

func (r *Registry) DeRegister(serviceID string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port) // consul地址

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(serviceID)
}
