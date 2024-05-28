package handler

import (
	"context"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

//设置库存
func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num
	global.DB.Save(&inv)
	return &empty.Empty{}, nil
}

//查询库存信息
func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{GoodsId: inv.Goods, Num: inv.Stocks}, nil
}

//扣减库存
//func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//本地事务（扣除的是多个商品，要么全部成功，要么全部失败，使用gorm手动事务）
//	//并发情况会出现超卖（分布式锁）
//	// 开始事务
//	tx := global.DB.Begin()
//	for _, goodInfo := range req.GoodsInfo {
//		var inv model.Inventory
//		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//			//回滚操作
//			tx.Rollback()
//			return nil, status.Errorf(codes.NotFound, "没有库存信息")
//		}
//		//判断库存是否充足
//		if inv.Stocks < goodInfo.Num {
//			//回滚操作
//			tx.Rollback()
//			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
//		}
//		//扣减库存
//		inv.Stocks -= goodInfo.Num
//		tx.Save(&inv)
//	}
//	//手动提交事务
//	tx.Commit()
//	return &empty.Empty{}, nil
//}

//扣减库存（互斥锁）
//var m sync.Mutex
//func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//本地事务（扣除的是多个商品，要么全部成功，要么全部失败，使用gorm手动事务）
//	//并发情况会出现超卖（分布式锁）
//	// 开始事务
//	tx := global.DB.Begin()
//	m.Lock() //获取锁 互斥锁最大问题是：性能问题。假设有10万并发，但并不是请求的同一件商品
//	for _, goodInfo := range req.GoodsInfo {
//		var inv model.Inventory
//		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//			//回滚操作
//			tx.Rollback()
//			return nil, status.Errorf(codes.NotFound, "没有库存信息")
//		}
//		//判断库存是否充足
//		if inv.Stocks < goodInfo.Num {
//			//回滚操作
//			tx.Rollback()
//			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
//		}
//		//扣减库存
//		inv.Stocks -= goodInfo.Num
//		tx.Save(&inv)
//	}
//	//手动提交事务
//	tx.Commit()
//	m.Unlock() //释放锁
//	return &empty.Empty{}, nil
//}

//扣减库存（mysql悲观锁 这里用行锁（条件：只锁住满足条件的行；筛选条件为索引，不是则升级为表锁） for update语句，使用gorm锁）
//func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//本地事务（扣除的是多个商品，要么全部成功，要么全部失败，使用gorm手动事务）
//	//并发情况会出现超卖（分布式锁）
//	// 开始事务
//	tx := global.DB.Begin()
//	for _, goodInfo := range req.GoodsInfo {
//		var inv model.Inventory
//		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//			//回滚操作
//			tx.Rollback()
//			return nil, status.Errorf(codes.NotFound, "没有库存信息")
//		}
//		//判断库存是否充足
//		if inv.Stocks < goodInfo.Num {
//			//回滚操作
//			tx.Rollback()
//			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
//		}
//		//扣减库存
//		inv.Stocks -= goodInfo.Num
//		tx.Save(&inv)
//	}
//	//手动提交事务
//	tx.Commit()
//	return &empty.Empty{}, nil
//}

//扣减库存（mysql乐观锁，加版本号，流程为：查询，业务，更新，重新执行上述流程）
//func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//本地事务（扣除的是多个商品，要么全部成功，要么全部失败，使用gorm手动事务）
//	//并发情况会出现超卖（分布式锁）
//	// 开始事务
//	tx := global.DB.Begin()
//	for _, goodInfo := range req.GoodsInfo {
//		var inv model.Inventory
//		//一直重复查询，直到版本号不相同
//		for {
//			if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//				//回滚操作
//				tx.Rollback()
//				return nil, status.Errorf(codes.NotFound, "没有库存信息")
//			}
//			//判断库存是否充足
//			if inv.Stocks < goodInfo.Num {
//				//回滚操作
//				tx.Rollback()
//				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
//			}
//			//扣减库存
//			inv.Stocks -= goodInfo.Num
//			// 更新语句 update inventory set stocks = stocks -1, version=version+1 where goods=goods and version=version
//			// 下面这种写法有瑕疵，不会更新0值，需要指定字段select
//			//tx.Model(&model.Inventory{}).Where("goods=? and version=?").Updates(model.Inventory{Goods: inv.Stocks, Version: inv.Version + 1})
//			if result := tx.Model(&model.Inventory{}).Select("stocks", "version").Where("goods=? and version=?",goodInfo.GoodsId,inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
//				zap.S().Info("库存扣减失败") //扣减失败就需要一直查询，用for循环
//			} else {
//				break
//			}
//		}
//	}
//	//手动提交事务
//	tx.Commit()
//	return &empty.Empty{}, nil
//}
//扣减库存（redis分布式锁）
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//本地事务（扣除的是多个商品，要么全部成功，要么全部失败，使用gorm手动事务）
	//并发情况会出现超卖（分布式锁）
	// 开始事务

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "127.0.0.1:6379",
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//获取分布式锁
		mutexName := fmt.Sprintf("goods_%d", goodInfo.GoodsId)
		mutex := rs.NewMutex(mutexName)

		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁失败")
		}
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			//回滚操作
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			//回滚操作
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减库存
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放分布式锁异常")
		}
	}
	//手动提交事务
	tx.Commit()
	return &empty.Empty{}, nil
}

//库存归还
func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//1 订单超时归还 2 订单创建失败，归还之前创建的库存 3 手动归还
	// 开始事务
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			//回滚操作
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}
		//归还库存
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	//手动提交事务
	tx.Commit()
	return &empty.Empty{}, nil
}
