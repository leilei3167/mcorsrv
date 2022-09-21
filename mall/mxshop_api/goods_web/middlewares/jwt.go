package middlewares

import (
	"errors"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth 是处理JWT的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//自定义token字段为 X-token,前端需要把token储存,需要和后端商量过期时间,可以约定刷新令牌或重新登录
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请登录",
			})
			c.Abort() //不继续执行以后的代码
			return
		}

		//验证token
		j := NewJWT()
		//将token解析
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			c.Abort()
			return
		}
		//验证token成功,将解析的数据存入到context中,便于后续处理使用
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

type JWT struct {
	SigningKey []byte //用于签名的密钥
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(global.ServerConfig.JWTInfo.SigningKey)} //密钥需要配置
}

// CreateToken 根据传入的数据结构,使用密钥创建token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //SigningMethodES256是非对称加密,HS256是对称加密(安全性差,不常用)
	return token.SignedString(j.SigningKey)
}

// ParseToken 将字符串尝试解析为数据结构
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil //解析方法使用这个回调函数来提供验证的密钥。该函数接收已解析但未验证的Token。
			// 这允许你使用令牌头部的属性（如`kid'）来确定使用哪个密钥。
		})
	if err != nil { //解析token出错,判断(参考库文档)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		if clams, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return clams, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() //修改claims的过期时间,并以此创建新的token
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
