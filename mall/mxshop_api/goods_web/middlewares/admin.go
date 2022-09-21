package middlewares

import (
	"mxshop_api/goods_web/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IsAdminAuth 验证是否是管理员
func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if ok {
			if currentUser, ok := claims.(*models.CustomClaims); ok {
				if currentUser.ID != 2 {
					//管理员
					c.JSON(http.StatusForbidden, gin.H{"msg": "无权限"})
					c.Abort()
					return
				}
				c.Next()
			}
		}
	}
}
