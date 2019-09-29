package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	"strconv"
	"strings"
	
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductFactory struct {
	eel.ServiceBase
}

func NewProductFactory(ctx context.Context) *ProductFactory {
	repository := new(ProductFactory)
	repository.Ctx = ctx
	return repository
}

//func (this *ProductFactory) addProductToCategories(product *Product, categoryIds []int) {
//	o := eel.GetOrmFromContext(this.Ctx)
//
//	if len(categoryIds) == 0 {
//		return
//	}
//
//	relationModels := make([]*m_product.CategoryHasProduct, 0)
//	//创建ProductHasProductCategory记录
//	for _, productCategoryId := range categoryIds {
//		relationModel := m_product.CategoryHasProduct{}
//		relationModel.ProductId = product.Id
//		relationModel.CategoryId = productCategoryId
//		relationModels = append(relationModels, &relationModel)
//	}
//
//	_, err := o.InsertMulti(len(relationModels), relationModels)
//	if err != nil {
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("product:set_category_fail", fmt.Sprintf("设置商品分类失败")))
//	}
//}

func (this *ProductFactory) addMediasToProduct(product *Product, mediaInfo *productMediaInfo) {
	if len(mediaInfo.Images) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	//创建ProductHasProductCategory记录
	models := make([]*m_product.ProductMedia, 0)
	for _, image := range mediaInfo.Images {
		model := m_product.ProductMedia{}
		model.ProductId = product.Id
		model.Type = m_product.PRODUCT_MEDIA_TYPE_IMAGE
		model.Url = image.Url
		models = append(models, &model)
	}
	
	//TODO: 替换成o.BatchInsert
	for _, model := range models {
		db := o.Create(model)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("product:set_media_fail", fmt.Sprintf("设置商品媒体资源失败")))
		}
	}
}

func (this *ProductFactory) addLogisticsInfoToProduct(product *Product, logisticsInfo *productLogisticsInfo) {
	if logisticsInfo == nil {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	//创建ProductHasProductCategory记录
	unifiedMoney, err := strconv.ParseFloat(logisticsInfo.UnifiedPostageMoney, 64)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("product:invalid_unified_postage_money", logisticsInfo.UnifiedPostageMoney))
	}
	model := m_product.ProductLogisticsInfo{
		ProductId: product.Id,
		PostageType: logisticsInfo.PostageType,
		UnifiedPostageMoney: unifiedMoney,
		LimitZoneType: logisticsInfo.LimitZoneType,
		LimitZoneId: logisticsInfo.LimitZoneId,
	}
	
	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:set_logistics_info_fail", fmt.Sprintf("设置商品物流信息失败")))
	}
}

func (this *ProductFactory) isNeedCreateProperty(propertyInfo *productSkuPropertyInfo) bool {
	return propertyInfo.PropertyId <= 0
}

func (this *ProductFactory) isNeedCreatePropertyValue(propertyInfo *productSkuPropertyInfo) bool {
	return propertyInfo.PropertyValueId <= 0
}

func (this *ProductFactory) EnsureSkuPropertyExists(corp business.ICorp, skuInfos []*productSkuInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	for _, skuInfo := range skuInfos {
		needRebuildSkuName := false
		for _, propertyInfo := range skuInfo.Properties {
			if this.isNeedCreateProperty(propertyInfo) {
				if o.Model(&m_product.ProductProperty{}).Where("name", propertyInfo.PropertyText).Exist() {
					var model m_product.ProductProperty
					db := o.Model(&m_product.ProductProperty{}).Where("name", propertyInfo.PropertyText).Take(&model)
					if db.Error != nil {
						eel.Logger.Error(db.Error)
					} else {
						propertyInfo.PropertyId = model.Id
					}
				} else {
					newProperty := NewProductProperty(this.Ctx, corp, propertyInfo.PropertyText)
					propertyInfo.PropertyId = newProperty.Id
				}
				
				needRebuildSkuName = true
			}
			
			if this.isNeedCreatePropertyValue(propertyInfo) {
				if o.Model(&m_product.ProductPropertyValue{}).Where("text", propertyInfo.PropertyValueText).Exist() {
					var model m_product.ProductPropertyValue
					db := o.Model(&m_product.ProductPropertyValue{}).Where("text", propertyInfo.PropertyValueText).Take(&model)
					if db.Error != nil {
						eel.Logger.Error(db.Error)
					} else {
						propertyInfo.PropertyValueId = model.Id
					}
				} else {
					newPropertyValue := NewProductPropertyValue(this.Ctx, propertyInfo.PropertyId, propertyInfo.PropertyValueText, "")
					propertyInfo.PropertyValueId = newPropertyValue.Id
				}
				needRebuildSkuName = true
			}
		}
		
		if needRebuildSkuName {
			skuInfo.RebuildName()
		}
	}
}

func (this *ProductFactory) addSkus(product *Product, skuInfos []*productSkuInfo) {
	if len(skuInfos) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	corp := account.GetCorpFromContext(this.Ctx)
	
	this.EnsureSkuPropertyExists(corp, skuInfos)
	
	isStandardSku := len(skuInfos) == 1 && skuInfos[0].Name == "standard"
	
	//创建<name, skuInfo>
	name2skuinfo := make(map[string]*productSkuInfo)
	
	//创建ProductSku记录
	models := make([]*m_product.ProductSku, 0)
	for _, skuInfo := range skuInfos {
		name2skuinfo[skuInfo.Name] = skuInfo
		model := m_product.ProductSku{}
		model.ProductId = product.Id
		model.CorpId = corp.GetId()
		model.Name = skuInfo.Name
		model.CostPrice = skuInfo.CostPrice
		model.Price = skuInfo.Price
		model.Stocks = skuInfo.Stocks
		model.Code = skuInfo.Code
		
		db := o.Create(&model)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("product:add_sku_fail_1", fmt.Sprintf("添加商品规格失败")))
		}
		models = append(models, &model)
	}
	
	
	if !isStandardSku {
		//创建ProductSkuHasPropertyValue记录
		relationModels := make([]*m_product.ProductSkuHasPropertyValue, 0)
		for _, model := range models {
			if skuInfo, ok := name2skuinfo[model.Name]; ok {
				for _, property := range skuInfo.Properties {
					relationModel := m_product.ProductSkuHasPropertyValue{}
					relationModel.PropertyId = property.PropertyId
					relationModel.PropertyValueId = property.PropertyValueId
					relationModel.SkuId = model.Id
					relationModels = append(relationModels, &relationModel)
				}
			}
		}
		
		//TODO: 替换成o.BatchInsert
		for _, model := range relationModels {
			db := o.Create(model)
			if db.Error != nil {
				eel.Logger.Error(db.Error)
				panic(eel.NewBusinessError("product:add_sku_fail_2", fmt.Sprintf("添加商品规格失败")))
			}
		}
		//_, err := o.InsertMulti(len(relationModels), relationModels)
		//if err != nil {
		//	eel.Logger.Error(err)
		//	panic(eel.NewBusinessError("product:add_sku_fail_2", fmt.Sprintf("添加商品规格失败")))
		//}
	}
}

func (this *ProductFactory) createProduct(baseInfo *productBaseInfo, mediaInfo *productMediaInfo) *Product {
	o := eel.GetOrmFromContext(this.Ctx)
	corp := account.GetCorpFromContext(this.Ctx)
	
	model := m_product.Product{}
	model.CorpId = corp.GetId()
	model.Type = baseInfo.Type
	model.Name = baseInfo.Name
	model.PromotionTitle = baseInfo.PromotionTitle
	model.Code = baseInfo.Code
	model.CategoryId = baseInfo.CategoryId
	model.DisplayIndex = 999999999
	model.Thumbnail = mediaInfo.Thumbnail
	model.IsDeleted = false
	
	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:create_fail", fmt.Sprintf("创建失败")))
	}
	
	//创建ProductDescription记录
	descriptionModel := m_product.ProductDescription{}
	descriptionModel.ProductId = model.Id
	descriptionModel.Introduction = ""
	descriptionModel.Detail = strings.TrimSpace(baseInfo.Detail)
	db = o.Create(&descriptionModel)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:create_description_fail", fmt.Sprintf("创建ProductDescription失败")))
	}
	
	return NewProductFromModel(this.Ctx, &model)
	
}

func (this *ProductFactory) CreateProduct(strBaseInfo string, strSkuInfos string, strMediaInfo string, strLogisticsInfo string) *PoolProduct {
	baseInfo := productBaseInfo{}
	err := json.Unmarshal([]byte(strBaseInfo), &baseInfo)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("product:parse_base_info_fail", "解析BaseInfo出错"))
	}

	//解析media info
	mediaInfo := productMediaInfo{}
	if strMediaInfo != "" {
		err = json.Unmarshal([]byte(strMediaInfo), &mediaInfo)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("product:parse_media_info_fail", "解析MediaInfo出错"))
		}
	}
	
	//解析logistics info
	var logisticsInfo *productLogisticsInfo
	if strLogisticsInfo != "" {
		logisticsInfo = &productLogisticsInfo{}
		err = json.Unmarshal([]byte(strLogisticsInfo), logisticsInfo)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("product:parse_logistics_info_fail", "解析LogisticsInfo出错"))
		}
	}
	
	//解析sku info
	skuInfos := make([]*productSkuInfo, 0)
	eel.Logger.Debug(strSkuInfos)
	err = json.Unmarshal([]byte(strSkuInfos), &skuInfos)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("product:parse_sku_info_fail", "解析SkuInfo出错"))
	}
	
	product := this.createProduct(&baseInfo, &mediaInfo)
	this.addMediasToProduct(product, &mediaInfo)
	this.addLogisticsInfoToProduct(product, logisticsInfo)
	this.addSkus(product, skuInfos)
	
	corp := account.GetCorpFromContext(this.Ctx)
	productPool := GetProductPoolForCorp(this.Ctx, corp)
	poolProduct := productPool.AddProduct(product, corp.GetId())
	
	return poolProduct
}

func init() {
}
