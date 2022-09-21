package main

import (
	"fmt"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/initialize"
	"mxshop_api/goods_web/utils"
	myvalidator "mxshop_api/goods_web/validator"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {

	//1.初始化logger
	initialize.InitLogger()
	//2.初始化配置
	initialize.InitConfig()

	//3.初始化router
	e := initialize.Routers()
	//4.初始化翻译器
	err := initialize.InitTrans("zh")
	if err != nil {
		panic(err)
	}
	//5.初始化srv的连接
	initialize.InitSrvConn()

	//6.动态获取端口号
	viper.AutomaticEnv()
	//如果是本地开发环境端口号固定，线上环境启动获取端口号
	release := viper.GetBool("MXSHOP_DEBUG")
	if release { //上线使用动态port
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port //修改
		}
	}

	//注册自定义字段验证,以及注册翻译
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}手机号码非法", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
	if err := e.Run(fmt.Sprintf("127.0.0.1:%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic(err)
	}
}
