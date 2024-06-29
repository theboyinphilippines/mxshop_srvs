package main

import (
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/handler"
	"mxshop_srvs/goods_srv/initialize"
	"mxshop_srvs/goods_srv/proto"
	"mxshop_srvs/goods_srv/utils"
	"mxshop_srvs/goods_srv/utils/otgrpc"
	"mxshop_srvs/goods_srv/utils/register/consul"
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
	initialize.InitEs()
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
		ServiceName: "mxshop-goods-srv",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(any(err))
	}
	opentracing.SetGlobalTracer(tracer)

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		//panic(any("fail to listen:" + err.Error()))
		panic(any("fail to listen:" + err.Error()))
	}
	// 将grpc服务 注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	//cfg := api.DefaultConfig()
	//cfg.Address = fmt.Sprintf("%s:%d",
	//	global.ServerConfig.ConsulInfo.Host,
	//	global.ServerConfig.ConsulInfo.Port)
	//client, err := api.NewClient(cfg)
	//if err != nil {
	//	panic(any(err))
	//}
	//
	//// 生成对应的检查对象
	//check := &api.AgentServiceCheck{
	//	GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
	//	Timeout:                        "5s",
	//	Interval:                       "5s",
	//	DeregisterCriticalServiceAfter: "10s",
	//}
	//
	//registration := new(api.AgentServiceRegistration)
	//registration.Name = global.ServerConfig.Name
	//serviceID := fmt.Sprintf("%s", uuid.NewV4())
	//registration.ID = serviceID
	//registration.Port = *Port
	//registration.Tags = global.ServerConfig.Tags
	//registration.Address = global.ServerConfig.Host
	//registration.Check = check
	//
	//// 生成注册对象
	//err = client.Agent().ServiceRegister(registration)
	//if err != nil {
	//	panic(any(err))
	//}
	//
	//go func() {
	//	err = server.Serve(lis)
	//	if err != nil {
	//		panic(any("fail to listen:" + err.Error()))
	//	}
	//}()
	//
	////接收终止信号， 优雅关闭服务
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//if err = client.Agent().ServiceDeregister(serviceID); err != nil {
	//	zap.S().Info("注销失败")
	//}
	//zap.S().Info("注销成功")

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

	//接收终止信号， 优雅关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_ = closer.Close()
	//注销服务
	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("注销失败：", err.Error())
	} else {
		zap.S().Info("注销成功")
	}

}
