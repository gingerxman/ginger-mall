package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"strconv"
	"strings"
)

type UpdateProductService struct {
	eel.ServiceBase
}

func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	service := new(UpdateProductService)
	service.Ctx = ctx
	return service
}

func (this *UpdateProductService) Update(productId int, strBaseInfo string, strSkuInfos string, strMediaInfo string, strLogisticsInfo string) {
	baseInfo := productBaseInfo{}
	err := json.Unmarshal([]byte(strBaseInfo), &baseInfo)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("update_product_service:parse_base_info_fail", "解析BaseInfo出错"))
	}
	
	//解析sku info
	skuInfos := make([]*productSkuInfo, 0)
	err = json.Unmarshal([]byte(strSkuInfos), &skuInfos)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("update_product_service:parse_sku_info_fail", "解析SkuInfo出错"))
	}
	
	//解析media info
	mediaInfo := productMediaInfo{}
	if strMediaInfo != "" {
		err = json.Unmarshal([]byte(strMediaInfo), &mediaInfo)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("update_product_service:parse_media_info_fail", "解析MediaInfo出错"))
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

	this.updateProduct(productId, &baseInfo, &mediaInfo)
	this.updateMedias(productId, &mediaInfo)
	this.updateSkus(productId, skuInfos)
	this.updateLogisticsInfo(productId, logisticsInfo)
}

func (this *UpdateProductService) updateProduct(productId int, baseInfo *productBaseInfo, mediaInfo *productMediaInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.Product{}).Where("id", productId).Update(gorm.Params{
		"name": strings.TrimSpace(baseInfo.Name),
		"promotion_title": strings.TrimSpace(baseInfo.PromotionTitle),
		"code": baseInfo.Code,
		"category_id": baseInfo.CategoryId,
		"thumbnail": strings.TrimSpace(mediaInfo.Thumbnail),
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:update_fail", "更新商品基本信息失败"))
	}
	
	db = o.Model(&m_product.ProductDescription{}).Where("product_id", productId).Update(gorm.Params{
		"detail": strings.TrimSpace(baseInfo.Detail),
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:update_fail", "更新商品详情失败"))
	}
}

func (this *UpdateProductService) addMediasToProduct(productId int, mediaInfo *productMediaInfo) {
	if len(mediaInfo.Images) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	//创建ProductHasProductCategory记录
	models := make([]*m_product.ProductMedia, 0)
	for _, image := range mediaInfo.Images {
		model := m_product.ProductMedia{}
		model.ProductId = productId
		model.Type = m_product.PRODUCT_MEDIA_TYPE_IMAGE
		model.Url = image.Url
		models = append(models, &model)
	}
	
	//TODO: 替换为o.BatchInsert方案
	for _, model := range models {
		db := o.Create(model)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("product:set_media_fail", fmt.Sprintf("设置商品媒体资源失败")))
		}
	}
	//_, err := o.InsertMulti(len(models), models)
	//if err != nil {
	//	eel.Logger.Error(err)
	//	panic(eel.NewBusinessError("product:set_media_fail", fmt.Sprintf("设置商品媒体资源失败")))
	//}
}

func (this *UpdateProductService) updateMedias(productId int, mediaInfo *productMediaInfo) {
	if mediaInfo == nil || len(mediaInfo.Images) == 0 {
		return
	}

	o := eel.GetOrmFromContext(this.Ctx)

	//删除老的media
	db := o.Where("product_id", productId).Delete(&m_product.ProductMedia{})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:update_fail", "删除商品媒体信息失败"))
	}
	
	this.addMediasToProduct(productId, mediaInfo)
}

//func (this *UpdateProductService) addProductToCategories(productId int, categoryIds []int) {
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
//		relationModel.ProductId = productId
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

//func (this *UpdateProductService) updateCategories(productId int, categoryIds []int) {
//	o := eel.GetOrmFromContext(this.Ctx)
//
//	//删除老的media
//	_, err := o.Model(&m_product.CategoryHasProduct{}).Where("product_id", productId).Delete()
//
//	if err != nil {
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("product:update_fail", "删除商品分类失败"))
//	}
//
//	if len(categoryIds) > 0 {
//		this.addProductToCategories(productId, categoryIds)
//	}
//}

func (this *UpdateProductService) updateSkus(productId int, skuInfos []*productSkuInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	var models []*m_product.ProductSku
	db := o.Model(&m_product.ProductSku{}).Where(eel.Map{
		"product_id": productId,
		"is_deleted": false,
	}).Find(&models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:update_fail", fmt.Sprintf("获取商品规格失败")))
	}
	
	//获取<new_sku_id, true>
	newSkuIds := make(map[int]bool)
	for _, skuInfo := range skuInfos {
		newSkuIds[skuInfo.Id] = true
	}
	
	//获取<existed_sku_id, true>
	existedSkuIds := make(map[int]bool)
	for _, model := range models {
		existedSkuIds[model.Id] = true
	}
	
	//获取需要增加和需要更新的sku
	needAddSkus := make([]*productSkuInfo, 0)
	needUpdateSkus := make([]*productSkuInfo, 0)
	for _, skuInfo := range skuInfos {
		if isExist, ok := existedSkuIds[skuInfo.Id]; ok && isExist {
			needUpdateSkus = append(needUpdateSkus, skuInfo)
		} else {
			needAddSkus = append(needAddSkus, skuInfo)
		}
	}
	
	//获取需要删除的sku
	needDeleteSkuIds := make([]int, 0)
	for _, model := range models {
		if _, exist := newSkuIds[model.Id]; !exist {
			needDeleteSkuIds = append(needDeleteSkuIds, model.Id)
		}
	}
	
	this.__deleteExistedSkus(productId, needDeleteSkuIds)
	this.__updateExistedSkus(productId, needUpdateSkus)
	this.__addSkus(productId, needAddSkus)
}

func (this *UpdateProductService) __updateExistedSkus(productId int, newSkuInfos []*productSkuInfo) {
	if len(newSkuInfos) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	for _, skuInfo := range newSkuInfos {
		db := o.Model(m_product.ProductSku{}).Where("id", skuInfo.Id).Update(gorm.Params{
			"price": skuInfo.Price,
			"stocks": skuInfo.Stocks,
			//"user_code": skuInfo.SkuCode,
			//"sku_code": skuInfo.SkuCode,
			"cost_price": skuInfo.CostPrice,
		})
		
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("product:update_sku_fail", "删除商品规格失败"))
		}
	}
}

func (this *UpdateProductService) __deleteExistedSkus(productId int, needDeleteIds []int) {
	if len(needDeleteIds) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductSku{}).Where("id__in", needDeleteIds).Update(gorm.Params{
		"is_deleted": true,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:delete_sku_fail", "删除商品规格失败"))
	}
}

func (this *UpdateProductService) __addSkus(productId int, skuInfos []*productSkuInfo) {
	if len(skuInfos) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	corp := account.GetCorpFromContext(this.Ctx)
	
	isStandardSku := len(skuInfos) == 1 && skuInfos[0].Name == "standard"
	
	//创建<name, skuInfo>
	name2skuinfo := make(map[string]*productSkuInfo)
	
	//创建ProductSku记录
	models := make([]*m_product.ProductSku, 0)
	for _, skuInfo := range skuInfos {
		name2skuinfo[skuInfo.Name] = skuInfo
		model := m_product.ProductSku{}
		model.ProductId = productId
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
			relationModel := m_product.ProductSkuHasPropertyValue{}
			if skuInfo, ok := name2skuinfo[model.Name]; ok {
				for _, property := range skuInfo.Properties {
					relationModel.PropertyId = property.PropertyId
					relationModel.PropertyValueId = property.PropertyValueId
					relationModel.SkuId = model.Id
					relationModels = append(relationModels, &relationModel)
				}
			}
		}
		
		//TODO: 替换为o.BatchInsert方案
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
func (this *UpdateProductService) updateLogisticsInfo(productId int, logisticsInfo *productLogisticsInfo) {
	if logisticsInfo == nil {
		return
	}

	o := eel.GetOrmFromContext(this.Ctx)

	unifiedMoney, err := strconv.ParseFloat(logisticsInfo.UnifiedPostageMoney, 64)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("product:invalid_unified_postage_money", logisticsInfo.UnifiedPostageMoney))
	}
	db := o.Model(&m_product.ProductLogisticsInfo{}).Where("product_id", productId).Update(gorm.Params{
		"postage_type": logisticsInfo.PostageType,
		"unified_postage_money": unifiedMoney,
		"limit_zone_type": logisticsInfo.LimitZoneType,
		"limit_zone_id": logisticsInfo.LimitZoneId,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:update_fail", "更新商品基本信息失败"))
	}
}

func init() {
}
