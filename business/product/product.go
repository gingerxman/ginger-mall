package product

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/mall/limit_zone"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)


type productBaseInfo struct {
	Name string `json:"name"` //商品名
	Type string `json:"type"` //商品类型
	Code string `json:"code"` //商品条码
	CategoryId int `json:"category_id"` //所属分类id
	PromotionTitle string `json:"promotion_title"` //促销标题
	Detail string `json:"detail"` //商品详情
}

type productImage struct {
	Url string `json:"url"`
}
type productMediaInfo struct {
	Images []productImage `json:"images"` //轮播图
	Thumbnail string `json:"thumbnail"` //缩略图
}

type productLogisticsInfo struct {
	PostageType string `json:"postage_type"` //运费类型
	UnifiedPostageMoney string `json:"unified_postage_money"` //统一运费金额
	LimitZoneType int `json:"limit_zone_type"` //仅售禁售类型
	LimitZoneId int `json:"limit_zone_id"` //仅售禁售模板id
}

type productSkuPropertyInfo struct {
	PropertyId int `json:"property_id"`
	PropertyValueId int `json:"property_value_id"`
}
type productSkuInfo struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	CostPrice float64 `json:"cost_price"`
	Stocks int `json:"stocks"`
	Code string `json:"code"`
	Properties []productSkuPropertyInfo `json:"properties"`
}

type Product struct {
	eel.EntityBase
	Id int
	CorpId int
	
	Type string
	CreateType string
	Name string
	PromotionTitle string
	Code string
	CategoryId int
	DisplayIndex int
	Thumbnail string
	IsDeleted bool
	IsEnabled bool
	CreatedAt time.Time
	
	//logistics info
	PostageType string
	UnifiedPostageMoney float64
	LimitZoneType int
	LimitZoneId int

	//foreign key
	Categories []*ProductCategory
	Labels []*ProductLabel
	Description *ProductDescription
	Medias []*ProductMedia
	ProductUsableImoneyId int //refer to product_usable_imoney
	ProductUsableImoney *ProductUsableImoney
	Skus []*ProductSku
}

func (this *Product) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.Product{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *Product) Enable() {
	this.enable(true)
}

func (this *Product) Disable() {
	this.enable(false)
}

//func (this *Product) HasStandardSku() bool {
//	if len(this.Skus) == 0 {
//		return true
//	}
//
//	if len(this.Skus) == 1 && this.Skus[0].IsStandardSku() {
//		return true
//	}
//
//	return false
//}

func (this *Product) UseUnifiedPostage() bool {
	return this.PostageType == "unified_postage_type"
}

func (this *Product) SetLabels(labels []*ProductLabel) {
	if len(labels) == 0 {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	var models []*m_product.ProductHasLabel
	
	db := o.Model(&m_product.ProductHasLabel{}).Where("product_id", this.Id).Find(&models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:set_label_fail_1", "获取标签列表失败"))
	}
	
	//构造<labelId, true>
	existedLabelIds := make(map[int]bool)
	for _, model := range models {
		existedLabelIds[model.LabelId] = true
	}
	
	//获取需要添加的label id集合
	needAddLabelIds := make([]int, 0)
	for _, label := range labels {
		if _, ok := existedLabelIds[label.Id]; !ok {
			needAddLabelIds = append(needAddLabelIds, label.Id)
		}
	}
	
	//创建数据
	if len(needAddLabelIds) == 0 {
		return
	}
	var newModels []interface{}
	for _, needAddLabelId := range needAddLabelIds {
		newModel := &m_product.ProductHasLabel{
			ProductId: this.Id,
			LabelId: needAddLabelId,
		}
		newModels = append(newModels, newModel)
	}
	
	_, err := o.BatchInsert(newModels)
	if err != nil {
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("product:set_label_fail_2", "创建标签失败"))
		}
	}
}

func (this *Product) RemoveLabel(label *ProductLabel) {
	if label == nil {
		return
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Where(eel.Map{
		"product_id": this.Id,
		"label_id": label.Id,
	}).Delete(&m_product.ProductHasLabel{})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product:remove_label_fail", fmt.Sprintf("删除标签%d失败", label.Id)))
	}
}

func (this *Product) GetSku(skuName string) *ProductSku {
	for _, sku := range this.Skus {
		if sku.Name == skuName {
			return sku
		}
	}
	
	return nil
}

func (this *Product) GetLimitZontTypeText() string {
	if this.LimitZoneType == m_product.LIMITZONE_TYPE_SALE {
		return "sale"
	} else if this.LimitZoneType == m_product.LIMITZONE_TYPE_FORBIDDEN {
		return "forbidden"
	} else {
		return "no_limit"
	}
}

func (this *Product) GetLimitZoneAreas() []*limit_zone.LimitProvince {
	if this.LimitZoneId > 0 {
		return limit_zone.NewLimitZoneRepository(this.Ctx).GetLimitZone(this.LimitZoneId).Provinces
	} else {
		return make([]*limit_zone.LimitProvince, 0)
	}
}

func (this *Product) HasStandardSku() bool {
	if len(this.Skus) == 0 {
		return true
	}
	
	if len(this.Skus) == 1 && this.Skus[0].IsStandardSku() {
		return true
	}
	
	return false
}

//根据model构建对象
func NewProductFromModel(ctx context.Context, model *m_product.Product) *Product {
	instance := new(Product)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Type = model.Type
	instance.Name = model.Name
	instance.PromotionTitle = model.PromotionTitle
	instance.Code = model.Code
	instance.CategoryId = model.CategoryId
	instance.DisplayIndex = model.DisplayIndex
	instance.Thumbnail = model.Thumbnail
	instance.IsDeleted = model.IsDeleted
	instance.CreatedAt = model.CreatedAt
	
	instance.Categories = make([]*ProductCategory, 0)
	instance.Skus = make([]*ProductSku, 0)
	instance.Medias = make([]*ProductMedia, 0)
	instance.Labels = make([]*ProductLabel, 0)

	return instance
}

func init() {
}
