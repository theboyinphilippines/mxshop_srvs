package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srv/proto"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func TestGoodsList() {
	rsp, err := brandClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin:    10,
		PriceMax:    30,
		Pages:       1,
		PagePerNums: 8,
		KeyWords:    "苹果",
		TopCategory: 130358,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, goodsInfo := range rsp.Data {
		fmt.Println(goodsInfo)
	}
}

func TestGetAllCategoryList() {
	rsp, err := brandClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
	fmt.Println(rsp.JsonData)
}

func TestBatchGetGoods() {
	rsp, err := brandClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}

func TestGetGoodsDetail() {
	rsp, err := brandClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Name)
}
func TestUpdateGoods() {
	_, err := brandClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:         847,
		CategoryId: 130358,
		BrandId:    614,
		Name:       "七七牌榴莲",
		GoodsBrief: "是在太好吃的东南亚进口好榴莲",
	})
	if err != nil {
		panic(err)
	}
}

func TestDeleteGoods() {
	_, err := brandClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: 847})
	if err != nil {
		panic(err)
	}
}

func TestCreateGoods() {
	_, err := brandClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            "妈妈们",
		GoodsSn:         "sdsd",
		MarketPrice:     27.54,
		ShopPrice:       25.24,
		GoodsBrief:      "嘻嘻嘻嘻",
		ShipFree:        false,
		Images:          []string{"http://www.amazon.com/iamges/01.jpg", "http://www.amazon.com/iamges/02.jpg"},
		DescImages:      []string{"http://www.amazon.com/iamges/01.jpg", "http://www.amazon.com/iamges/02.jpg"},
		GoodsFrontImage: "http://www.amazon.com/iamges",
		IsNew:           false,
		IsHot:           false,
		OnSale:          false,
		CategoryId:      130358,
		BrandId:         614,
	})
	if err != nil {
		panic(err)
	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	//TestGetBrandList()
	//TestGetAllCategoryList()
	//TestGoodsList()
	//TestBatchGetGoods()
	//TestGetGoodsDetail()
	//TestUpdateGoods()
	//TestDeleteGoods()
	TestCreateGoods()
	conn.Close()
}
