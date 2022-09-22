package global

import (
	"mxshop_api/goods_web/config"
	"mxshop_api/goods_web/proto"

	ut "github.com/go-playground/universal-translator"
)

// 定义全局变量.

var (
	ServerConfig   = &config.ServerConfig{} // 全局配置
	Trans          ut.Translator
	GoodsSrvClient proto.GoodsClient

	NacosConfig = &config.NacosConfig{}
)
