package product

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel"
	"strings"
)

type EncodePoolProductService struct {
	eel.ServiceBase
}

func NewEncodePoolProductService(ctx context.Context) *EncodePoolProductService {
	service := new(EncodePoolProductService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodePoolProductService) Encode(poolProduct *PoolProduct) *RPoolProduct {
	if poolProduct == nil {
		return nil
	}
	
	//对BaseInfo进行编码
	product := poolProduct.Product
	rProductBaseInfo := RProductBaseInfo{}
	rProductBaseInfo.Id = product.Id
	rProductBaseInfo.Name = product.Name
	rProductBaseInfo.Type = product.Type
	rProductBaseInfo.CreateType = poolProduct.Type
	rProductBaseInfo.PromotionTitle = product.PromotionTitle
	rProductBaseInfo.Code = product.Code
	if product.Description != nil {
		rProductBaseInfo.Detail = product.Description.Detail
	}
	rProductBaseInfo.ShelveType = poolProduct.Status
	rProductBaseInfo.Thumbnail = product.Thumbnail
	if strings.Index(rProductBaseInfo.Thumbnail, ".oss-") != -1 {
		thumbnail := strings.Replace(rProductBaseInfo.Thumbnail, "vxiaocheng-jh.oss-cn-beijing.aliyuncs.com", "resource.vxiaocheng.com", 1)
		thumbnail = fmt.Sprintf("%s?x-oss-process=image/resize,w_100/quality,q_80/interlace,1/format,jpg", thumbnail)
		rProductBaseInfo.Thumbnail = thumbnail
	}
	
	//对Medias进行编码
	rProductMedias := make([]*RProductMedia, 0)
	for _, media := range product.Medias {
		rProductMedias = append(rProductMedias, &RProductMedia{
			Id: media.Id,
			Type: media.Type,
			Url: media.Url,
		})
	}
	
	//对Skus进行编码
	rSkus := make([]*RProductSku, 0)
	for _, sku := range product.Skus {
		rSku := &RProductSku{
			Id: sku.Id,
			Name: sku.Name,
			Code: sku.Code,
			Price: sku.Price,
			CostPrice: sku.CostPrice,
			Stocks: sku.Stocks,
		}
		
		rPropertyValues := make([]*RProductPropertyValue, 0)
		buf := make([]string, 0)
		for _, propertyValue := range sku.PropertyValues {
			rPropertyValue := &RProductPropertyValue{
				Id: propertyValue.Id,
				PropertyId: propertyValue.PropertyId,
				PropertyName: propertyValue.PropertyName,
				Text: propertyValue.Text,
				Image: propertyValue.Image,
			}
			rPropertyValues = append(rPropertyValues, rPropertyValue)
			buf = append(buf, propertyValue.Text)
		}
		rSku.PropertyValues = rPropertyValues
		rSku.DisplayName = strings.Join(buf, " ")
		
		rSkus = append(rSkus, rSku)
	}
	
	//对Labels进行编码
	rLabels := make([]*RProductLabel, 0)
	if len(product.Labels) > 0 {
		rLabels = NewEncodeProductLabelService(this.Ctx).EncodeMany(product.Labels)
	} else {
		rLabels = make([]*RProductLabel, 0)
	}
	
	//对Category进行编码
	var rCategory *RLintProductCategory
	if len(product.Categories) > 0 {
		//TODO: 将product.Categories改为product.Category
		category := product.Categories[0]
		rCategory = &RLintProductCategory{
			Id: category.Id,
			Name: category.Name,
		}
	}
	
	//对logistics info进行编码
	//supplierCorp := account.NewCorpFromOnlyId(this.Ctx, poolProduct.SupplierId)
	//postageConfig := postage.NewPostageConfigRepository(this.Ctx).GetActivePostageConfigInCorp(supplierCorp)
	//if postageConfig == nil {
	//	postage.NewPostageConfigFactory(this.Ctx).MakeSureDefaultPostageConfigExits(supplierCorp)
	//	postageConfig = postage.NewPostageConfigRepository(this.Ctx).GetActivePostageConfigInCorp(supplierCorp)
	//}
	//postage.NewFillPostageConfigService(this.Ctx).FillOne(postageConfig, eel.FillOption{})
	//rPostageConfig := postage.NewEncodePostageConfigService(this.Ctx).Encode(postageConfig)
	rLogisticsInfo := RProductLogisticsInfo{
		PostageType: product.PostageType,
		UnifiedPostageMoney: product.UnifiedPostageMoney,
		//PostageConfig: rPostageConfig,
		LimitZoneType: product.LimitZoneType,
		LimitZoneTypeCode: product.GetLimitZontTypeText(),
		LimitZoneId: product.LimitZoneId,
		LimitZoneAreas: product.GetLimitZoneAreas(),
	}

	return &RPoolProduct{
		Id: poolProduct.Id,
		CorpId: poolProduct.CorpId,
		Type: poolProduct.Type,
		BaseInfo: &rProductBaseInfo,
		LogisticsInfo: &rLogisticsInfo,
		Medias: rProductMedias,
		Category: rCategory,
		SoldCount: poolProduct.SoldCount,
		VisitInfo: &RVisitInfo{
			UserCount: 5,
			ViewCount: 20,
		},
		Skus: rSkus,
		Labels: rLabels,
		IsDeleted: false,
		Status: poolProduct.Status,
		CreatedAt: poolProduct.CreatedAt.Format("2006-01-02 15:04"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodePoolProductService) EncodeMany(poolProducts []*PoolProduct) []*RPoolProduct {
	rDatas := make([]*RPoolProduct, 0)
	for _, poolProduct := range poolProducts {
		rDatas = append(rDatas, this.Encode(poolProduct))
	}
	
	return rDatas
}

func init() {
}
