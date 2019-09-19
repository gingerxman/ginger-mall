package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type FillProductService struct {
	eel.ServiceBase
	Corp business.ICorp
}

func NewFillProductService(ctx context.Context) *FillProductService {
	service := new(FillProductService)
	service.Ctx = ctx
	return service
}

func NewFillProductServiceForCorp(ctx context.Context, corp business.ICorp) *FillProductService{
	service := new(FillProductService)
	service.Ctx = ctx
	service.Corp = corp
	return service
}

func (this *FillProductService) Fill(products []*Product, option eel.FillOption) {
	if len(products) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, product := range products {
		ids = append(ids, product.Id)
	}

	if enableOption, ok := option["with_category"]; ok && enableOption {
		this.fillCategory(products, ids)
	}

	if enableOption, ok := option["with_label"]; ok && enableOption {
		this.fillProductLabel(products, ids)
	}

	if enableOption, ok := option["with_description"]; ok && enableOption {
		this.fillDescription(products, ids)
	}
	
	if enableOption, ok := option["with_logistics"]; ok && enableOption {
		this.fillLogistics(products, ids)
	}

	//if enableOption, ok := option["with_product_usable_imoney"]; ok && enableOption {
	//	this.fillProductUsableImoney(products, ids)
	//}

	if enableOption, ok := option["with_media"]; ok && enableOption {
		this.fillProductMedia(products, ids)
	}
	
	isEnableSkuProperty := false
	if option, ok := option["with_sku_property"]; ok {
		isEnableSkuProperty = option
	}
	if enableOption, ok := option["with_sku"]; ok && enableOption {
		this.fillSku(products, ids, isEnableSkuProperty)
	}
	return
}


func (this *FillProductService) fillCategory(products []*Product, ids []int) {
	categoryIds := make([]int, 0)
	for _, product := range products {
		categoryIds = append(categoryIds, product.CategoryId)
	}
	
	categories := NewProductCategoryRepository(this.Ctx).GetProductCategories(eel.Map{
		"id__in": categoryIds,
	})
	
	//构建<id, category>
	id2category := make(map[int]*ProductCategory)
	for _, category := range categories {
		id2category[category.Id] = category
	}
	
	for _, product := range products {
		if category, ok := id2category[product.CategoryId]; ok {
			product.Categories = append(product.Categories, category)
		}
	}
}



func (this *FillProductService) fillProductLabel(products []*Product, ids []int) {
	//构建<id, product>
	id2entity := make(map[int]*Product)
	for _, product := range products {
		id2entity[product.Id] = product
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	//从db中获取relation models
	var relationModels []*m_product.ProductHasLabel
	db := o.Model(&m_product.ProductHasLabel{}).Where("product_id__in", ids).Find(&relationModels)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	if len(relationModels) == 0 {
		return
	}
	
	//获取关联的id集合
	productLabelIds := make([]int, 0)
	for _, relationModel := range relationModels {
		productLabelIds = append(productLabelIds, relationModel.LabelId)
	}
	//从db中获取数据集合
	var models []*m_product.ProductLabel
	db = o.Model(&m_product.ProductLabel{}).Where("id__in", productLabelIds).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	//构建<id, model>
	id2model := make(map[int]*m_product.ProductLabel)
	for _, model := range models {
		id2model[model.Id] = model
	}
	
	//填充product的ProductLabels对象
	for _, relationModel := range relationModels {
		productId := relationModel.ProductId
		productLabelId := relationModel.LabelId
		
		if product, ok := id2entity[productId]; ok {
			if model, ok2 := id2model[productLabelId]; ok2 {
				product.Labels = append(product.Labels, NewProductLabelFromModel(this.Ctx, model))
			}
		}
	}
}

func (this *FillProductService) fillDescription(products []*Product, ids []int) {
	//从db中获取数据集合
	var models []*m_product.ProductDescription
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductDescription{}).Where("product_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}

	//构建<id, product>
	id2product := make(map[int]*Product)
	for _, product := range products {
		id2product[product.Id] = product
	}

	//填充product的ProductDescription对象
	for _, model := range models {
		if product, ok := id2product[model.ProductId]; ok {
			product.Description = NewProductDescriptionFromModel(this.Ctx, model)
		}
	}
}

func (this *FillProductService) fillLogistics(products []*Product, ids []int) {
	//从db中获取数据集合
	var models []*m_product.ProductLogisticsInfo
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductLogisticsInfo{}).Where("product_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	//构建<id, product>
	id2product := make(map[int]*Product)
	for _, product := range products {
		id2product[product.Id] = product
	}
	
	//填充product的ProductDescription对象
	for _, model := range models {
		if product, ok := id2product[model.ProductId]; ok {
			product.PostageType = model.PostageType
			product.UnifiedPostageMoney = model.UnifiedPostageMoney
			product.LimitZoneType = model.LimitZoneType
			product.LimitZoneId = model.LimitZoneId
		}
	}
}



//func (this *FillProductService) fillProductUsableImoney(products []*Product, ids []int) {
//	//获取关联的id集合
//	productUsableImoneyIds := make([]int, 0)
//	for _, product := range products {
//		productUsableImoneyIds = append(productUsableImoneyIds, product.ProductUsableImoneyId)
//	}
//
//	//从db中获取数据集合
//	var models []*m_product.ProductUsableImoney
//	o := eel.GetOrmFromContext(this.Ctx)
//	_, err := o.Model(&m_product.ProductUsableImoney{}).Where("id__in", productUsableImoneyIds).All(&models)
//	if err != nil {
//		eel.Logger.Error(err)
//		return
//	}
//
//	//构建<id, model>
//	id2model := make(map[int]*m_product.ProductUsableImoney)
//	for _, model := range models {
//		id2model[model.Id] = model
//	}
//
//	//填充product的ProductUsableImoney对象
//	for _, product := range products {
//		if model, ok := id2model[product.ProductUsableImoneyId]; ok {
//			product.ProductUsableImoney = NewProductUsableImoneyFromModel(this.Ctx, model)
//		}
//	}
//}



func (this *FillProductService) fillProductMedia(products []*Product, ids []int) {
	//从db中获取数据集合
	var models []*m_product.ProductMedia
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductMedia{}).Where("product_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	//构建<id, product>
	id2product := make(map[int]*Product)
	for _, product := range products {
		id2product[product.Id] = product
	}
	
	//填充product的ProductDescription对象
	for _, model := range models {
		if product, ok := id2product[model.ProductId]; ok {
			product.Medias = append(product.Medias, NewProductMediaFromModel(this.Ctx, model))
			
			if model.Type == m_product.PRODUCT_MEDIA_TYPE_IMAGE && product.Thumbnail == "" {
				product.Thumbnail = model.Url
			}
		}
	}
}



func (this *FillProductService) fillSku(products []*Product, ids []int, isEnableSkuProperty bool) {
	firstProduct := products[0]
	if len(firstProduct.Skus) > 0 {
		//已经完成过填充，再次进入，跳过填充
		return
	}
	
	NewProductSkuGenerator(this.Ctx).FillSkusForProducts(products, isEnableSkuProperty)
}

func init() {
}
