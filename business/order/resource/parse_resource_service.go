package resource

import (
	"context"
	"encoding/json"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ParseResourceService struct {
	eel.ServiceBase
}

func NewParseResourceService(ctx context.Context) *ParseResourceService {
	service := new(ParseResourceService)
	service.Ctx = ctx
	return service
}

// 从订单数据中反序列化资源 todo
func (this *ParseResourceService) ParseFromOrderResources(resourcesData []map[string]interface{}) []business.IResource{
	resources := make([]business.IResource, 0)
	productArray := make([]interface{}, 0)
	imoneyArray := make([]interface{}, 0)
	couponArray := make([]interface{}, 0)
	for _, resourceData := range resourcesData{
		resourceType := resourceData["type"]
		switch resourceType {
		case RESOURCE_TYPE_PRODUCT:
			productArray = append(productArray, resourceData)
		case RESOURCE_TYPE_IMONEY:
			imoneyArray = append(imoneyArray, resourceData)
		case RESOURCE_TYPE_COUPON:
			couponArray = append(couponArray, resourceData)
		}
	}

	productResources := this.parseProductFromOrderResource(productArray)
	for _, productResource := range productResources{
		resources = append(resources, productResource)
	}

	supplierId := productResources[0].GetPoolProduct().SupplierId
	user := account.GetUserFromContext(this.Ctx)
	imoneyResources := this.parseIMoneyResource(imoneyArray, user.Id, supplierId)
	for _, imoneyResource := range imoneyResources{
		resources = append(resources, imoneyResource)
	}
	//couponResources := this.parseCouponResource(couponArray)

	for _, resource := range resources{
		resource.SetAllocated()
	}

	return resources
}

func (this *ParseResourceService) parseProductFromOrderResource(productsArray []interface{}) []*ProductResource {
	productResources := make([]*ProductResource, 0)
	for _, item := range productsArray {
	data := item.(map[string]interface{})

	productResource := NewProductResource(this.Ctx)
	productData := data["product"].(map[string]interface{})
	poolProductId, _ := productData["pool_product_id"].(json.Number).Int64()
	productResource.PoolProductId = int(poolProductId)

	count, _ := data["count"].(json.Number).Int64()
	productResource.Count = int(count)

	price, _ := data["price"].(json.Number).Float64()
	productResource.Price = price

	productResource.Sku = productData["sku_name"].(string)

	productResources = append(productResources, productResource)
	}

	this.fillProductResources(productResources)
	return productResources
}

func (this *ParseResourceService) Parse(salesmanId int, productsArray []interface{}, imoneysArray []interface{}, couponUsage *CouponUsage) []business.IResource {
	resources := make([]business.IResource, 0)
	
	//解析product resource
	productResources := this.parseProductResource(productsArray)
	for _, productResource := range productResources {
		resources = append(resources, productResource)
	}
	for _, productResource := range productResources {
		productResource.SalesmanId = salesmanId
	}
	
	supplierId := productResources[0].GetPoolProduct().SupplierId
	user := account.GetUserFromContext(this.Ctx)
	//解析imoney resource
	imoneyResources := this.parseIMoneyResource(imoneysArray, user.Id, supplierId)
	for _, imoneyResource := range imoneyResources {
		resources = append(resources, imoneyResource)
	}
	
	if couponUsage != nil {
		couponResource := this.parseCouponResource(couponUsage, productResources)
		if couponResource != nil {
			resources = append(resources, couponResource)
		}
	}
	
	return resources
}

func (this *ParseResourceService) parseProductResource(productsArray []interface{}) []*ProductResource {
	productResources := make([]*ProductResource, 0)
	for _, item := range productsArray {
		data := item.(map[string]interface{})

		productResource := NewProductResource(this.Ctx)

		poolProductId, _ := data["id"].(json.Number).Int64()
		productResource.PoolProductId = int(poolProductId)

		count, _ := data["count"].(json.Number).Int64()
		productResource.Count = int(count)

		price, _ := data["price"].(json.Number).Float64()
		productResource.Price = price

		productResource.Sku = data["sku"].(string)
		productResource.Ctx = this.Ctx

		productResources = append(productResources, productResource)
	}

	this.fillProductResources(productResources)
	return productResources
}

func (this *ParseResourceService) fillProductResources(productResources []*ProductResource) {
	poolProductIds := make([]int, 0)
	for _, productResource := range productResources {
		poolProductIds = append(poolProductIds, productResource.PoolProductId)
	}
	
	poolProducts := product.GetGlobalProductPool(this.Ctx).GetPoolProductsByIds(poolProductIds)
	product.NewFillPoolProductService(this.Ctx).Fill(poolProducts, eel.FillOption{
		"with_sku": true,
		"with_logistics": true,
	})
	
	//为productResource设置poolProduct
	id2product := make(map[int]*product.PoolProduct)
	for _, poolProduct := range poolProducts {
		id2product[poolProduct.Id] = poolProduct
	}
	
	for _, productResource := range productResources {
		if poolProduct, ok := id2product[productResource.PoolProductId]; ok {
			productResource.SetPoolProduct(poolProduct)
		}
	}
}

func (this *ParseResourceService) parseIMoneyResource(imoneysArray []interface{}, sourceUserId int, destUserId int) []*IMoneyResource {
	resources := make([]*IMoneyResource, 0)
	for _, item := range imoneysArray {
		data := item.(map[string]interface{})

		resource := NewIMoneyResource(this.Ctx)

		count, _ := data["count"].(json.Number).Float64()
		resource.Count = count
		resource.Code = data["code"].(string)
		resource.DestUserId = destUserId
		resource.Ctx = this.Ctx

		if _, ok := data["source_user_id"]; ok {
			id, _ := data["source_user_id"].(json.Number).Int64()
			sourceUserId = int(id)
		}
		resource.SourceUserId = sourceUserId


		resources = append(resources, resource)
	}

	return resources
}

func (this *ParseResourceService) parseCouponResource(couponUsage *CouponUsage, productResources []*ProductResource) *CouponResource {
	poolProductIds := make([]int, 0)
	for _, productResource := range productResources {
		poolProduct := productResource.GetPoolProduct()
		poolProductId := 0
		if poolProduct.IsSelfProduct() {
			poolProductId = poolProduct.Id
		} else {
			poolProductId = poolProduct.SourcePoolProductId
		}
		
		poolProductIds = append(poolProductIds, poolProductId)
	}
	
	resource := CouponResource{
		Code: couponUsage.Code,
		poolProductIds: poolProductIds,
	}
	resource.Ctx = this.Ctx
	
	return &resource
}

func init() {
}
