package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/opentracing/opentracing-go"
	"github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/handler"
	"mxshop_srvs/order_srv/initialize"
	"mxshop_srvs/order_srv/proto"
	"mxshop_srvs/order_srv/utils"
	"mxshop_srvs/order_srv/utils/otgrpc"
	"mxshop_srvs/order_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 命令行输入ip port来启动
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	// 初始化第三方微服务连接
	initialize.InitSrvConn()
	zap.S().Info("全局配置是:", global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip:", *IP)
	//获取空闲的端口号
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port:", *Port)

	//链路追踪配置
	cfg := jaegercfg.Configuration{
		//采样器设置
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		//jaeger agent设置
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.0.101:6831",
		},
		ServiceName: "mxshop-order-srv",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(any(err))
	}
	opentracing.SetGlobalTracer(tracer)

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterOrderServer(server, &handler.OrderServer{})
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

	//消费库存归还消息order_timeout
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.0.101:9876"}), consumer.WithGroupName("order"))
	if err != nil {
		zap.S().Errorf("消费者初始化失败：%v", err)
	}
	err = c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeOut)
	if err != nil {
		zap.S().Errorf("消费者订阅消息失败：%v", err)
	}
	_ = c.Start()

	//接收终止信号， 优雅关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = closer.Close()
	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("注销失败：", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
