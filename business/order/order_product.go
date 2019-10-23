package order

import (
	"context"
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)


type OrderProduct struct {
	eel.EntityBase
	
	Name string
	Thumbnail string
	Sku string
	SkuDisplayName string
	PurchaseCount int
	Price float64
	Weight float64
	
	OrderId int
	SupplierId int
	ProductId int
}

func NewOrderProductFromModel(ctx context.Context, model *m_order.OrderHasProduct) *OrderProduct {
	return nil
}

func init() {
}
