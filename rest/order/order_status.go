package order

import (
	"github.com/gingerxman/eel"
	b_order "github.com/gingerxman/ginger-mall/business/order"
)

type OrderStatus struct {
	eel.RestResource
}

func (this *OrderStatus) Resource() string {
	return "order.order_status"
}

func (this *OrderStatus) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"bid"},
	}
}

func (this *OrderStatus) Get(ctx *eel.Context) {
	req := ctx.Request
	bid := req.GetString("bid")
	
	order := b_order.NewOrderRepository(ctx.GetBusinessContext()).GetOrderByBid(bid)
	
	ctx.Response.JSON(eel.Map{
		"status": order.GetStatusText(),
	})
}
