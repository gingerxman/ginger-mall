package order

import (
	"fmt"
	"github.com/gingerxman/eel"
	b_order "github.com/gingerxman/ginger-mall/business/order"
)

type CanceledOrder struct {
	eel.RestResource
}

func (this *CanceledOrder) Resource() string {
	return "order.canceled_order"
}

func (this *CanceledOrder) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{
			"bid",
			"reason",
		},
	}
}

func (this *CanceledOrder) Put(ctx *eel.Context) {
	req := ctx.Request
	bid := req.GetString("bid")
	reason := req.GetString("reason")
	
	bCtx := ctx.GetBusinessContext()
	order := b_order.NewOrderRepository(bCtx).GetOrderByBid(bid)
	
	if order == nil {
		ctx.Response.Error("canceled_order:invalid_order", fmt.Sprintf("invalid bid(%s)", bid))
	} else {
		order.Cancel(reason)
		ctx.Response.JSON(eel.Map{})
	}
}
