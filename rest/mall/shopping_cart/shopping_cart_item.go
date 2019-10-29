package shopping_cart

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/shopping_cart"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ShoppingCartItem struct {
	eel.RestResource
}

func (this *ShoppingCartItem) Resource() string {
	return "mall.shopping_cart_item"
}

func (this *ShoppingCartItem) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"pool_product_id:int", "sku_name", "count:int"},
		"DELETE": []string{"id:int"},
	}
}

func (this *ShoppingCartItem) Put(ctx *eel.Context) {
	req := ctx.Request
	poolProductId, _ := req.GetInt("pool_product_id")
	skuName := req.GetString("sku_name")
	count, _ := req.GetInt("count")

	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	productPool := product.GetProductPoolForCorp(bCtx, corp)
	poolProduct := productPool.GetPoolProduct(poolProductId)
	
	shoppingCart := shopping_cart.NewShoppingCartRepository(bCtx).GetShoppingCartForUserInCorp(user, corp)
	err := shoppingCart.AddProduct(poolProduct, skuName, count)
	
	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("shopping_cart_item:create_fail", err.Error())
	} else {
		ctx.Response.JSON(eel.Map{
			"count": shoppingCart.GetProductCount(),
		})
	}
}

func (this *ShoppingCartItem) Delete(ctx *eel.Context) {
	req := ctx.Request
	
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	corp := account.GetCorpFromContext(bCtx)

	id, _ := req.GetInt("id")
	shoppingCart := shopping_cart.NewShoppingCartRepository(bCtx).GetShoppingCartForUserInCorp(user, corp)
	err := shoppingCart.DeleteItem(id)
	
	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("shopping_cart_item:delete_fail", err.Error())
	} else {
		ctx.Response.JSON(eel.Map{})
	}
}

