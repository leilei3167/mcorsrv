package router

import (
	"mxshop_api/goods_web/api/goods"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	GoodsRouter := r.Group("goods")

	{
		GoodsRouter.GET("/list", goods.List)
	}

}
