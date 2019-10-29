package shopping_cart

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/shopping_cart"
)

type ShoppingCart struct {
	eel.RestResource
}

func (this *ShoppingCart) Resource() string {
	return "mall.shopping_cart"
}

func (this *ShoppingCart) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *ShoppingCart) Get(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	if corp == nil {
		ctx.Response.Error("shopping_cart:invalid_corp", fmt.Sprintf("%d", corp.Id))
		return
	}
	
	shoppingCart := shopping_cart.NewShoppingCartRepository(bCtx).GetShoppingCartForUserInCorp(user, corp)
	
	if shoppingCart == nil {
		ctx.Response.Error("shopping_cart:invalid_shopping_cart", fmt.Sprintf("user(%d), corp(%d)", user.GetId(), corp.GetId()))
	} else {
		encodeService := shopping_cart.NewEncodeShoppingCartService(bCtx)
		respData := encodeService.Encode(shoppingCart)
		
		ctx.Response.JSON(respData)
	}
}
