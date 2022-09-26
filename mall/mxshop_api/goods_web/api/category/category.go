package category

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	r, err := global.GoodsSrvClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[List] 查询 【分类列表】失败： ", err.Error())
	}

	c.JSON(http.StatusOK, data)
}

func Detail(c *gin.Context) { // 查询某个分类的详细情况(子分类)
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	// 回复interface在使用gin时非常方便
	reMap := make(map[string]interface{})
	subCategorys := make([]interface{}, 0)
	if r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	} else {
		// 处理其子分类
		for _, value := range r.SubCategorys {
			subCategorys = append(subCategorys, map[string]interface{}{
				"id":              value.Id,
				"name":            value.Name,
				"level":           value.Level,
				"parent_category": value.ParentCategory,
				"is_tab":          value.IsTab,
			})
		}
		reMap["id"] = r.Info.Id
		reMap["name"] = r.Info.Name
		reMap["level"] = r.Info.Level
		reMap["parent_category"] = r.Info.ParentCategory
		reMap["is_tab"] = r.Info.IsTab
		reMap["sub_categorys"] = subCategorys // 使用interface的好处 就是非常灵活

		c.JSON(http.StatusOK, reMap)
	}
}

func Delete(c *gin.Context) { // 删除也需要结合业务来确定,如是否需要删除该分类的所有子分类
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// 1. 先查询出该分类写的所有子分类
	// 2. 将所有的分类全部逻辑删除
	// 3. 将该分类下的所有的商品逻辑删除
	_, err = global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func New(c *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	c.JSON(http.StatusOK, request)
}

func Update(c *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}
