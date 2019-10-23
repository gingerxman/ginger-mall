package shopping_cart

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/product"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
)

type ShoppingCartProduct struct {
	eel.EntityBase
	Id int
	UserId int
	PoolProductId int
	ProductSkuName string
	ProductSkuDisplayName string
	Count int
	
	PoolProduct *product.PoolProduct
	
	_isFillValidity bool
	_isValid bool
}

func (this *ShoppingCartProduct) isValid() bool {
	if !this.PoolProduct.CanPurchase() {
		return false
	}
	
	sku := this.PoolProduct.GetSku(this.ProductSkuName)
	if sku == nil || sku.IsDeleted || !sku.HasStocks() {
		return false
	}
	
	if this.PoolProduct.Product.IsDeleted {
		return false
	}
	
	return true
}

// IsValid 判断购物车商品是否可以购买
func (this *ShoppingCartProduct) IsValid() bool {
	if !this._isFillValidity {
		this._isValid = this.isValid()
	}
	
	return this._isValid
}

//根据model构建对象
func NewShoppingCartProductFromModel(ctx context.Context, model *m_mall.ShoppingCartItem) *ShoppingCartProduct {
	instance := new(ShoppingCartProduct)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.UserId = model.UserId
	instance.PoolProductId = model.ProductId
	instance.ProductSkuName = model.ProductSkuName
	instance.ProductSkuDisplayName = model.ProductSkuDisplayName
	instance.Count = model.Count

	return instance
}

func init() {
}
