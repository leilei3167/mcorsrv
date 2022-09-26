package router

import (
	"mxshop_api/goods_web/api/goods"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	goodsRouter := r.Group("goods")

	{
		goodsRouter.GET("", goods.List) // 商品列表
		// GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) // 该接口需要管理员权限

		goodsRouter.POST("", goods.New)               // 该接口需要管理员权限
		goodsRouter.GET("/:id", goods.Detail)         // 使用到路径变量
		goodsRouter.DELETE("/:id", goods.Delete)      // 需要管理员权限
		goodsRouter.PUT("/:id", goods.Update)         // 更新商品信息
		goodsRouter.PATCH("/:id", goods.UpdateStatus) // 更新状态
	}
}
