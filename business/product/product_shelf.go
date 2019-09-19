package product

import (
	"context"
	"github.com/kfchen81/beego/orm"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/vanilla"
	"gpeanut/business"
	m_product "gpeanut/models/product"
	"time"
)

type ProductShelf struct {
	vanilla.RepositoryBase
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
	
	o := vanilla.GetOrmFromContext(this.Ctx)
	_, err := o.QueryTable(&m_product.PoolProduct{}).Filter(vanilla.Map{
		"id__in": ids,
		"status__gt": m_product.PP_STATUS_DELETE,
	}).Update(orm.Params{
		"status": status,
		"display_index": NEW_PRODUCT_DISPLAY_INDEX,
		"sync_at": time.Now(),
	})
	
	if err != nil {
		beego.Error(err)
		panic(vanilla.NewBusinessError("product_shelf:add_product_fail", "向货架添加商品失败"))
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

func (this *ProductShelf) GetPagedProducts(filters vanilla.Map, page *vanilla.PageInfo) ([]*PoolProduct, vanilla.INextPageInfo) {
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
