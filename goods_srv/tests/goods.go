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
		Id:         845,
		CategoryId: 225638,
		BrandId:    618,
		Name:       "百香果1",
	})
	if err != nil {
		panic(err)
	}

}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
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
	TestUpdateGoods()
	conn.Close()
}
