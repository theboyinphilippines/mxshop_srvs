package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func modelToGoodsResponse(good model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              good.ID,
		CategoryId:      good.CategoryID,
		Name:            good.Name,
		GoodsSn:         good.GoodsSn,
		ClickNum:        good.ClickNum,
		SoldNum:         good.SoldNum,
		FavNum:          good.FavNum,
		MarketPrice:     good.MarketPrice,
		ShopPrice:       good.ShopPrice,
		GoodsBrief:      good.GoodsBrief,
		ShipFree:        good.ShipFree,
		Images:          good.Images,
		DescImages:      good.DescImages,
		GoodsFrontImage: good.GoodsFrontImage,
		IsNew:           good.IsNew,
		IsHot:           good.IsHot,
		OnSale:          good.OnSale,
		//AddTime:good.BaseModel.CreatedAt
		Category: &proto.CategoryBriefInfoResponse{
			Id:   good.CategoryID,
			Name: good.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   good.BrandsID,
			Name: good.Brands.Name,
			Logo: good.Brands.Logo,
		},
	}
}

//商品列表(通过点击一二三级类目去查询出商品)
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var goodsListResponse proto.GoodsListResponse
	// 条件搜索
	var goods []model.Goods
	//这里全局变量db要换成局部变量 拼凑sql语句
	localDB := global.DB.Model(&model.Goods{})
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}
	// IsHot IsNew 默认为false
	if req.IsHot {
		localDB = localDB.Where(&model.Goods{IsHot: req.IsHot})
	}

	if req.IsNew {
		localDB = localDB.Where(&model.Goods{IsNew: req.IsNew})
	}
	if req.IsHot {
		localDB = localDB.Where(&model.Goods{IsNew: req.IsNew})
	}
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.Brand > 0 {
		localDB = localDB.Where(&model.Goods{BrandsID: req.Brand})
	}
	subQuery := ""
	if req.TopCategory > 0 {
		//通过点击一二三级category去查询商品

		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		//判断点击的是一级二级还是三级类目，提取出所有符合条件的商品的category_id，再去goods中查找出数据
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category where parent_category_id = %d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in %d", req.TopCategory)

		} else if category.Level == 3 {
			//localDB = localDB.Where(&model.Goods{CategoryID: category.ID}).Find(&goods)
			subQuery = fmt.Sprintf("select id from category where id = %d", req.TopCategory)
		}
		localDB = localDB.Where(fmt.Sprintf("category_id IN (%s)", subQuery)).Find(&goods)
		var total int64
		localDB.Count(&total)
		goodsListResponse.Total = int32(total)
		localDB.Preload("Category").Preload("Brands").Where(fmt.Sprintf("category_id IN (%s)", subQuery)).Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	}
	var goodsInfoResponse []*proto.GoodsInfoResponse
	for _, good := range goods {
		goodResponse := modelToGoodsResponse(good)
		goodsInfoResponse = append(goodsInfoResponse, &goodResponse)
	}
	goodsListResponse.Data = goodsInfoResponse
	return &goodsListResponse, nil

	//提取出category_id，自定义结构体接收数据
	//使用原生sql查询，用scan和自定义结构体Result
	//type Result struct {
	//	ID int32
	//}
	//var results []Result
	//var categoryIds []int32
	//global.DB.Model(&model.Category{}).Raw(subQuery).Scan(&results)
	//for _, result := range results {
	//	categoryIds = append(categoryIds,result.ID)
	//}
	//用提取出的category_id 去goods中查询商品

}
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goods []model.Goods
	result := global.DB.Where("id IN ?", req.Id).Preload("Category").Preload("Brands").Find(&goods)
	//result := global.DB.Find(&goods,req.Id) （根据主键查询）
	var goodsListResponse proto.GoodsListResponse
	goodsListResponse.Total = int32(result.RowsAffected)
	for _, good := range goods {
		goodsInfoResponse := modelToGoodsResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return &goodsListResponse, nil
}
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var good model.Goods
	good.Name = req.Name
	good.GoodsSn = req.GoodsSn
	good.MarketPrice = req.MarketPrice
	good.ShopPrice = req.ShopPrice
	good.GoodsBrief = req.GoodsBrief
	good.ShipFree = req.ShipFree
	good.Images = req.Images
	good.DescImages = req.DescImages
	good.GoodsFrontImage = req.GoodsFrontImage
	good.IsNew = req.IsNew
	good.IsHot = req.IsHot
	good.OnSale = req.OnSale
	good.CategoryID = req.CategoryId
	good.BrandsID = req.BrandId
	global.DB.Create(&good)
	goodsInfoResponse := modelToGoodsResponse(good)
	return &goodsInfoResponse, nil
}
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	var good model.Goods
	if result := global.DB.First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	global.DB.Delete(&model.Goods{}, req.Id)
	return &empty.Empty{}, nil
}
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var good model.Goods
	good.Name = req.Name
	good.GoodsSn = req.GoodsSn
	good.MarketPrice = req.MarketPrice
	good.ShopPrice = req.ShopPrice
	good.GoodsBrief = req.GoodsBrief
	good.ShipFree = req.ShipFree
	good.Images = req.Images
	good.DescImages = req.DescImages
	good.GoodsFrontImage = req.GoodsFrontImage
	good.IsNew = req.IsNew
	good.IsHot = req.IsHot
	good.OnSale = req.OnSale
	good.CategoryID = req.CategoryId
	good.BrandsID = req.BrandId
	global.DB.Save(&good)
	return &empty.Empty{}, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var good model.Goods
	if result := global.DB.Preload("Category").Preload("Brands").First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := modelToGoodsResponse(good)
	return &goodsInfoResponse, nil
}
