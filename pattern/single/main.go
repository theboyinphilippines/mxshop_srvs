package main

import "sync"

/*
设计模式之单例模式：主要用于实例化全局变量
*/
//实例化一个连接池
type DBPool struct {
}

var dbPool *DBPool
var once sync.Once

func GetDBPool() *DBPool {
	once.Do(func() {
		dbPool = &DBPool{}
	})

	return dbPool
}

func main() {
	GetDBPool()
}
