package handler

import (
	"context"
	"fmt"

	"mxshop_srv/inventory_srv/global"
	"mxshop_srv/inventory_srv/model"
	"mxshop_srv/inventory_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

var _ proto.InventoryServer = (*InventoryServer)(nil)

func (i *InventoryServer) SetInv(_ context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 设置库存,先查询
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId // 直接设置 并更新
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	// 先查询
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// 扣减库存是库存服务的重点,涉及到分布式的事务.

func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 依次扣减库存
	// v1:
	// for _, goodInfo := range req.GoodsInfo { // 此种写法有几个隐患:1.并发问题;2.本地事务都没有实现(全部扣减成功或全部失败)
	// 	var inv model.Inventory
	// 	if result := global.DB.First(&inv, goodInfo.GoodsId); result.RowsAffected == 0 {
	// 		return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
	// 	}

	// 	if inv.Stocks < goodInfo.Num { // 库存小于了要扣减的数量
	// 		return nil, status.Errorf(codes.InvalidArgument, "库存不足")
	// 	}

	// 	// 扣减库存,并更新
	// 	inv.Stocks -= goodInfo.Num
	// 	global.DB.Save(&inv)
	// }

	// v2:本地事务,k遇到错误立刻Rollback;执行完毕后必须Commit
	// tx := global.DB.Begin()
	// for _, goodInfo := range req.GoodsInfo { // 此种写法不发解决并发情况的问题:超卖;因为可以多个请求同时执行扣减库存,会引起数据不一致
	// 	// for update 悲观锁
	// 	var inv model.Inventory
	// 	// if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
	// 	// 	tx.Rollback()
	// 	// 	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
	// 	// }

	// 	// 乐观锁,原理是开始时先查询到一个版本数据,之后更新时对比版本是否还是自己之前的
	// 	for {
	// 		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
	// 			tx.Rollback()
	// 			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
	// 		}

	// 		if inv.Stocks < goodInfo.Num { // 库存小于了要扣减的数量
	// 			tx.Rollback()
	// 			return nil, status.Errorf(codes.InvalidArgument, "库存不足")
	// 		}

	// 		// 扣减库存,并更新
	// 		inv.Stocks -= goodInfo.Num
	// 		if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods=? AND version =?",
	// 			goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
	// 			// 当抢购商品库存只有1时,第一个执行的g会将库存量设置为0,即成了零值,会被gorm忽略掉;使用Select方法指定强制更新的字段
	// 			zap.S().Info("库存扣减失败")
	// 		} else {
	// 			break // 更新成功 退出循环
	// 		}

	// 	}

	// }
	// tx.Commit()

	// v3: 基于redis的分布式锁

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		mName := fmt.Sprintf("mu_goods_%d", goodInfo.GoodsId)
		mutex := global.RS.NewMutex(mName)

		if err := mutex.Lock(); err != nil {
			// zap.S().Errorf("获取redis锁错误:%v", err)
			return nil, status.Errorf(codes.Internal, "获取redis锁错误")
		}

		m--
		fmt.Println("[m]:", m)

		// FIXME:此处并发查询会取得相同的库存值

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		zap.S().Debugf("got stock:[%v]", inv.Stocks)

		if inv.Stocks < goodInfo.Num { // 库存小于了要扣减的数量
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存不足")
		}

		// 扣减库存,并更新
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis锁错误")
		}
	}
	tx.Commit()

	return &emptypb.Empty{}, nil
}

var m = 100 // mu sync.Mutex

func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 库存归还: 1.订单超时归还;2.订单创建失败时归还之前扣减的库存;3.收到归还(可选)
	// 归还也是批量的归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		// 扣减库存,并更新
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
