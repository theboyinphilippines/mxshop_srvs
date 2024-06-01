package global

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"mxshop_srvs/goods_srv/config"
)

// init方法：被引用这个包时，会自动调用init方法
var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	EsClient     *elastic.Client
)
