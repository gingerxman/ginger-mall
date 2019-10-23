package order

import (
	"context"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/product"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type OrderProductRepository struct {
	eel.ServiceBase
}

func NewOrderProductRepository(ctx context.Context) *OrderProductRepository {
	service := new(OrderProductRepository)
	service.Ctx = ctx
	return service
}

func (this *OrderProductRepository) GetOrderProducts(invoiceIds []int) []*OrderProduct {
	o := eel.GetOrmFromContext(this.Ctx)
	
	var models []*m_order.OrderHasProduct
	db := o.Model(&m_order.OrderHasProduct{}).Where("order_id__in", invoiceIds).Find(&models)
	err := db.Error
	
	if err != nil {
		eel.Logger.Error(err)
		return make([]*OrderProduct, 0)
	}
	
	orderProducts := make([]*OrderProduct, 0)
	productIds := make([]int, 0)
	for _, model := range models {
		product := &OrderProduct{}
		product.ProductId = model.ProductId
		product.Name = model.ProductName
		product.Sku = model.ProductSkuName
		product.SkuDisplayName = model.ProductSkuDisplayName
		product.Price = model.Price
		product.PurchaseCount = model.Count
		product.Thumbnail = model.Thumbnail
		product.Weight = model.Weight
		product.OrderId = model.OrderId
		
		orderProducts = append(orderProducts, product)
		productIds = append(productIds, model.ProductId)
	}
	
	//填充SupplierId
	products := product.NewProductRepository(this.Ctx).GetProductsByIds(productIds)
	id2product := make(map[int]*product.Product, 0)
	for _, product := range products {
		id2product[product.Id] = product
	}
	
	for _, orderProduct := range orderProducts {
		if product, ok := id2product[orderProduct.ProductId]; ok {
			orderProduct.SupplierId = product.CorpId
		}
	}
	
	return orderProducts
}

func init() {
}
