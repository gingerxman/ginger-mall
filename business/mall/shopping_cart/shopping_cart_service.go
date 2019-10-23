package shopping_cart

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
)

type ShoppingCartService struct {
	eel.ServiceBase
}

func NewShoppingCartService(ctx context.Context) *ShoppingCartService {
	service := new(ShoppingCartService)
	service.Ctx = ctx
	return service
}

func (this *ShoppingCartService) DeleteShoppingCartItems(ids[] int) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.ShoppingCartItem{}).Where("id__in", ids).Delete(&m_mall.ShoppingCartItem{})
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
	}
}


func init() {
}
