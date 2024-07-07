package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/inventory_srv/proto"
	"sync"
)

var inventoryClient proto.InventoryClient
var conn *grpc.ClientConn

func TestSetInv(goodsId int32, num int32) {
	_, err := inventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     num,
	})
	if err != nil {
		panic(any(err))
	}
}

func TestSell(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := inventoryClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
			{GoodsId: 422, Num: 2},
			{GoodsId: 423, Num: 3},
			{GoodsId: 424, Num: 4},
		},
	})
	if err != nil {
		panic(any(err))
	}
	fmt.Println("库存扣减成功")
}

func TestReback() {
	_, err := inventoryClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 20},
		},
	})
	if err != nil {
		panic(any(err))
	}
	fmt.Println("库存归还成功")
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50053", grpc.WithInsecure())
	if err != nil {
		panic(any(err))
	}
	inventoryClient = proto.NewInventoryClient(conn)
}

func main() {
	Init()
	//var i int32
	//for i = 421; i < 841; i++ {
	//	TestSetInv(i,100)
	//}

	//TestSell()
	//TestReback()

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go TestSell(&wg)
	}
	wg.Wait()
	conn.Close()
}
