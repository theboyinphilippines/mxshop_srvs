package main

import "fmt"

/*
设计模式之函数选项模式
*/
//创建一个数据库结构体
type dbConfig struct {
	Host     string
	Port     int
	Name     string
	Password string
	DbName   string
}

// 自定义一个option的函数类型，接收的参数是dbConfig
type Option func(*dbConfig)

// 这个函数用来设置host
func WithHost(host string) Option {
	return func(db *dbConfig) {
		db.Host = host
	}
}

// 这个函数用来设置port
func WithPort(port int) Option {
	return func(db *dbConfig) {
		db.Port = port
	}
}

// new db client函数，可以传入很多选项
func NewDBClient(options ...Option) *dbConfig {
	//设置默认值
	db := &dbConfig{
		Host: "127.0.0.1",
		Port: 3306,
	}
	for _, option := range options {
		//调用option函数，并传入参数dbConfig的实例
		option(db)
	}
	return db
}
func main() {
	//创建一个db client，并且传入要自定义设置的值
	dbClient := NewDBClient(WithHost("192.168.0.101"), WithPort(8856))
	fmt.Println(dbClient)
	//不传递，使用默认值
	dbClient2 := NewDBClient()
	fmt.Println(dbClient2)

}
