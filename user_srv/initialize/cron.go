package initialize

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"time"
)

// 初始化定时任务
func InitCron() {
	//添加时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc))
	//添加任务
	_, _ = c.AddFunc("0 9 * * *", VipGrade) //早上9点执行
	//每秒执行
	//_, _ = c.AddFunc("@every 1s", VipGrade)
	//每小时执行
	//_, _ = c.AddFunc("* */9 * * *", VipGrade)
	c.Start()

	select {}
}

// vip升级定时任务
func VipGrade() {
	//批量更新，vip等级加1， 下面两种都可以，一种更新多列，一种更新单列
	//global.DB.Model(&model.UserVip{}).Where("is_upgrade = 1").Updates(map[string]interface{}{"vip_level": gorm.Expr("vip_level + ?", 1)})
	//global.DB.Model(&model.UserVip{}).Where("is_upgrade = 1").UpdateColumn("vip_level", gorm.Expr("vip_level +?", 1))

	// 事务保证
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()
	err := tx.Model(&model.UserVip{}).Where("is_upgrade = 1").Update("vip_level", gorm.Expr("vip_level +?", 1)).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	zap.S().Info("定时任务9点升级成功")
}
