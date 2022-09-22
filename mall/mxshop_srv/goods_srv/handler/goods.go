package handler

import (
	"context"
	"fmt"

	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/model"
	"mxshop_srv/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer // 可以临时使用,快速启动grpcserver
}

var _ proto.GoodsServer = (*GoodsServer)(nil)

// Paginate 官方文档中的分页.
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{ // 简化的信息,不需要所有
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

// GoodsList 商品列表页,关键是需要涉及到多种的过滤条件,如 关键词搜索,查询新品,查询热门商品,通过价格区间查询,根据分类筛选等等等
// 需要根据传入的字段来构建sql语句.
func (g *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods
	localDB := global.DB.Model(&model.Goods{}) // 根据传入的条件,一步一步增加查询语句
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
	}
	if req.IsNew {
		localDB = localDB.Where(model.Goods{IsNew: true})
	}

	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}

	if req.Brand > 0 {
		localDB = localDB.Where("brands_id = ?", req.Brand)
	}

	// 通过category设置查找条件;当设置为1级分类时,需要查询2,3级下的商品,用到子查询即可

	if req.TopCategory > 0 { // 查询确定该分类存在
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		var subQuery string
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		localDB = localDB.Where(fmt.Sprintf("category_id IN (%s)", subQuery))
	}
	var count int64
	localDB.Count(&count)

	goodsListResponse.Total = int32(count)

	// 需要取出外键的数据
	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}
	// 转换
	for _, good := range goods {
		goodsInfoRsp := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoRsp)
	}
	return goodsListResponse, nil
}

func (g *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	// 传入id切片 查询一批
	var goods []model.Goods
	result := global.DB.Find(&goods, req.Id)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	// 先确定是否已存在
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}
	global.DB.Save(&goods)
	return &proto.GoodsInfoResponse{Id: goods.ID}, nil
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	// 先确定是否已存在
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}

	var category model.Category // 查询待修改的 category 确保存在
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	global.DB.Save(&goods)

	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}
