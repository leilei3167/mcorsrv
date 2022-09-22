package router

import (
	"mxshop_api/goods_web/api/goods"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	GoodsRouter := r.Group("goods")

	{
		GoodsRouter.GET("", goods.List) // 商品列表
		// GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) // 该接口需要管理员权限
		GoodsRouter.POST("", goods.New)       // 该接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail) // 使用到路径变量
	}
}
