package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// 设计服务注册中心接口（go风格-鸭子模式）

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

type Registry struct {
	Host string
	Port int
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (reg *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", reg.Host, reg.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)

	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	// 生成注册对象
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil

}

func (reg *Registry) DeRegister(serviceId string) error {
	return nil
}
