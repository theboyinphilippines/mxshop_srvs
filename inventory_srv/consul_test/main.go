package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// 这里是测试consul连接，与项目无关
func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.101:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(any(err))

	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.1.101:8021/health",
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
		panic(any(err))
	}
	return nil

}

func Allservice() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.101:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(any(err))
	}
	// 获取所有服务
	data, err := client.Agent().Services()
	if err != nil {
		panic(any(err))
	}
	for key, _ := range data {
		fmt.Println(key)

	}

}

func FilterService() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.101:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(any(err))
	}
	data, err := client.Agent().ServicesWithFilter(`Service==inventory_srv`)
	if err != nil {
		return
	}
	for key, _ := range data {
		fmt.Println(key)
	}

}

func main() {
	_ = Register("192.168.1.101", 8021,
		"user-web", []string{"mxshop-web", "shy"},
		"user-web")

}
