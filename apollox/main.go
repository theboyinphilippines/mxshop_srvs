package main

import (
	"fmt"
	"github.com/philchia/agollo/v4"
	"log"
	"time"
)

func main() {
	agollo.Start(&agollo.Conf{
		AppID:           "1005784",
		Cluster:         "default",
		NameSpaceNames:  []string{"application.properties", "TEST1.Srv.Namespace.properties"},
		MetaAddr:        "http://localhost:8080",
		AccesskeySecret: "6dbd7a532b9847298d4955639b37aad6",
	})

	agollo.OnUpdate(func(event *agollo.ChangeEvent) {
		// 监听配置变更
		log.Printf("event:%#v\n", event)
	})
	log.Println("初始化Apollo配置成功")

	for {
		// 从默认的application.properties命名空间获取key的值
		val := agollo.GetString("mysql.host")
		log.Printf("mysql.host的值是：%v", val)
		// 获取命名空间下所有key
		keys := agollo.GetAllKeys(agollo.WithNamespace("application.properties"))
		fmt.Println(keys)
		// 获取指定一个命令空间下key的值
		other := agollo.GetString("mysql.port", agollo.WithNamespace("application.properties"))
		log.Println(other)
		batchSize := agollo.GetString("sender.batchSize", agollo.WithNamespace("TEST1.Srv.Namespace.properties"))
		log.Printf("batchSize是: %v\n", batchSize)
		// 获取指定命名空间下的所有内容
		namespaceContent := agollo.GetContent(agollo.WithNamespace("application.properties"))
		log.Println(namespaceContent)
		time.Sleep(time.Second * 3)
	}

}
