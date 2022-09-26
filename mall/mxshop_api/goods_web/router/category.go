package router

import (
	"mxshop_api/goods_web/api/category"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(r *gin.RouterGroup) {
	categoryRouter := r.Group("categorys")
	{
		categoryRouter.GET("", category.List)          // 列表
		categoryRouter.DELETE("/:id", category.Delete) // 删除
		categoryRouter.GET("/:id", category.Detail)    // 查询
		categoryRouter.POST("", category.New)
		categoryRouter.PUT("/:id", category.Update)
	}
}
