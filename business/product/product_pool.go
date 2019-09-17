package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type ProductPool struct {
	eel.RepositoryBase
	Corp business.ICorp
}

func (this *ProductPool) GetPoolProducts(filters eel.Map, orderExprs string) []*PoolProduct {
	return make([]*PoolProduct, 0)
	//o := vanilla.GetOrmFromContext(this.Ctx)
	//qs := o.QueryTable(&m_product.PoolProduct{})
	//
	//var models []*m_product.PoolProduct
	//if this.Corp != nil && this.Corp.IsValid() {
	//	filters["corp_id"] = this.Corp.GetId()
	//}
	//if len(filters) > 0 {
	//	qs = qs.Filter(filters)
	//}
	//if len(orderExprs) > 0 {
	//	qs = qs.Order(orderExprs...)
	//}
	//_, err := qs.All(&models)
	//if err != nil {
	//	beego.Error(err)
	//	return nil
	//}
	//
	//poolProducts := make([]*PoolProduct, 0)
	//for _, model := range models {
	//	poolProducts = append(poolProducts, NewPoolProductFromModel(this.Ctx, model))
	//}
	//return poolProducts
}

func GetProductPoolForCorp(ctx context.Context, corp business.ICorp) *ProductPool {
	instance := new(ProductPool)
	instance.Ctx = ctx
	instance.Corp = corp
	return instance
}

func GetGlobalProductPool(ctx context.Context) *ProductPool {
	instance := new(ProductPool)
	instance.Ctx = ctx
	//instance.Corp = account.NewInvalidCorp(ctx)
	return instance
}

func init() {
}
