package handler

import (
	"gorm.io/gorm"
	"mxshop_srv/goods_srv/proto"
)

// Paginate 官方文档中的分页
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

// ModelToResponse 方便将model.User转换为proto的数据
/*func ModelToResponse(user model.User) proto.UserInfoResponse {
	//grpc中message有默认值,不能随便赋值
	//要搞清哪些字段有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		//BirthDay: user.Birthday, //注意时间的转换,Birthday可能为nil 本身是*Time类型
		Gender: user.Gender,
		Role:   int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}
*/

type GoodsServer struct {
	proto.UnimplementedGoodsServer //可以临时使用,快速启动grpcserver
}

var _ proto.GoodsServer = (*GoodsServer)(nil)

/*
func (g *GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
*/
