package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/handler"
	"mxshop_srvs/inventory_srv/initialize"
	"mxshop_srvs/inventory_srv/proto"
	"mxshop_srvs/inventory_srv/utils"
	"mxshop_srvs/inventory_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 命令行输入ip port来启动
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50053, "端口号")

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
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(any("fail to listen:" + err.Error()))
	}

	// 将grpc服务 注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	//服务注册到consul中
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = registerClient.Register(global.ServerConfig.Host,
		*Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败：", err.Error())
	}
	zap.S().Debugf("启动服务器，端口：%d", *Port)

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(any("fail to listen:" + err.Error()))
		}
	}()

	//消费[库存归还]消息order_reback
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.0.101:9876"}), consumer.WithGroupName("inventory"))
	if err != nil {
		zap.S().Errorf("消费者初始化失败：%v", err)
	}
	err = c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback)
	if err != nil {
		zap.S().Errorf("消费者订阅消息失败：%v", err)
	}
	err = c.Start()
	if err != nil {
		zap.S().Errorf("消费者开始失败：%v", err)
	}

	//接收终止信号， 优雅关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("注销失败：", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
