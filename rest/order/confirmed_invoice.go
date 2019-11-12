package order

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_order "github.com/gingerxman/ginger-mall/business/order"
)

type ConfirmedInvoice struct {
	eel.RestResource
}

func (this *ConfirmedInvoice) Resource() string {
	return "order.confirmed_invoice"
}

func (this *ConfirmedInvoice) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{
			"bid",
		},
	}
}

func (this *ConfirmedInvoice) Put(ctx *eel.Context) {
	req := ctx.Request
	bid := req.GetString("bid")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	invoice := b_order.NewOrderRepository(bCtx).GetInvoiceByBidForCorp(corp, bid)
	
	if invoice == nil {
		ctx.Response.Error("confirmed_invoice:invalid_invoice", fmt.Sprintf("invalid bid(%s)", bid))
	} else {
		invoice.Confirm()
		
		ctx.Response.JSON(eel.Map{})
	}
}
