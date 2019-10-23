package shopping_cart

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type ShoppingCartRepository struct {
	eel.RepositoryBase
}

func NewShoppingCartRepository(ctx context.Context) *ShoppingCartRepository {
	repository := new(ShoppingCartRepository)
	repository.Ctx = ctx
	return repository
}

//GetShipInfoInCorp 根据id和user获得ShipInfo对象
func (this *ShoppingCartRepository) GetShoppingCartForUserInCorp(user business.IUser, corp business.ICorp) *ShoppingCart {
	shoppingCart := ShoppingCart{
		UserId: user.GetId(),
		CorpId: corp.GetId(),
	}
	shoppingCart.Ctx = this.Ctx
	
	return &shoppingCart
}


func init() {
}
