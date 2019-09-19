package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)

type ProductShelf struct {
	eel.RepositoryBase
	Type string
	Corp business.ICorp
}


func (this *ProductShelf) AddProducts(poolProducts []*PoolProduct) {
	if len(poolProducts) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, poolProduct := range poolProducts {
		ids = append(ids, poolProduct.Id)
	}
	
	status := m_product.PP_STATUS_ON
	if this.Type == "for_sale" {
		status = m_product.PP_STATUS_OFF
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.PoolProduct{}).Where(eel.Map{
		"id__in": ids,
		"status__gt": m_product.PP_STATUS_DELETE,
	}).Update(gorm.Params{
		"status": status,
		"display_index": NEW_PRODUCT_DISPLAY_INDEX,
		"sync_at": time.Now(),
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_shelf:add_product_fail", "向货架添加商品失败"))
	}
}

func (this *ProductShelf) MoveProducts(poolProducts []*PoolProduct) {
	if len(poolProducts) > 0 {
		this.AddProducts(poolProducts)
	}
}

func (this *ProductShelf) AddProduct(poolProduct *PoolProduct) {
	this.AddProducts([]*PoolProduct{poolProduct})
}

func (this *ProductShelf) GetPagedProducts(filters eel.Map, page *eel.PageInfo) ([]*PoolProduct, eel.INextPageInfo) {
	productPool := GetProductPoolForCorp(this.Ctx, this.Corp)
	
	if this.Type == "in_sale" {
		filters["status"] = m_product.PP_STATUS_ON
	} else if this.Type == "for_sale" {
		filters["status"] = m_product.PP_STATUS_OFF
	}
	return productPool.GetPagedPoolProducts(filters, page, "-id")
}

func GetInSaleProductShelfForCorp(ctx context.Context, corp business.ICorp) *ProductShelf {
	instance := new(ProductShelf)
	instance.Ctx = ctx
	instance.Type = "in_sale"
	instance.Corp = corp
	return instance
}

func GetForSaleProductShelfForCorp(ctx context.Context, corp business.ICorp) *ProductShelf {
	instance := new(ProductShelf)
	instance.Ctx = ctx
	instance.Type = "for_sale"
	instance.Corp = corp
	return instance
}

func init() {
}
