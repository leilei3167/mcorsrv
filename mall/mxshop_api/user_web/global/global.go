package global

import (
	"mxshop_api/user_web/config"
	"mxshop_api/user_web/proto"

	ut "github.com/go-playground/universal-translator"
)

//定义全局变量

var (
	ServerConfig  = &config.ServerConfig{} //全局配置
	Trans         ut.Translator
	UserSrvClient proto.UserClient

	NacosConfig = &config.NacosConfig{}
)
