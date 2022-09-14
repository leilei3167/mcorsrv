package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// BaseModel 用于替换gorm.Model,更灵活
type BaseModel struct {
	//为什么使用int32? 如果一个表的外键指向主键,且他们类型不一致,就会出问题;因此为确保一致性,所以必须要确定类型
	//此处将主键同意设置为int类型
	ID        int32          `gorm:"primary_key;type:int"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt //为了继承软删除
	IsDeleted bool           `gorm:"column:is_deleted"`
}

type GormList []string //实现gorm接口,作为自定义的类型

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(src any) error {
	return json.Unmarshal(src.([]byte), &g)
}