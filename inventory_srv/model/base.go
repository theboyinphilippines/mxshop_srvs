package model

import (
	"gorm.io/gorm"
	"time"
)

// 基本模型字段
type BaseModel struct {
	ID        int32     `gorm:"primarykey;type:int" json:"id"` //为什么使用int32， bigint
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	//gorm中添加该字段为软删除，数据库中需要加上该字段deleted_at
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}
