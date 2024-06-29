package main

import (
	"mxshop_srvs/logx/self_define"
)

func main() {
	//初始化日志
	self_define.Init(self_define.NewOptions())
	//使用日志
	self_define.GetLogger().Debug("debug info")
	self_define.Debug("debug info")

}
