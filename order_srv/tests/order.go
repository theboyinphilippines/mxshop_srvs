package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/order_srv/proto"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func TestCreateCartItem() {
	_, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		GoodsId: 422,
		Nums:    30,
	})
	if err != nil {
		panic(err)
	}
}
func TestCartItemList() {
	rsp, _ := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: 1,
	})
	fmt.Println(rsp.Total)
}

func TestUpdateCartItem() {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      1,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder() {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "beijing",
		Name:    "shy",
		Mobile:  "18758782858",
		Post:    "kuaidian",
	})
	if err != nil {
		panic(err)
	}
}

func TestOrderDetail() {
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: 2,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Goods)
	fmt.Println(rsp.OrderInfo)
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func main() {
	Init()
	//TestCreateCartItem()
	//TestCartItemList()
	//TestUpdateCartItem()
	//TestCreateOrder()
	TestOrderDetail()
	conn.Close()
}
