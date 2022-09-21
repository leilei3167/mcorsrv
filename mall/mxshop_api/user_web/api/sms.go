package api

import (
	"context"
	"fmt"
	"math/rand"
	"mxshop_api/user_web/forms"
	"mxshop_api/user_web/global"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"

	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// SendSms 提供短信验证使用,阿里云自带的sdk,需要验证form表单的问题,如手机号码合法性等
func SendSms(c *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	zap.S().Debugf("sms配置:%#v", global.ServerConfig.AliSmsInfo)
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect)
	/* use STS Token
	credential := credentials.NewStsTokenCredential("<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		panic(err)
	}
	mobile := sendSmsForm.Mobile
	smsCode := GenerateSmsCode(6)

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"

	request.SignName = "阿里云短信测试"
	request.TemplateCode = "SMS_154950909"
	request.PhoneNumbers = mobile
	request.TemplateParam = "{\"code\":" + smsCode + "}"

	response, err := client.SendSms(request)
	if err != nil {
		zap.S().Errorf("发送短信错误:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "发送短信错误"})
		return
	}
	zap.S().Debugf("response is %#v\n", response)

	//以上代码可以发送出短信,需要将发出的验证码保存起来,便于后面来进行验证
	//一般手机号为key,验证码为value存入redis中
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	c.JSON(http.StatusOK, gin.H{"msg": "发送成功"})
}

// GenerateSmsCode 生成width长度的验证码
func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
