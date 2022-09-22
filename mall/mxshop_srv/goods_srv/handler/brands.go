package handler

import (
	"context"

	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/model"
	"mxshop_srv/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListRsp := proto.BrandListResponse{}
	var brands []model.Brands
	r := global.DB.Scopes(model.Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if r.Error != nil {
		return nil, r.Error
	}
	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands { // 需要将model的结构转换为proto定义的
		brandResponse := proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		}
		brandResponses = append(brandResponses, &brandResponse)
	}
	brandListRsp.Data = brandResponses

	// 查询总数
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListRsp.Total = int32(total)
	// brandListRsp.Total = int32(r.RowsAffected) //返回的应该是品牌的总数,而不是这次分页的品牌个数
	return &brandListRsp, nil
}

func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 创建时应该先查询,确保品牌名称不能重复
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{Name: req.Name, Logo: req.Logo}
	global.DB.Save(brand) // Save可以用Create;Save既可以创建又可以更新

	return &proto.BrandInfoResponse{Id: brand.ID, Name: brand.Name, Logo: brand.Logo}, nil
}

func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if res := global.DB.Delete(&model.Brands{}, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "该品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	global.DB.Save(&brands)

	return &emptypb.Empty{}, nil
}
