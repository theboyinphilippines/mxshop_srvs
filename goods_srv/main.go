package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/handler"
	"mxshop_srvs/goods_srv/initialize"
	"mxshop_srvs/goods_srv/proto"
	"mxshop_srvs/goods_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 命令行输入ip port来启动
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50052, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info("全局配置是:", global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip:", *IP)
	//获取空闲的端口号
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port:", *Port)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("fail to listen:" + err.Error())
	}
	// 将grpc服务 注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	// 生成注册对象
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("fail to listen:" + err.Error())
		}
	}()

	//接收终止信号， 优雅关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

}
