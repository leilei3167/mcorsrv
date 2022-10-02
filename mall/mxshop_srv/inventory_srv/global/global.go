package global

import (
	"mxshop_srv/inventory_srv/config"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

var (
	// DB 全局的数据库实例,handler层直接依赖,应该考虑解耦,方便后期更换数据库.
	DB           *gorm.DB
	RS           *redsync.Redsync
	ServerConfig *config.ServerConfig
	NacosConfig  = &config.NacosConfig{}
)
