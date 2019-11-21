package product

import (
	"context"
	"errors"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)

type PoolProduct struct {
	eel.EntityBase
	Id int
	CorpId int
	ProductId int
	ProductType string
	SupplierId int
	Status string
	Type string
	SourcePoolProductId int
	DisplayIndex int
	SoldCount int
	SyncAt time.Time
	CreatedAt time.Time

	//foreign key
	Product *Product
}

//Update 更新对象
func (this *PoolProduct) Update(
	userId int,
	productId int,
	name string,
	supplierId int,
	status int,
	productType int,
	displayIndex int,
	syncAt string,
) error {
	var model m_product.PoolProduct
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"product_id": productId,
		"name": name,
		"supplier_id": supplierId,
		"status": status,
		"product_type": productType,
		"display_index": displayIndex,
		"sync_at": syncAt,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("pool_product:update_fail")
	}

	return nil
}

func (this *PoolProduct) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.PoolProduct{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *PoolProduct) Enable() {
	this.enable(true)
}

func (this *PoolProduct) Disable() {
	this.enable(false)
}

func (this *PoolProduct) RemoveLabel(label *ProductLabel) {
	product := NewProductRepository(this.Ctx).GetProduct(this.ProductId)
	product.RemoveLabel(label)
}

func (this *PoolProduct) SetLabels(labels []*ProductLabel) {
	product := NewProductRepository(this.Ctx).GetProduct(this.ProductId)
	product.SetLabels(labels)
}

func (this *PoolProduct) GetSku(skuName string) *ProductSku {
	return this.Product.GetSku(skuName)
}

func (this *PoolProduct) UseUnifiedPostage() bool {
	return this.Product.UseUnifiedPostage()
}

func (this *PoolProduct) GetUnifiedPostageMoney() int {
	return this.Product.UnifiedPostageMoney
}

func (this *PoolProduct) CanPurchase() bool {
	return this.Status == m_product.PPSTATUS2STR[m_product.PP_STATUS_ON]
}

func (this *PoolProduct) IsSelfProduct() bool {
	return this.Type == "create"
}

//工厂方法
//根据model构建对象
func NewPoolProductFromModel(ctx context.Context, model *m_product.PoolProduct) *PoolProduct {
	instance := new(PoolProduct)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.ProductId = model.ProductId
	instance.ProductType = model.ProductType
	instance.SupplierId = model.SupplierId
	instance.SoldCount = model.SoldCount
	instance.Status = m_product.PPSTATUS2STR[model.Status]
	instance.Type = m_product.PPTYPE2STR[model.Type]
	instance.SourcePoolProductId = model.SourcePoolProductId
	instance.DisplayIndex = model.DisplayIndex
	instance.SyncAt = model.SyncAt
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
