package order

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_order "github.com/gingerxman/ginger-mall/business/order"
)

type CanceledInvoice struct {
	eel.RestResource
}

func (this *CanceledInvoice) Resource() string {
	return "order.canceled_invoice"
}

func (this *CanceledInvoice) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{
			"bid",
			"reason",
		},
	}
}

func (this *CanceledInvoice) Put(ctx *eel.Context) {
	req := ctx.Request
	bid := req.GetString("bid")
	reason := req.GetString("reason")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	invoice := b_order.NewOrderRepository(bCtx).GetInvoiceByBidForCorp(corp, bid)
	
	if invoice == nil {
		ctx.Response.Error("finished_invoice:invalid_invoice", fmt.Sprintf("invalid bid(%s)", bid))
	} else {
		invoice.Cancel(reason)
		
		ctx.Response.JSON(eel.Map{})
	}
}
