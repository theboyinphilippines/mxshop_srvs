package handler

import (
	"context"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

//func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
//	//var res proto.BrandListResponse
//	//var brands []model.Brands
//	//result := global.DB.Find(&brands)
//	//var brandList []*proto.BrandInfoResponse
//	//for _, brand := range brands {
//	//	brandInfoResponse := proto.BrandInfoResponse{
//	//		Id:   brand.ID,
//	//		Name: brand.Name,
//	//		Logo: brand.Logo,
//	//	}
//	//	brandList = append(brandList, &brandInfoResponse)
//	//}
//	//
//	//res.Total = int32(result.RowsAffected)
//	//res.Data = brandList
//	//return &res, nil
//
//	brandListResponse := proto.BrandListResponse{}
//
//	var brands []model.Brands
//	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//
//	var total int64
//	global.DB.Model(&model.Brands{}).Count(&total)
//	brandListResponse.Total = int32(total)
//
//	var brandResponses []*proto.BrandInfoResponse
//	for _, brand := range brands {
//		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
//			Id:  brand.ID,
//			Name: brand.Name,
//			Logo: brand.Logo,
//		})
//	}
//	brandListResponse.Data = brandResponses
//	return &brandListResponse, nil
//
//}

func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}
	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)
	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}
