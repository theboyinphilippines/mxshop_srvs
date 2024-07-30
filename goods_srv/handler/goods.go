package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

// 商品列表中涉及es搜索，增删改商品，需要增加es操作
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

// 商品列表(通过点击一二三级类目去查询出商品)
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//使用es的目的是搜索出商品的id来，通过id拿到具体的字段信息是通过mysql来完成
	//我们使用es是用来做搜索的， 是否应该将所有的mysql字段全部在es中保存一份
	//es用来做搜索，这个时候我们一般只把搜索和过滤的字段信息保存到es中
	//es可以用来当做mysql使用， 但是实际上mysql和es之间是互补的关系， 一般mysql用来做存储使用，es用来做搜索使用
	//es想要提高性能， 就要将es的内存设置的够大， 或写入少点字段1k 2k

	var goodsListResponse proto.GoodsListResponse
	// 条件搜索
	var goods []model.Goods
	//这里全局变量db要换成局部变量 拼凑sql语句
	localDB := global.DB.Model(&model.Goods{})
	// 定义es复合查询 bool查询
	q := elastic.NewBoolQuery()
	if req.KeyWords != "" {
		//localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	// IsHot IsNew 默认为false
	if req.IsHot {
		//localDB = localDB.Where(&model.Goods{IsHot: req.IsHot})
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		//localDB = localDB.Where(&model.Goods{IsNew: req.IsNew})
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsHot))
	}
	if req.PriceMin > 0 {
		//localDB = localDB.Where("shop_price >= ?", req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		//localDB = localDB.Where("shop_price <= ?", req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		//localDB = localDB.Where(&model.Goods{BrandsID: req.Brand})
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	parentSpan := opentracing.SpanFromContext(ctx)
	topCategorySpan := opentracing.GlobalTracer().StartSpan("TopCategory", opentracing.ChildOf(parentSpan.Context()))
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
		//localDB = localDB.Where(fmt.Sprintf("category_id IN (%s)", subQuery)).Find(&goods)
		//var total int64
		//localDB.Count(&total)
		//goodsListResponse.Total = int32(total)

		//去数据库中查询出category ids
		categoryIds := make([]interface{}, 0)
		type Result struct {
			ID int32
		}
		var result []Result
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&result)
		for _, r := range result {
			categoryIds = append(categoryIds, r.ID)
		}
		// 生成满足categoryIds的匹配条件（切片传入多个值，用terms）
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}
	topCategorySpan.Finish()

	//用category ids去es中查询出goods，取出goods ids
	//分页
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	esSearchSpan := opentracing.GlobalTracer().StartSpan("esSearchSpan", opentracing.ChildOf(parentSpan.Context()))
	rsp, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}
	esSearchSpan.Finish()

	goodsListResponse.Total = int32(rsp.Hits.TotalHits.Value)
	goodsIds := make([]int32, 0)
	for _, value := range rsp.Hits.Hits {
		var esGoods model.EsGoods
		_ = json.Unmarshal(value.Source, &esGoods)
		goodsIds = append(goodsIds, esGoods.ID)
	}

	// 用goods ids 去数据库中查询出goods list
	//localDB.Preload("Category").Preload("Brands").Where(fmt.Sprintf("category_id IN (%s)", subQuery)).Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	goodsListSpan := opentracing.GlobalTracer().StartSpan("goodsListSpan", opentracing.ChildOf(parentSpan.Context()))
	localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	goodsListSpan.Finish()

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
	var result *gorm.DB
	if result = global.DB.Where("id IN ?", req.Id).Preload("Category").Preload("Brands").Find(&goods); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "没有查询到批量商品")
	}
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
	//判断category id是否存在
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	//判断brand id是否存在
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
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

	//mysql保存商品时，需要将商品保存到es，为了耦合性（方便使用和删除es功能），这里用gorm的钩子，方便
	//用事务来保持mysql，es的数据保存
	tx := global.DB.Begin()
	result := tx.Create(&good)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	goodsInfoResponse := modelToGoodsResponse(good)
	return &goodsInfoResponse, nil
}
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	var good model.Goods
	if result := global.DB.First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	tx := global.DB.Begin()
	result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &empty.Empty{}, nil
}
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var good model.Goods
	if result := global.DB.First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

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

	tx := global.DB.Begin()
	result := global.DB.Save(&good)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
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

// 测试下es的聚合
func TestESAggs() {
	q := elastic.NewBoolQuery()
	q = q.Filter(elastic.NewRangeQuery("market_price").Gte(10).Lte(20))
	martetPriceAggs := elastic.NewSumAggregation().Field("market_price")
	builder := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q)
	builder = builder.Aggregation("martetAggs", martetPriceAggs)
	searchResult, _ := builder.Pretty(true).Do(context.Background())

	total := searchResult.Hits.TotalHits.Value
	for _, value := range searchResult.Hits.Hits {
		zap.S().Infof("这是TestESAggs.value: %v", string(value.Source))
	}
	zap.S().Infof("这是TestESAggs.total: %d", total)
	for _, v := range searchResult.Aggregations {
		zap.S().Infof("这是TestESAggs.aggs.value: %v", string(v))
	}
}
