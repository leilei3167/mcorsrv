package main

import (
	"context"
	"fmt"
	"sync"

	"mxshop_srv/inventory_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	invClient proto.InventoryClient
	conn      *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("172.29.101.222:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	invClient = proto.NewInventoryClient(conn)
}

func main() {
	Init()
	// var i int32
	// for i = 421; i <= 840; i++ {
	// 	TestSetInv(i, 100)
	// }
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go TestSell(&wg)
	}

	wg.Wait()
}

func TestSetInv(goodsId, Num int32) {
	_, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     Num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestSell(wg *sync.WaitGroup) {
	/*
		1. 第一件扣减成功： 第二件： 1. 没有库存信息 2. 库存不足
		2. 两件都扣减成功
	*/
	defer wg.Done()
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
			{GoodsId: 422, Num: 1},
			// {GoodsId: 423, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}
