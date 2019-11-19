package resource

import (
	"context"
	"errors"
	"fmt"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/product"
	b_product "github.com/gingerxman/ginger-mall/business/product"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"strconv"
	"strings"
)

type ProductResource struct {
	eel.EntityBase
	Resource
	
	PoolProductId int
	Sku string
	Count int
	SalesmanId int
	Price int
	poolProduct *product.PoolProduct
}

func (this *ProductResource) GetType() string {
	return RESOURCE_TYPE_PRODUCT
}

func (this *ProductResource) CanSplit() bool {
	return false
}

func (this *ProductResource) GetDeductionMoney(deductableMoney int) int {
	return 0
}

func (this *ProductResource) GetPrice() int {
	return this.poolProduct.GetSku(this.Sku).Price * this.Count
}

func (this *ProductResource) GetPostage() int {
	if this.poolProduct.UseUnifiedPostage() {
		return this.poolProduct.GetUnifiedPostageMoney()
	} else {
		return 0
	}
}

func (this *ProductResource) GetRawResourceObject() interface{} {
	return this
}

func (this *ProductResource) IsNeedLockWhenConsume() bool {
	return true
}

func (this *ProductResource) GetLockName() string {
	return fmt.Sprintf("%d", this.PoolProductId)
}

func (this *ProductResource) ToMap() map[string]interface{} {
	productInfo := make(map[string]interface{})
	productResourceInfo := make(map[string]interface{})
	
	poolProduct := this.poolProduct
	product := poolProduct.Product
	sku := poolProduct.GetSku(this.Sku)
	
	productInfo["id"] = product.Id
	productInfo["pool_product_id"] = poolProduct.Id
	productInfo["name"] = product.Name
	productInfo["thumbnail"] = product.Thumbnail
	productInfo["price"] = sku.Price
	productInfo["sku_name"] = this.Sku
	productInfo["sku_display_name"] = sku.Name
	
	productResourceInfo["type"] = this.GetType()
	productResourceInfo["count"] = this.Count
	productResourceInfo["total_price"] = this.Count * sku.Price
	productResourceInfo["price"] = sku.Price
	productResourceInfo["deduction_money"] = this.GetDeductionMoney(0)
	productResourceInfo["product"] = productInfo
	
	return productResourceInfo
}

func (this *ProductResource) SaveForOrder(order business.IOrder) error {
	poolProduct := this.poolProduct
	product := poolProduct.Product
	sku := poolProduct.GetSku(this.Sku)
	
	model := &m_order.OrderHasProduct{}
	model.OrderId = order.GetId()
	model.ProductId = product.Id
	model.ProductName = product.Name
	model.Thumbnail = product.Thumbnail
	model.Count = this.Count
	model.Price = sku.Price
	model.ProductSkuName = this.Sku
	
	//获得sku display name
	if this.Sku == "standard" {
		model.ProductSkuDisplayName = "standard"
	} else {
		names := make([]string, 0)
		propertyRepository := b_product.NewProductPropertyRepository(this.Ctx)
		items := strings.Split(sku.Name, "_")
		for _, item := range items {
			valueId, _ := strconv.Atoi(strings.Split(item, ":")[1])
			propertyValue := propertyRepository.GetProductPropertyValue(valueId)
			names = append(names, propertyValue.Text)
		}
		model.ProductSkuDisplayName = strings.Join(names, " ")
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Create(model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_resource:save_fail", "保存product resource失败"))
	}
	
	return nil
}

func (this *ProductResource) IsValid() error {
	if this.Count < 0 {
		eel.Logger.Error(fmt.Sprintf("购买数量错误(%d)", this.Count))
		return errors.New("invalid_purchase_count")
	}
	
	if !this.poolProduct.CanPurchase() {
		eel.Logger.Error(fmt.Sprintf("商品已下架(%d)", this.Count))
		return errors.New("product_off_shelve")
	}
	
	sku := this.poolProduct.GetSku(this.Sku)
	if sku == nil {
		eel.Logger.Error(fmt.Sprintf("错误的sku(%s)", this.Sku))
		return errors.New("invalid_sku")
	}
	
	if !sku.CanAffordStock(this.Count) {
		eel.Logger.Error(fmt.Sprintf("库存不足(%d)", this.Count))
		return errors.New("not_enough_stocks")
	}
	
	if sku.Price != this.Price {
		eel.Logger.Error(fmt.Sprintf("价格发生变动(%s)", this.Price))
		return errors.New("price_change")
	}
	
	return nil
}

func (this *ProductResource) IsAllocated() bool{
	return this.Resource.IsAllocated()
}

func (this *ProductResource) SetAllocated(){
	this.Resource.SetAllocated()
}

func (this *ProductResource) ResetAllocation(){
	this.Resource.ResetAllocation()
}

func (this *ProductResource) GetPoolProduct() *product.PoolProduct {
	if (this.poolProduct == nil) {
		this.poolProduct = product.GetGlobalProductPool(this.Ctx).GetPoolProduct(this.PoolProductId)
	}
	
	return this.poolProduct
}

func (this *ProductResource) SetPoolProduct(poolProduct *product.PoolProduct) {
	this.poolProduct = poolProduct;
}

func NewProductResource(ctx context.Context) *ProductResource {
	instance := &ProductResource{}
	instance.Ctx = ctx
	
	return instance
}

func init() {
}
