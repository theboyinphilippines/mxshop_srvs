package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/order_srv/config"
	"mxshop_srvs/order_srv/proto"
)

// init方法：被引用这个包时，会自动调用init方法
var (
	DB                 *gorm.DB
	ServerConfig       config.ServerConfig
	NacosConfig        config.NacosConfig
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
