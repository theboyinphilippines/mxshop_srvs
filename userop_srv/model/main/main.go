package main

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"mxshop_srvs/userop_srv/model"
	"os"
	"time"
)

// md5加密
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))

}

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	dsn := "root:1234qwer!@tcp(127.0.0.1:3306)/mxshop_userop_srv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		//表名为单数
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic(err)
	}

	// 设置全局logger,打印sql语句

	//迁移schema
	_ = db.AutoMigrate(&model.LeavingMessages{}, &model.Address{}, &model.UserFav{})

	//fmt.Println(genMd5("123456"))

}
