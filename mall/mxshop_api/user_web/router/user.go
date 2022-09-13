package router

import (
	"user_api/user_web/api"
	"user_api/user_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	UserGroup := r.Group("user")

	{
		//	UserGroup.GET("/list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList) //api添加jwt验证
		UserGroup.GET("/list", middlewares.JWTAuth(), api.GetUserList) //api添加jwt验证
		UserGroup.POST("/pwd_login", api.PasswordLogin)
		UserGroup.POST("/register", api.Register)
	}

}
