package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// api包内使用的验证码储存方式.
var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80) // 数字验证码,可配置高度宽度等
	cp := base64Captcha.NewCaptcha(driver, store)               // 通过driver放入store中
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成验证码错误"})
		return
	}
	// 将生成的验证码返回
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,   // 此次验证码的id
		"picPath":   b64s, // 验证码图片的base64编码,前端拿到可直接转换为图片
	})
}
