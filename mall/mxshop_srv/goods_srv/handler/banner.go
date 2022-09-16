package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/model"
	"mxshop_srv/goods_srv/proto"
)

func (g *GoodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	rsp := proto.BannerListResponse{}
	var banners []model.Banner
	result := global.DB.Find(&banners)     //结果写入[]model.Banner中
	rsp.Total = int32(result.RowsAffected) //检查总数

	var bannerRsp []*proto.BannerResponse //将model转换为proto定义的数据结构
	for _, banner := range banners {
		bannerRsp = append(bannerRsp, &proto.BannerResponse{
			Id:    banner.ID,
			Index: banner.Index,
			Image: banner.Image,
			Url:   banner.Url,
		})
	}

	rsp.Data = bannerRsp
	return &rsp, nil
}

func (g *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	//banner可以重名
	banner := model.Banner{Image: req.Image, Index: req.Index, Url: req.Url}

	global.DB.Save(&banner) //save当没有主键时创建,有主键时更新

	return &proto.BannerResponse{Id: banner.ID}, nil //更新成功返回一个主键
}

func (g *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	//先查询
	var banner model.Banner

	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}
	global.DB.Save(&banner)
	return &emptypb.Empty{}, nil

}
