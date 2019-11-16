package order

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/order"
)

type CorpInvoices struct {
	eel.RestResource
}

func (this *CorpInvoices) Resource() string {
	return "order.corp_invoices"
}

func (this *CorpInvoices) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?filters:json"},
	}
}


func (this *CorpInvoices) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	filters := req.GetOrmFilters()
	spew.Dump(filters)
	pageInfo := req.GetPageInfo()
	
	corp := account.GetCorpFromContext(bCtx)
	orders, nextPageInfo := order.NewOrderRepository(bCtx).GetPagedInvoicesForCorp(corp, filters, pageInfo, "-created_at")

	fillOptions := eel.Map{}
	fillOptions["with_invoice"] = map[string]interface{}{
		"with_products": true,
	}
	order.NewFillOrderService(bCtx).Fill(orders, fillOptions)

	rows := order.NewEncodeOrderService(bCtx).EncodeMany(orders)
	ctx.Response.JSON(eel.Map{
		"orders": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}

