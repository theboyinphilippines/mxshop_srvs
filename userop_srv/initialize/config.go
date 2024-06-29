package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_srvs/userop_srv/global"
	"os"
)

// 配置环境变量，根据环境变量来决定用开发还是生产的配置文件
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)

}

func InitConfig() {
	// 从配置文件中读取对应的配置

	debug := GetEnvInfo("MZSHOP_DEBUG")
	fmt.Printf("debug是：%v\n", debug)
	configFilePrefix := "config"
	fmt.Println(os.Getwd())
	configFileName := fmt.Sprintf("userop_srv/%s-pro.yaml", configFilePrefix)
	//configFileName := fmt.Sprintf("./%s-pro.yaml", configFilePrefix)
	if debug {
		//configFileName = fmt.Sprintf("./%s-debug.yaml", configFilePrefix)
		configFileName = fmt.Sprintf("userop_srv/%s-debug.yaml", configFilePrefix)
	}
	zap.S().Infof("配置文件路径为：%v", configFileName)

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(any(err))
	}

	// serverConfig对象，其他文件中也要使用配置，所以声明为全局变量
	//serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(any(err))
	}
	zap.S().Infof("配置信息：%v", global.NacosConfig)
	fmt.Printf("服务名称是：%v", v.Get("name"))

	// 动态监控配置文件变化 （nacos中已经动态监控）
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	zap.S().Infof("配置文件产生变化：%v", e.Name)
	//	_ = v.ReadInConfig()
	//	_ = v.Unmarshal(global.ServerConfig)
	//	zap.S().Infof("配置信息：%v", global.ServerConfig)
	//})

	// 从nacos中读取配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(any(err))
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(any(err))
	}

	//将从nacos中获取的配置数据绑定到结构体中
	fmt.Println("这是content", content)
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	fmt.Println("这是global.ServerConfig", &global.ServerConfig)

}
