package resource

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type MergeSameSkuProductService struct {
	eel.ServiceBase
}

func NewMergeSameSkuProductService(ctx context.Context) *MergeSameSkuProductService {
	service := new(MergeSameSkuProductService)
	service.Ctx = ctx
	return service
}

func (this *MergeSameSkuProductService) Merge(resources []business.IResource) []*SkuMergedProduct {
	product2merged := make(map[int]*SkuMergedProduct)
	
	for _, productResource := range resources {
		rawProductResource := productResource.GetRawResourceObject().(*ProductResource)
		poolProduct := rawProductResource.GetPoolProduct()
		
		var ok bool
		var merged *SkuMergedProduct
		productId := poolProduct.Id
		if merged, ok = product2merged[productId]; !ok {
			merged = &SkuMergedProduct{
				TotalCount: 0,
				TotalPrice: 0.0,
				// TotalWeight: 0.0,
				PoolProduct: poolProduct,
			}
			product2merged[productId] = merged
		}
		
		// merged.TotalWeight += poolProduct.GetSku(rawProductResource.Sku).Weight * float64(rawProductResource.Count)
		merged.TotalCount += rawProductResource.Count
		//merged.TotalPrice +=
	}
	
	//构造返回结果
	mergedProducts := make([]*SkuMergedProduct, 0)
	for _, mergedProduct := range product2merged {
		mergedProducts = append(mergedProducts, mergedProduct)
	}
	
	return mergedProducts
}


func init() {
}
