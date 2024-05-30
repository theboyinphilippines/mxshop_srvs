package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/userop_srv/proto"
)

var userFavClient proto.UserFavClient
var messageClient proto.MessageClient
var addressClient proto.AddressClient
var conn *grpc.ClientConn

func TestCreateMessage() {
	message, err := messageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      1,
		MessageType: 2,
		Subject:     "你好",
		Message:     "我想知道怎么购物",
		File:        "http:baidu.com",
	})
	if err != nil {
		return
	}
	fmt.Println(message)
}

func TestAddUserFav() {
	_, err := userFavClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  1,
		GoodsId: 421,
	})
	if err != nil {
		return
	}
}

func TestDeleteUserFav() {
	_, err := userFavClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  1,
		GoodsId: 421,
	})
	if err != nil {
		return
	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	addressClient = proto.NewAddressClient(conn)
	messageClient = proto.NewMessageClient(conn)
	userFavClient = proto.NewUserFavClient(conn)
}

func main() {
	Init()
	//TestCreateMessage()
	//TestAddUserFav()
	TestDeleteUserFav()
	conn.Close()
}
