package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

//获取用户的购物车信息
func (o *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse
	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}

//将商品添加到购物车
func (o *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//有该商品记录，就加数量，没有记录就添加
	var shopCart model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{User: req.Id, Goods: req.GoodsId}).First(&shopCart); result.RowsAffected == 0 {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	} else {
		shopCart.Nums += req.Nums
	}
	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

//更新购物车记录（更新数量和是否选中）
func (o *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.First(&shopCart, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	shopCart.Checked = req.Checked
	global.DB.Save(&shopCart)
	return &empty.Empty{}, nil
}

//删除购物车记录
func (o *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &empty.Empty{}, nil
}

//订单列表
func (o *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	//分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		})
	}
	return &rsp, nil
}

//订单详情
func (o *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	rsp.OrderInfo = &proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Post:    order.Post,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SingerMobile,
	}

	var orderGoods []model.OrderGoods
	global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods)
	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			Id:         orderGood.ID,
			OrderId:    orderGood.Order,
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsImage: orderGood.GoodsImage,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Nums,
		})
	}
	return &rsp, nil
}
func GenerateOrderSn(userId int32) string {
	/*
		订单号生成规则
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(),
		now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, rand.Intn(90)+10)
	return orderSn
}

//新建订单（未使用分布式事务方案）
//func (o *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
//	/*
//		1.从购物车中获取选中的商品
//		2。商品的价格自己查询（用最新的价格结算）-访问商品服务（跨微服务）
//		3。库存的扣减-访问库存服务（跨微服务）
//		4。订单的基本信息表-订单的商品信息表（插入order和ordergood表）
//		5。从购物车中删除已购买的记录
//	*/
//
//	//从购物车中查询选中的商品
//	var shopCarts []model.ShoppingCart
//	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
//		return nil, status.Errorf(codes.InvalidArgument, "没有选中商品")
//	}
//	//获取选中的商品的goodsid，去goods_srv中查询价格
//	var goodsIds []int32
//	//根据商品id存储对应的数量，方便后面算订单总金额
//	goodNumMap := make(map[int32]int32)
//	for _, shopCart := range shopCarts {
//		goodsIds = append(goodsIds, shopCart.Goods)
//		goodNumMap[shopCart.Goods] = shopCart.Nums
//	}
//	//跨服务调用goods_srv，批量查询商品价格，算出订单总金额
//	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "批量查询商品失败")
//	}
//	var orderAmount float32
//	//顺便拿到订单商品（创建订单时，不只是插入order表，还有ordergood表）
//	var orderGoods []*model.OrderGoods
//	//顺便拿到goodsInfo，扣减库存服务需要
//	var goodsInvInfo []*proto.GoodsInvInfo
//	for _, good := range goods.Data {
//		orderAmount += good.ShopPrice * float32(goodNumMap[good.Id])
//		orderGoods = append(orderGoods, &model.OrderGoods{
//			Goods:      good.Id,
//			GoodsName:  good.Name,
//			GoodsImage: good.GoodsFrontImage,
//			GoodsPrice: good.ShopPrice,
//			Nums:       goodNumMap[good.Id],
//		})
//		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
//			GoodsId: good.Id,
//			Num:     goodNumMap[good.Id],
//		})
//	}
//
//	//跨服务调用inventory_srv，扣减库存
//	_, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsInvInfo})
//	if err != nil {
//		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
//	}
//
//	//插入订单表
//	//本地事务用于订单表，订单商品表，购物车表的操作（除了查询，增删改操作都要求事务）
//	tx := global.DB.Begin()
//	order := model.OrderInfo{
//		User:         req.UserId,
//		OrderSn:      GenerateOrderSn(req.UserId),
//		OrderMount:   orderAmount,
//		Address:      req.Address,
//		SignerName:   req.Name,
//		SingerMobile: req.Mobile,
//		Post:         req.Post,
//	}
//	if result := tx.Save(&order); result.RowsAffected == 0 {
//		tx.Rollback()
//		return nil, status.Errorf(codes.Internal, "创建订单失败")
//	}
//
//	//插入订单商品表
//	//插入订单表生成的order id 加入到订单商品表
//	for _, orderGood := range orderGoods {
//		orderGood.Order = order.ID
//	}
//
//	//批量插入订单商品表 CreateInBatches用于数据大时，分批插入
//	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
//		tx.Rollback()
//		return nil, status.Errorf(codes.Internal, "创建订单失败")
//	}
//
//	//删除购物车记录
//	if result := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
//		tx.Rollback()
//		return nil, status.Errorf(codes.Internal, "创建订单失败")
//	}
//	tx.Commit()
//	return &proto.OrderInfoResponse{
//		Id:      order.ID,
//		OrderSn: order.OrderSn,
//		Total:   order.OrderMount,
//	}, nil
//}

type OrderListener struct {
	code       codes.Code
	detail     string
	ID         int32
	OrderMount float32
}

//业务
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	//从购物车中查询选中的商品
	var shopCarts []model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		o.code = codes.InvalidArgument
		o.detail = "没有选中商品"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.InvalidArgument, "没有选中商品")
	}
	//获取选中的商品的goodsid，去goods_srv中查询价格
	var goodsIds []int32
	//根据商品id存储对应的数量，方便后面算订单总金额
	goodNumMap := make(map[int32]int32)
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodNumMap[shopCart.Goods] = shopCart.Nums
	}
	//跨服务调用goods_srv，批量查询商品价格，算出订单总金额
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.code = codes.Internal
		o.detail = "批量查询商品失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.Internal, "批量查询商品失败")
	}
	var orderAmount float32
	//顺便拿到订单商品（创建订单时，不只是插入order表，还有ordergood表）
	var orderGoods []*model.OrderGoods
	//顺便拿到goodsInfo，扣减库存服务需要
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodNumMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodNumMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodNumMap[good.Id],
		})
	}

	//跨服务调用inventory_srv，扣减库存
	_, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsInvInfo, OrderSn: orderInfo.OrderSn})
	if err != nil {
		o.code = codes.ResourceExhausted
		o.detail = "扣减库存失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}

	//插入订单表
	//本地事务用于订单表，订单商品表，购物车表的操作（除了查询，增删改操作都要求事务）
	tx := global.DB.Begin()
	//order := model.OrderInfo{
	//	User:         req.UserId,
	//	OrderSn:      GenerateOrderSn(req.UserId),
	//	OrderMount:   orderAmount,
	//	Address:      req.Address,
	//	SignerName:   req.Name,
	//	SingerMobile: req.Mobile,
	//	Post:         req.Post,
	//}
	orderInfo.OrderMount = orderAmount
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.code = codes.Internal
		o.detail = "创建订单失败"
		return primitive.CommitMessageState
		//return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	o.OrderMount = orderAmount
	o.ID = orderInfo.ID

	//插入订单商品表
	//插入订单表生成的order id 加入到订单商品表
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	//批量插入订单商品表 CreateInBatches用于数据大时，分批插入
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.code = codes.Internal
		o.detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
		//return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	//删除购物车记录
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.code = codes.Internal
		o.detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}

	//订单生成后，发送一条订单延迟消息，记录假如订单超时，也要将库存归还
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.101:9876"}))
	if err != nil {
		zap.S().Errorf("初始化延迟消息失败：%v", err)
	}
	err = p.Start()
	if err != nil {
		zap.S().Errorf("开始延迟消息失败：%v", err)
	}
	orderInfoString, _ := json.Marshal(&orderInfo)
	msg = primitive.NewMessage("order_timeout", orderInfoString)
	msg.WithDelayTimeLevel(16)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("延迟消息发送失败：%v", err)
		//延迟消息发送失败，也要回滚本地事务
		tx.Rollback()
		o.code = codes.Internal
		o.detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	tx.Commit()
	//整个事务提交后，创建状态码，可在createorder中判断状态

	o.code = codes.OK
	return primitive.RollbackMessageState
}

//回查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	if result := global.DB.Where(&model.OrderInfo{User: orderInfo.User, OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

//新建订单（使用分布式事务方案-基于事务消息的最终一致性）
func (o *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		1.从购物车中获取选中的商品
		2。商品的价格自己查询（用最新的价格结算）-访问商品服务（跨微服务）
		3。库存的扣减-访问库存服务（跨微服务）
		4。订单的基本信息表-订单的商品信息表（插入order和ordergood表）
		5。从购物车中删除已购买的记录
	*/

	//发送事务半消息
	var ordeInfo model.OrderInfo
	ordeInfo.User = req.UserId
	ordeInfo.OrderSn = GenerateOrderSn(req.UserId)
	ordeInfo.Address = req.Address
	ordeInfo.SignerName = req.Name
	ordeInfo.SingerMobile = req.Mobile
	ordeInfo.Post = req.Post
	ordeInfoString, _ := json.Marshal(&ordeInfo)

	orderListener := OrderListener{}
	p, err := rocketmq.NewTransactionProducer(&orderListener, producer.WithNameServer([]string{"192.168.0.101:9876"}))
	if err != nil {
		zap.S().Errorf("初始化事务消息失败：%v", err)
		return nil, err
	}
	err = p.Start()
	if err != nil {
		zap.S().Errorf("事务消息开始失败：%v", err)
		return nil, err
	}
	msg := primitive.NewMessage("order_reback", ordeInfoString)
	_, err = p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("事务消息发送失败：%v", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}
	//用于判断整个本地事务是否执行成功，不成功就返回错误玛
	if orderListener.code != codes.OK {
		return nil, status.Error(orderListener.code, orderListener.detail)
	}
	//成功就返回数据
	return &proto.OrderInfoResponse{
		Id:      orderListener.ID,
		OrderSn: ordeInfo.OrderSn,
		Total:   orderListener.OrderMount,
	}, nil
}

//更新订单状态
func (o *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &empty.Empty{}, nil
}

//库存归还（消费库存归还消息）
func OrderTimeOut(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type orderDetail struct {
		OrderSn string
	}
	for i := range msgs {
		var order orderDetail
		_ = json.Unmarshal(msgs[i].Body, &order)
		//tx := global.DB.Begin()
		//var sellDetail model.SellDetail
		//if result := tx.Where(&model.SellDetail{OrderSn: order.orderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
		//	return consumer.ConsumeRetryLater, nil
		//}
		//for _, detail := range sellDetail.Detail {
		//	if result := tx.Where(&model.Inventory{Goods: detail.GoodsId}).Update("stocks", gorm.Expr("stocks + ?", detail.Num)); result.RowsAffected == 0 {
		//		return consumer.ConsumeRetryLater, nil
		//	}
		//}
		//tx.Commit()

		//查看订单表中该订单的状态是否成功，如果不是成功状态就设置为超时关闭，发送归还库存消息让库存服务归还订单
		var orderInfo model.OrderInfo
		if result := global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: order.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
			return consumer.ConsumeRetryLater, nil
		}
		if orderInfo.Status != "TRADE_SUCCESS" {
			orderInfo.Status = "TRADE_CLOSED"
			global.DB.Save(&orderInfo)

			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.101:9876"}))
			if err != nil {
				zap.S().Errorf("初始化消息失败：%v", err)
			}
			err = p.Start()
			if err != nil {
				zap.S().Errorf("开始消息失败：%v", err)
			}
			msg := primitive.NewMessage("order_reback", msgs[i].Body)
			_, err = p.SendSync(context.Background(), msg)
			if err != nil {
				zap.S().Errorf("发送消息失败：%v", err)
				return consumer.ConsumeRetryLater, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
