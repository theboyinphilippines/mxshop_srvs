package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

//获取商品分类
func (s *GoodsServer) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var category []model.Category
	//先查询一级类目，再查询该类目下的子类目
	result := global.DB.Model(&model.Category{}).Where("level=1").Preload("SubCategory.SubCategory").Find(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	b, _ := json.Marshal(&category)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

//获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var category model.Category
	//先判断是否有该分类
	result := global.DB.First(&category, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//这里proto.SubCategoryListResponse返回类型中没有SubCategory字段，所以这里只返回下一级的类目，不用返回下下一级的类目
	//preload := "SubCategory"
	//if req.Level == 1 {
	//	preload = "SubCategory.SubCategory"
	//}
	var categoryInfoResponse proto.CategoryInfoResponse
	categoryInfoResponse.Id = category.ID
	categoryInfoResponse.Name = category.Name
	categoryInfoResponse.Level = category.Level
	categoryInfoResponse.IsTab = category.IsTab
	categoryInfoResponse.ParentCategory = category.ParentCategoryID

	var subCategoryListResponse proto.SubCategoryListResponse
	subCategoryListResponse.Info = &categoryInfoResponse

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	//global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload("SubCategory").Find(&subCategorys)
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	subCategoryListResponse.SubCategorys = subCategoryResponse
	return &subCategoryListResponse, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	var category model.Category
	category.ID = req.Id
	category.Name = req.Name
	category.Level = req.Level
	category.IsTab = req.IsTab
	if req.Level != 1 {
		category.ParentCategoryID = req.ParentCategory
	}
	result := global.DB.Create(&category)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "服务目前不可用")
	}
	return &proto.CategoryInfoResponse{Id: category.ID}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	//先判断分类是否存在
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	global.DB.Delete(&model.Category{}, req.Id)
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	//先判断分类是否存在
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	// 判断数据库的category表里是否是默认值false
	category.IsTab = req.IsTab
	global.DB.Save(&category)
	return &empty.Empty{}, nil
}

// 品牌分类
func (s *GoodsServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var categoryBrands []model.GoodsCategoryBrand
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Preload("Brands").Preload("Category").Find(&categoryBrands)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total)
	var categoryBrandListResponse proto.CategoryBrandListResponse
	categoryBrandListResponse.Total = int32(total)

	var categoryBrandResponse []*proto.CategoryBrandResponse
	for _, categoryBrand := range categoryBrands {
		categoryBrandResponse = append(categoryBrandResponse, &proto.CategoryBrandResponse{
			Id: categoryBrand.ID,
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.BrandsID,
				Name: categoryBrand.Brands.Name,
				Logo: categoryBrand.Brands.Logo,
			},
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrand.CategoryID,
				Name:           categoryBrand.Category.Name,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
			},
		})
	}
	categoryBrandListResponse.Data = categoryBrandResponse
	return &categoryBrandListResponse, nil
}

//通过category获取brands
func (s GoodsServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var goodsCategoryBrands []model.GoodsCategoryBrand
	result := global.DB.Where(&model.GoodsCategoryBrand{CategoryID: req.Id}).Preload("Brands").Find(&goodsCategoryBrands)
	var brandListResponse proto.BrandListResponse
	brandListResponse.Total = int32(result.RowsAffected)
	var brandInfoResponse []*proto.BrandInfoResponse
	for _, goodsCategoryBrand := range goodsCategoryBrands {
		brandInfoResponse = append(brandInfoResponse, &proto.BrandInfoResponse{
			Id:   goodsCategoryBrand.BrandsID,
			Name: goodsCategoryBrand.Brands.Name,
			Logo: goodsCategoryBrand.Brands.Logo,
		})
	}
	brandListResponse.Data = brandInfoResponse
	return &brandListResponse, nil
}

func (s GoodsServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	//校验分类是否存在
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//校验品牌是否存在
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	var goodsCategoryBrand model.GoodsCategoryBrand
	goodsCategoryBrand.CategoryID = req.CategoryId
	goodsCategoryBrand.BrandsID = req.BrandId
	global.DB.Create(&goodsCategoryBrand)
	return &proto.CategoryBrandResponse{Id: goodsCategoryBrand.ID}, nil
}

func (s GoodsServer) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.GoodsCategoryBrand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s GoodsServer) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//校验分类是否存在
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//校验品牌是否存在
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	var goodsCategoryBrand model.GoodsCategoryBrand
	goodsCategoryBrand.CategoryID = req.CategoryId
	goodsCategoryBrand.BrandsID = req.BrandId
	global.DB.Save(&goodsCategoryBrand)
	return &emptypb.Empty{}, nil
}
