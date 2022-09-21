package goods

import (
	"context"
	"errors"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
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

func List(c *gin.Context) {
	//商品列表主要是要过滤,大部分都是通过url的参数来传递;重要的是和前端约定好参数名称
	req := &proto.GoodsFilterRequest{}

	priceMin := c.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin) //当传入非法查询参数时,可以直接忽略而不是报错;默认为0(即不过滤)
	req.PriceMin = int32(priceMinInt)

	priceMax := c.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	req.PriceMax = int32(priceMaxInt)

	isHot := c.DefaultQuery("ih", "0")
	if isHot == "1" {
		req.IsHot = true
	}

	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		req.IsNew = true
	}

	isTab := c.DefaultQuery("it", "0")
	if isTab == "1" {
		req.IsTab = true
	}

	categoryId := c.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	req.TopCategory = int32(categoryIdInt)

	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	req.Pages = int32(pagesInt)

	perNums := c.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PagePerNums = int32(perNumsInt)

	keywords := c.DefaultQuery("q", "")
	req.KeyWords = keywords

	brandId := c.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	req.Brand = int32(brandIdInt)

	//调用服务
	r, err := global.GoodsSrvClient.GoodsList(context.Background(), req)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//proto的结果映射为业务层返回数据(因为有些数据不方便发给客户)
	reMap := make(map[string]any)
	reMap["total"] = r.Total
	reMap["data"] = r.Data
	c.JSON(http.StatusOK, reMap)
}
