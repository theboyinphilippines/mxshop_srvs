package model

import "time"
import "gorm.io/gorm"

// 基本模型字段
type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	//软删除
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel `json:"base_model"`
	Mobile    string `gorm:"index:idx_mobile;unique;type:varchar(11);not null" json:"mobile,omitempty"`
	Password  string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	NickName  string `gorm:"type:varchar(20)" json:"nick_name,omitempty"`
	// 日期保存容易报错，这里用指针类型
	Birthday *time.Time `gorm:"type:datetime" json:"birthday,omitempty"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'" json:"gender,omitempty"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'" json:"role,omitempty"`
}

type UserVip struct {
	BaseModel `json:"base_model"`
	UserID    int32 `gorm:"type:int;not null" json:"user_id"`
	User      User
	VipLevel  int32 `gorm:"type:int;default:0;not null" json:"vip_level"`
	IsUpgrade bool  `gorm:"type:bool;default:false" json:"is_upgrade"`
}
