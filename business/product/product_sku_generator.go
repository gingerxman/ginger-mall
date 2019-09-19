package product

import (
	"context"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"sort"
)

type ProductSkuGenerator struct {
	eel.ServiceBase
}

func NewProductSkuGenerator(ctx context.Context) *ProductSkuGenerator {
	service := new(ProductSkuGenerator)
	service.Ctx = ctx
	return service
}

func (this *ProductSkuGenerator) FillSkusForProducts(products []*Product, isEnableSkuProperty bool) {
	productIds := make([]int, 0)
	for _, product := range products {
		productIds = append(productIds, product.Id)
	}
	
	id2property, id2value := this.getAllProperty(products, isEnableSkuProperty)
	
	//获取所有skus
	var models []*m_product.ProductSku
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductSku{}).Where(eel.Map{
		"product_id__in": productIds,
		"is_deleted": false,
	}).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
	if len(models) == 0 {
		return
	}
	
	pid2skus := make(map[int][]*ProductSku)
	for _, model := range models {
		sku := NewProductSkuFromModel(this.Ctx, model)
		sku.FillProperties(id2property, id2value)
		var skus []*ProductSku
		if tmpSkus, ok := pid2skus[model.ProductId]; ok {
			skus = tmpSkus
		} else {
			skus = make([]*ProductSku, 0)
		}
		skus = append(skus, sku)
		pid2skus[model.ProductId] = skus
	}
	
	//填充poolProducts
	for _, product := range products {
		product.Skus = pid2skus[product.Id]
	}
}

//getAllProperty 获取系统中所有的商品规格属性信息
func (this *ProductSkuGenerator) getAllProperty(products []*Product, isEnableSkuProperty bool) (map[int]*ProductProperty, map[int]*ProductPropertyValue) {
	corpIds := make([]int, 0)
	for _, product := range products {
		corpIds = append(corpIds, product.CorpId)
	}
	
	properties := NewProductPropertyRepository(this.Ctx).GetProductPropertiesInCorps(corpIds)
	NewFillProductPropertyService(this.Ctx).Fill(properties, eel.FillOption{
		"with_value": true,
	})
	id2property := make(map[int]*ProductProperty)
	id2value := make(map[int]*ProductPropertyValue)
	for _, property := range properties {
		id2property[property.Id] = property
		for _, value := range property.Values {
			id2value[value.Id] = value
		}
	}
	
	return id2property, id2value
}

func (this *ProductSkuGenerator) fillUsedProperty(product *Product, id2property map[int]*ProductProperty, id2value map[int]*ProductPropertyValue) {
	/*
	填充商品中使用了的商品规格属性的信息
	
	从models中构建used_system_model_properties，
	假如商品有以下两个规格
	1. {property:'颜色', value:'红色'}, {property:'尺寸', value:'M'}
	2. {property:'颜色', value:'黄色'}, {property:'尺寸', value:'M'}
	
	则合并后的used_system_model_properties为:
	[{
		property: '颜色',
		values: ['红色', '黄色']
	}, {
		property: '尺寸',
		values: ['M']
	}]
	*/
	if product.HasStandardSku() {
		return
	}
	
	//收集排序后的propertyIds, valueIds
	propertyIds := make([]int, 0)
	valueIds := make([]int, 0)
	propertyid2bool := make(map[int]bool)
	valueid2bool := make(map[int]bool)
	for _, sku := range product.Skus {
		for _, propertyValue := range sku.PropertyValues {
			propertyId := propertyValue.PropertyId
			if _, ok := propertyid2bool[propertyId]; !ok {
				propertyIds = append(propertyIds, propertyId)
				propertyid2bool[propertyId] = true
			}
			
			valueId := propertyValue.Id
			if _, ok := valueid2bool[valueId]; !ok {
				valueIds = append(valueIds, valueId)
				valueid2bool[valueId] = true
			}
		}
	}
	
	sort.Sort(sort.IntSlice(propertyIds))
	sort.Sort(sort.IntSlice(valueIds))
}


func init() {
}
