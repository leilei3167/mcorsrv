package model

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 用于替换gorm.Model,更灵活
type BaseModel struct {
	ID        int32          `gorm:"primary_key"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt //为了继承软删除
	IsDeleted bool           `gorm:"column:is_deleted"`
}

type User struct {
	BaseModel
	//以手机号创建索引,唯一索引,varchar类型11位,不能为空
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	//密码不是明文,所以要长一些
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"` //nickname可选
	//指针类型的Time,为了能够避免一些错误
	Birthday *time.Time
	Gender   string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male为男'"`
	Role     int    `gorm:"column:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}
