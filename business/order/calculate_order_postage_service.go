package order

import (
	"context"
	"fmt"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/postage"
	"github.com/gingerxman/ginger-mall/business/order/resource"
	"math"
	"strconv"
	"strings"
)

// CalculateOrderPostageService 针对同一个供货商的商品集合，计算器运费
type CalculateOrderPostageService struct {
	eel.ServiceBase
}

func NewCalculateOrderPostageService(ctx context.Context) *CalculateOrderPostageService {
	service := new(CalculateOrderPostageService)
	service.Ctx = ctx
	return service
}

func (this *CalculateOrderPostageService) Calculate(productResources []business.IResource, purchaseInfo *PurchaseInfo) float64 {
	var useSystemPostageConfigProducts []*resource.SkuMergedProduct
	var totalPostageMoney float64 = 0.0
	
	mergedProducts := resource.NewMergeSameSkuProductService(this.Ctx).Merge(productResources)
	
	for _, mergedProduct := range mergedProducts {
		poolProduct := mergedProduct.PoolProduct
		
		if poolProduct.UseUnifiedPostage() {
			totalPostageMoney += poolProduct.GetUnifiedPostageMoney()
		} else {
			useSystemPostageConfigProducts = append(useSystemPostageConfigProducts, mergedProduct)
		}
	}
	
	if len(useSystemPostageConfigProducts) > 0 {
		supplierId := mergedProducts[0].PoolProduct.SupplierId
		supplierCorp := account.NewCorpFromOnlyId(this.Ctx, supplierId)
		postageConfig := postage.NewPostageConfigRepository(this.Ctx).GetActivePostageConfigInCorp(supplierCorp)
		if postageConfig == nil {
			eel.Logger.Error(fmt.Sprintf("no active postage config for corp(%d)"), supplierId)
		} else {
			postage.NewFillPostageConfigService(this.Ctx).FillOne(postageConfig, eel.FillOption{})
			
			freeConfig := this.getMatchedFreeConfig(useSystemPostageConfigProducts, purchaseInfo.ShipInfo, postageConfig)
			if freeConfig == nil {
				//没有匹配免邮条件，继续检查特殊地区设置
				specialAreaConfig := this.getMatchedSpecialAreaConfig(useSystemPostageConfigProducts, purchaseInfo.ShipInfo, postageConfig)
				if specialAreaConfig == nil {
					//没有特殊地区免邮条件
					totalPostageMoney += this.calculatePostage(useSystemPostageConfigProducts, postageConfig)
				} else {
					//使用SpecialAreaConfig作为PostageConfig
					tmpPostageConfig := &postage.PostageConfig{
						FirstWeight: specialAreaConfig.FirstWeight,
						FirstWeightPrice: specialAreaConfig.FirstWeightPrice,
						AddedWeight: specialAreaConfig.AddedWeight,
						AddedWeightPrice: specialAreaConfig.AddedWeightPrice,
					}
					totalPostageMoney += this.calculatePostage(useSystemPostageConfigProducts, tmpPostageConfig)
				}
			} else {
				//满足免邮条件
			}
		}
	}
	
	return totalPostageMoney
}

func (this *CalculateOrderPostageService) getMatchedSpecialAreaConfig(mergedProducts []*resource.SkuMergedProduct, shipInfo *ShipInfo, postageConfig *postage.PostageConfig) *postage.SpecialAreaPostageConfig {
	if !postageConfig.IsEnableSpecialConfig {
		return nil
	}
	
	provinceId := shipInfo.GetArea().Province.Id
	for _, specialAreaConfig := range postageConfig.SpecialAreaPostageConfigs {
		if this.isProvinceInDestinations(provinceId, specialAreaConfig.Destination) {
			return specialAreaConfig
		}
	}
	
	return nil
}

func (this *CalculateOrderPostageService) getMatchedFreeConfig(mergedProducts []*resource.SkuMergedProduct, shipInfo *ShipInfo, postageConfig *postage.PostageConfig) *postage.FreePostageConfig {
	if !postageConfig.IsEnableFreeConfig {
		return nil
	}
	
	provinceId := shipInfo.GetArea().Province.Id
	
	purchaseCount := 0
	for _, mergedProduct := range mergedProducts {
		purchaseCount += mergedProduct.TotalCount
	}
	
	for _, freeConfig := range postageConfig.FreePostageConfigs {
		if this.isProvinceInDestinations(provinceId, freeConfig.Destination) {
			if freeConfig.Condition == "count" {
				expectedCount, err := strconv.Atoi(freeConfig.ConditionValue)
				if err != nil {
					eel.Logger.Error(err)
				}
				
				if purchaseCount >= expectedCount {
					return freeConfig
				}
			}
		}
	}
	
	return nil
}

func (this *CalculateOrderPostageService) isProvinceInDestinations(provinceId string, destinations string) bool {
	items := strings.Split(destinations, ",")
	for _, item := range items {
		if item == provinceId {
			return true
		}
	}
	
	return false
}

func (this *CalculateOrderPostageService) calculatePostage(mergedProducts []*resource.SkuMergedProduct, postageConfig *postage.PostageConfig) float64 {
	firstWeight := postageConfig.FirstWeight
	firstWeightPrice := postageConfig.FirstWeightPrice
	addedWeight := postageConfig.AddedWeight
	addedWeightPrice := postageConfig.AddedWeightPrice
	
	if firstWeight == 0 && addedWeight == 0 {
		//免运费
		return 0.0
	}
	
	productWeight := 0.0
	for _, mergedProduct := range mergedProducts {
		productWeight += mergedProduct.TotalWeight
	}
	
	if productWeight <= firstWeight {
		return firstWeightPrice
	} else {
		productAddedWeight := productWeight - firstWeight
		addedWeightFactor := math.Ceil(productAddedWeight / addedWeight)
		
		return firstWeightPrice + addedWeightFactor * addedWeightPrice
	}
}


func init() {
}
