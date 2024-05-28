package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/inventory_srv/config"
)

// init方法：被引用这个包时，会自动调用init方法
var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)
