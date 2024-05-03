package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// 基本模型字段
type BaseModel struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"` //为什么使用int32， bigint
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	//gorm中添加该字段为软删除，数据库中需要加上该字段deleted_at
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool `json:"-"`
}

// gorm中实现自定义数据类型
type GormList []string

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}
// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte),&g)
}