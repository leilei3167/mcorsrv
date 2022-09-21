package api

import (
	"context"
	"errors"
	"fmt"
	"mxshop_api/user_web/forms"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/global/response"
	"mxshop_api/user_web/middlewares"
	"mxshop_api/user_web/models"
	"mxshop_api/user_web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandleGRPCErrorToHttp 为了方便,可以将GRPC的内部错误码转换为http的面向用户的错误码
func HandleGRPCErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			switch s.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": s.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户服务不可用"})
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "该电话号码已注册"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "其他错误:" + s.Message() + "code:" + s.Code().String()})
			}
			return
		}
	}
}

func GetUserList(c *gin.Context) {

	if info, ok := c.Get("claims"); ok {
		zap.S().Infof("收到调用:%#v", info.(*models.CustomClaims)) //获取调用者的信息
	}
	//1.获取分页信息
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.UserSrvClient.GetUserList(c, &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList]查询用户列表失败")
		HandleGRPCErrorToHttp(err, c)
		return
	}

	//将结果反馈至前端,
	result := make([]any, 0)
	for _, value := range rsp.Data {

		//更规范
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Unix(int64(value.BirthDay), 0),
			Birthday: time.Unix(int64(value.BirthDay), 0).Format("2006-01-02"), //日期格式化打印
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)

}

func PasswordLogin(c *gin.Context) {
	//密码登录
	//1.表单验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		fmt.Printf("%t: %#v\n", err, err)
		HandleValidatorError(c, err)
		return
	}

	//此处进行验证码(store是包内的全局储存),第三个参数为true代表使用一次后将清除验证码缓存(即判断后验证码就失效)
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
		c.JSON(http.StatusBadRequest, gin.H{"captcha": "验证码错误"})
		return
	}

	if rsp, err := global.UserSrvClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		//查询号码没问题,检查密码;注意rsp中已经取得了哈希过的密码
		if _, passErr := global.UserSrvClient.CheckPassWord(c, &proto.CheckPasswordInfo{
			Password: passwordLoginForm.Password, EncryptedPassword: rsp.Password}); passErr != nil {
			if e, ok := status.FromError(passErr); ok {
				switch e.Code() {
				case codes.InvalidArgument:
					c.JSON(http.StatusBadRequest, gin.H{
						"password": "密码错误",
					})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "内部错误",
					})
				}
				return
			}
		} else { //登录成功
			j := middlewares.NewJWT()
			claims := models.CustomClaims{
				ID:         uint(rsp.Id),
				NickName:   rsp.NickName,
				AutorityId: uint(rsp.Role),
				StandardClaims: jwt.StandardClaims{
					//设置生效时间
					NotBefore: time.Now().Unix(),                     //签名的生效时间(立即生效)
					ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //过期时间,24小时
					Issuer:    "leilei",
				},
			}
			token, err := j.CreateToken(claims)
			if err != nil {
				zap.S().Debugf("%v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败", "err:": err.Error()})
				return
			}

			//将token写入并返回
			c.JSON(http.StatusOK, gin.H{
				"id":         rsp.Id,
				"nick_name":  rsp.NickName,
				"token":      token,
				"expired_at": time.Now().Add(time.Hour * 24).Unix(),
			})
			return
		}
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	//1.先获取用户注册所需的信息
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	//2.验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err != nil {
		zap.S().Errorf("redis.Get错误:%v", err)
		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "验证码错误"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			return
		}
	} else if value != registerForm.Code { //验证码错误
		c.JSON(http.StatusBadRequest, gin.H{"msg": "验证码错误"})
		return
	}

	//3.创建用户

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile, //默认nickname为电话
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register]新建用户失败:%s", err)
		HandleGRPCErrorToHttp(err, c)
		return
	}

	//4.注册成功,写token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:         uint(user.Id),
		NickName:   user.NickName,
		AutorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			//设置生效时间
			NotBefore: time.Now().Unix(),                     //签名的生效时间(立即生效)
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //过期时间,24小时
			Issuer:    "leilei",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		zap.S().Debugf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败", "err:": err.Error()})
		return
	}
	//将token写入并返回
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": time.Now().Add(time.Hour * 24).Unix(),
	})

}

func HandleValidatorError(c *gin.Context, err error) {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": removeTopStruct(errs.Translate(global.Trans)), //错误进行翻译
			//"error": errs.Translate(global.Trans), //错误进行翻译
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{ //非字段验证类型错误
		"msg": err.Error(),
	})
}

// "PasswordLoginForm.password": "password长度不能超过20个字符",删除前面的结构体名
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
