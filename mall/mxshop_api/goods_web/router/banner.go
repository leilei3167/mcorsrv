package router

import (
	"mxshop_api/goods_web/api/banner"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banner.List)          // 轮播图列表页
		BannerRouter.DELETE("/:id", banner.Delete) // 删除轮播图
		BannerRouter.POST("", banner.New)          // 新建轮播图
		BannerRouter.PUT("/:id", banner.Update)    // 修改轮播图信息
	}
}
