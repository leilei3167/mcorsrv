package goods

import (
	"context"
	"net/http"
	"strconv"

	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	// 商品列表主要是要过滤,大部分都是通过url的参数来传递;重要的是和前端约定好参数名称
	req := &proto.GoodsFilterRequest{}

	priceMin := c.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin) // 当传入非法查询参数时,可以直接忽略而不是报错;默认为0(即不过滤)
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

	// 调用服务
	r, err := global.GoodsSrvClient.GoodsList(context.Background(), req)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	// proto的结果映射为业务层返回数据(因为有些数据不方便发给客户)
	reMap := make(map[string]any)
	reMap["total"] = r.Total

	goodsList := make([]any, 0)
	for _, r := range r.Data {
		goodsList = append(goodsList, map[string]any{
			"id":          r.Id,
			"name":        r.Name,
			"goods_brief": r.GoodsBrief,
			"desc":        r.GoodsDesc,
			"ship_free":   r.ShipFree,
			"images":      r.Images,
			"des_images":  r.DescImages,
			"front_image": r.GoodsFrontImage,
			"shop_price":  r.ShopPrice,
			"category": map[string]any{
				"id":   r.Category.Id,
				"name": r.Category.Name,
			},
			"brand": map[string]any{
				"id":   r.Brand.Id,
				"name": r.Brand.Name,
				"logo": r.Brand.Logo,
			},
			"is_hot":  r.IsHot,
			"is_new":  r.IsNew,
			"on_sale": r.OnSale,
		})
	}
	reMap["data"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

// New 新建商品,注意需要管理员权限.
func New(c *gin.Context) {
	// 先拿要创建的商品的信息
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBind(&goodsForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(c, &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree, // 表单验证需要传递指针
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		// 在新建时无需设置是否热门 促销等等,应该放在更新
		CategoryId: goodsForm.CategoryId,
		BrandId:    goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	// TODO: 如何设置库存? 需要单独的商品库存服务
	// 对于数据库的写入,实际是非常复杂的,尤其是涉及到多个服务的数据库,涉及到分布式事务
	c.JSON(http.StatusOK, rsp)
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	r, err := global.GoodsSrvClient.GetGoodsDetail(c, &proto.GoodInfoRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	// 对于有些数据是来自于库存微服务的,就需要单独去库存服务查(可以单独写一个单独的查询库存的web接口,供前端异步调用即可)
	// 而不是在一个接口中集成所有,大型项目不可能一个接口拿完所有部分的数据
	rsp := map[string]any{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"des_images":  r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"category": map[string]any{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]any{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}

	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(c, &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}

func UpdateStatus(c *gin.Context) {
	// 和创建类似
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func Stock(c *gin.Context) { // 库存服务作为一个单独的库存接口,用于获取商品的库存
	id := c.Param("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	// TODO: 编写库存服务
}
