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

func TestGetCategoryList() {
	rsp, err := brandClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.JsonData)
}

func TestGetSubCategoryList() {
	rsp, err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{Id: 130358, Level: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Info)
	fmt.Println(rsp.SubCategorys)
}
func TestCreateCategory() {
	rsp, err := brandClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           "寒带水果",
		Level:          2,
		IsTab:          false,
		ParentCategory: 130358,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
	fmt.Println(rsp.Name)
}

func TestDeleteCategory() {
	_, err := brandClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
		Id: 238010,
	})
	if err != nil {
		panic(err)
	}
}

func TestUpdateCategory() {
	_, err := brandClient.UpdateCategory(context.Background(), &proto.CategoryInfoRequest{
		Id: 238009,
		Name: "二颗莓",
		ParentCategory: 135487,
		Level: 3,
		IsTab: false,
	})
	if err != nil {
		panic(err)
	}
}

func TestCategoryBrandList() {
	_, err := brandClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
	Pages: 1,
	PagePerNums: 10,
	})
	if err != nil {
		panic(err)
	}
}
func TestGetCategoryBrandList(){
	rsp, err := brandClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: 130368,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)

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
	//TestCreateUser()
	TestGetSubCategoryList()
	//TestGetCategoryList()
	//TestCreateCategory()
	//TestDeleteCategory()
	//TestUpdateCategory()
	//TestCategoryBrandList()
	//TestGetCategoryBrandList()
	conn.Close()
}
