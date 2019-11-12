package order

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_order "github.com/gingerxman/ginger-mall/business/order"
)

type FinishedInvoice struct {
	eel.RestResource
}

func (this *FinishedInvoice) Resource() string {
	return "order.finished_invoice"
}

func (this *FinishedInvoice) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{
			"bid",
		},
	}
}

func (this *FinishedInvoice) Put(ctx *eel.Context) {
	req := ctx.Request
	bid := req.GetString("bid")
	
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	invoice := b_order.NewOrderRepository(bCtx).GetInvoiceByBidForUser(user, bid)
	
	if invoice == nil {
		ctx.Response.Error("finished_invoice:invalid_invoice", fmt.Sprintf("invalid bid(%s)", bid))
	} else {
		invoice.Finish()
		
		ctx.Response.JSON(eel.Map{})
	}
}
