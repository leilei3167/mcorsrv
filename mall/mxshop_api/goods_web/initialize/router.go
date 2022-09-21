package initialize

import (
	"mxshop_api/goods_web/middlewares"
	"mxshop_api/goods_web/router"

	"github.com/gin-gonic/gin"
)

// Routers 负责初始化各种路由
func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors()) //配置全局跨域
	ApiGroup := r.Group("/g/v1")
	//向router中添加路由分组
	router.InitUserRouter(ApiGroup) // v1/user/

	return r
}
