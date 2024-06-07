package model

import (
	"database/sql/driver"
	"encoding/json"
)

//库存表
type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的mysql乐观锁
}

//订单库存扣减情况表
type SellDetail struct {
	BaseModel
	OrderSn string        `gorm:"type:varchar(30);uniqueIndex"`
	Detail  GormGoodsList `gorm:"type:varchar(200)"`
	Status  int32         `gorm:"type:int"` //1代表已扣减，2代表已归还
}

func (SellDetail) TableName() string {
	return "selldetail"
}

type GoodsNum struct {
	GoodsId int32
	Num     int32
}

// gorm中实现自定义数据类型
type GormGoodsList []GoodsNum

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormGoodsList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormGoodsList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
