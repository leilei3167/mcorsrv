package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/model"
	"mxshop_srv/goods_srv/proto"
)

// GetAllCategorysList 获取所有的分类数据,如何将数据组织成前端易于使用的数据? 体现层级关系而不是一股脑列出
//一般来说service层不应该负责数据组织的功能,应该尽量简单,由web层去组织;此处 返回组织好的json数据和原始数据,供调用处选择
func (g *GoodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {

	var categorys []model.Category
	//global.DB.Find(&categorys) //拿原始数据很简单
	//global.DB.Preload("SubCategory").Find(&categorys) //预加载,其实就是反向查询,这种方式有2个问题: 1.他会将所有分级的的结果返回;2.只会向下查询一级,也就是3级无法体现

	//嵌套preload
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	//其实底层执行了3条sql语句,1.首先where查询所有level=1的一级类目;2.查询parent_id在第一步结果中的二级类目;3.查询parent_id在第二步结果中的三级类目
	b, _ := json.Marshal(categorys)

	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListRsp := proto.SubCategoryListResponse{}

	//先查询该分类是否存在
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//如果商品分类存在,写入查询结果
	categoryListRsp.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	//然后查询该分类的子分类
	var subCategorys []model.Category
	var subCategorysRsp []*proto.CategoryInfoResponse
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		//转换
		subCategorysRsp = append(subCategorysRsp, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			ParentCategory: subCategory.ParentCategoryID,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
		})
	}

	categoryListRsp.SubCategorys = subCategorysRsp
	return &categoryListRsp, nil
}

func (g *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}

	cMap := map[string]any{} //使用create和map 的形式
	cMap["name"] = req.Name
	cMap["level"] = req.Level
	cMap["is_tab"] = req.IsTab
	if req.Level != 1 {
		cMap["parent_category_id"] = req.ParentCategory
	}
	tx := global.DB.Model(&model.Category{}).Create(cMap)
	fmt.Println(tx)

	////使用Save 无法解决外键的问题
	//category.Name = req.Name
	//category.Level = req.Level
	//category.IsTab = req.IsTab
	//if req.Level != 1 { //如果他本身不是一级类目,则设置他的父级目录
	//	category.ParentCategoryID = req.ParentCategory
	//}
	//global.DB.Save(&category)

	return &proto.CategoryInfoResponse{Id: category.ID}, nil

}

func (g *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	category := model.Category{}

	//更新时要先查询该分类是否存在
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	//获取要更新的字段,判断的原因是,避免将req中未设置的零值存入数据库
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)
	return &emptypb.Empty{}, nil

}
