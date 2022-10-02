package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 用于替换gorm.Model,更灵活.
type BaseModel struct {
	ID        int32          `gorm:"primary_key"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt // 为了继承软删除
	IsDeleted bool           `gorm:"column:is_deleted"`
}

type Inventory struct {
	BaseModel

	Goods   int32 `gorm:"type:int;index"` // 商品,添加索引
	Stocks  int32 `gorm:"type:int"`       // 对应的库存;仓库和库存是有对应关系的,一个商品可能有多个仓库都有库存,可以引入仓库表...
	Version int32 `gorm:"type:int"`       // 用于乐观锁
}

// Type Stock struct { // 仓库
// 	BaseModel
// 	Name    string
// 	Address string
// }.
