package initialize

import (
	"net/http"

	"mxshop_api/user_web/middlewares"
	"mxshop_api/user_web/router"

	"github.com/gin-gonic/gin"
)

// Routers 负责初始化各种路由.
func Routers() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	r.Use(middlewares.Cors()) // 配置全局跨域
	ApiGroup := r.Group("/u/v1")
	// 向router中添加路由分组
	router.InitUserRouter(ApiGroup) // v1/user/
	// 添加baseRouter(验证码)
	router.InitBaseRouter(ApiGroup) // v1/base/captcha
	return r
}
