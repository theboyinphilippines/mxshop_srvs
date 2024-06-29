package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"os"
	"strconv"
	"time"
)

// md5加密
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))

}

func main() {
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold:             time.Second, // Slow SQL threshold
	//		LogLevel:                  logger.Info, // Log level
	//		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	//		ParameterizedQueries:      true,        // Don't include params in the SQL log
	//		Colorful:                  true,        // Disable color
	//	},
	//)
	//
	//dsn := "root:1234qwer!@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: newLogger,
	//	//表名为单数
	//	NamingStrategy: schema.NamingStrategy{SingularTable: true},
	//})
	//if err != nil {
	//	panic(any(err))
	//}
	//
	//// 设置全局logger,打印sql语句
	//
	////迁移schema
	//_ = db.AutoMigrate(&model.Category{}, model.Goods{},
	//	model.Goods{}, model.GoodsCategoryBrand{}, model.Banner{})

	mysqlToEs()

}

// 将mysql中的goods表中的数据同步到es中
func mysqlToEs() {
	dsn := "root:1234qwer!@tcp(192.168.0.101:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(any(err))
	}

	host := "http://192.168.0.101:9200"
	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(any(err))
	}

	var goods []model.Goods
	db.Find(&goods)
	for _, good := range goods {
		gooddata := model.EsGoods{
			ID:          good.ID,
			CategoryID:  good.CategoryID,
			BrandsID:    good.BrandsID,
			OnSale:      good.OnSale,
			ShipFree:    good.ShipFree,
			IsNew:       good.IsNew,
			IsHot:       good.IsHot,
			Name:        good.Name,
			ClickNum:    good.ClickNum,
			SoldNum:     good.SoldNum,
			FavNum:      good.FavNum,
			MarketPrice: good.MarketPrice,
			GoodsBrief:  good.GoodsBrief,
			ShopPrice:   good.ShopPrice,
		}
		_, err := global.EsClient.Index().Index("goods").Id(strconv.Itoa(int(good.ID))).BodyJson(gooddata).Do(context.Background())
		if err != nil {
			return
		}

	}

}
