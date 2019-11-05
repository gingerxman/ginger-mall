package order

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/order"
	"github.com/davecgh/go-spew/spew"
)

type UserOrders struct {
	eel.RestResource
}

func (this *UserOrders) Resource() string {
	return "order.user_orders"
}

func (this *UserOrders) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?filters:json"},
	}
}


func (this *UserOrders) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	filters := req.GetOrmFilters()
	pageInfo := req.GetPageInfo().Desc()
	spew.Dump(pageInfo)

	corp := account.GetCorpFromContext(bCtx)
	user := account.GetUserFromContext(bCtx)

	orders, nextPageInfo := order.NewOrderRepository(bCtx).GetPagedOrdersForUserInCorp(user, corp, filters, pageInfo, "-id")
	
	fillOptions := eel.Map{}
	fillOptions["with_invoice"] = map[string]interface{}{
		"with_products": true,
	}
	order.NewFillOrderService(bCtx).Fill(orders, fillOptions)

	rows := order.NewEncodeOrderService(bCtx).EncodeMany(orders)
	ctx.Response.JSON(eel.Map{
		"orders": rows,
		"page_info": nextPageInfo.ToMap(),
	})
}

