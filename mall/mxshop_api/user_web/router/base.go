package router

import (
	"github.com/gin-gonic/gin"
	"user_api/user_web/api"
)

func InitBaseRouter(r *gin.RouterGroup) {
	BaseRouter := r.Group("/base")
	{
		BaseRouter.GET("/captcha", api.GetCaptcha) //生成图片验证码
		BaseRouter.POST("/send_sms", api.SendSms)  //生成短信验证码
	}
}
