package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 自定义参数验证器.

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则判断是否合法
	ok, _ := regexp.MatchString(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, mobile)
	return ok
}
