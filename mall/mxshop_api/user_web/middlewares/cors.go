package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		//支持全域名访问，不安全，部署后需要固定限制为客户端网址
		c.Header("Access-Control-Allow-Origin", "*")
		//支持的Headers字段
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		//支持的Method
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")

		//如果想拿到非基本字段的其他字段，就必须在Access-Control-Expose-Headers里面指定
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		//该字段可选。它的值是一个布尔值，表示是否允许发送Cookie。默认情况下，Cookie不包括在CORS请求之中。设为true，即表示服务器明确许可，Cookie可以包含在请求中，一起发给服务器。
		//这个值也只能设为true，如果服务器不要浏览器发送Cookie，删除该字段即可。
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
