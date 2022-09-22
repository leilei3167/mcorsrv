package initialize

import (
	"net/http"

	"mxshop_api/goods_web/middlewares"
	"mxshop_api/goods_web/router"

	"github.com/gin-gonic/gin"
)

// Routers 负责初始化各种路由.
func Routers() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	r.Use(middlewares.Cors()) // 配置全局跨域
	ApiGroup := r.Group("/g/v1")
	// 向router中添加路由分组
	router.InitGoodsRouter(ApiGroup) // v1/user/

	return r
}
