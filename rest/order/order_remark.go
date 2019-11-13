package order

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_order "github.com/gingerxman/ginger-mall/business/order"
)

type OrderRemark struct {
	eel.RestResource
}

func (this *OrderRemark) Resource() string {
	return "order.order_remark"
}

func (this *OrderRemark) GetParameters() map[string][]string {
	return map[string][]string{
		"POST": []string{"bid", "remark"},
	}
}

func (this *OrderRemark) Post(ctx *eel.Context) {
	//get order
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	bid := req.GetString("bid")
	remark := req.GetString("remark")
	corp := account.GetCorpFromContext(bCtx)
	invoice := b_order.NewOrderRepository(bCtx).GetInvoiceByBidForCorp(corp, bid)
	invoice.UpdateRemark(remark)
	
	ctx.Response.JSON(eel.Map{})
}
