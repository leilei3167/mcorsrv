package main

import (
	"fmt"
	"user_api/user_web/global"
	"user_api/user_web/initialize"
	myvalidator "user_api/user_web/validator"

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
