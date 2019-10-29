package shopping_cart

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/shopping_cart"
)

type ProductCount struct {
	eel.RestResource
}

func (this *ProductCount) Resource() string {
	return "mall.shopping_cart_product_count"
}

func (this *ProductCount) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *ProductCount) Get(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	if corp == nil {
		ctx.Response.Error("mall.shopping_cart_product_count:invalid_corp", "")
		return
	}

	user := account.GetUserFromContext(bCtx)
	shoppingCart := shopping_cart.NewShoppingCartRepository(bCtx).GetShoppingCartForUserInCorp(user, corp)
	count := shoppingCart.GetProductCount()
	
	ctx.Response.JSON(eel.Map{
		"count": count,
	})
}
