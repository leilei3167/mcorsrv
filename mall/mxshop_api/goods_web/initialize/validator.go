package initialize

import (
	"fmt"
	"mxshop_api/goods_web/global"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

//翻译器

func InitTrans(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		//分别对应,备用的语言环境,之后的2个是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		global.Trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("无法获取翻译器")
		}
		switch local {
		case "en":
			en_trans.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zh_trans.RegisterDefaultTranslations(v, global.Trans)
		default:
			en_trans.RegisterDefaultTranslations(v, global.Trans)
		}
	}
	return
}
