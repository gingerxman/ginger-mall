package order

import (
	"encoding/json"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_order "github.com/gingerxman/ginger-mall/business/order"
	"github.com/gingerxman/ginger-mall/business/order/params"
)

type ShippedInvoice struct {
	eel.RestResource
}

func (this *ShippedInvoice) Resource() string {
	return "order.shipped_invoice"
}

func (this *ShippedInvoice) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"ship_infos:json-array"},
	}
}

func (this *ShippedInvoice) Put(ctx *eel.Context) {
	req := ctx.Request
	shipInfos := req.GetJSONArray("ship_infos")
	
	// 解析LogisticsParams
	bids := make([]string, 0)
	bid2shipInfo := make(map[string]*params.LogisticsParams, 0)
	for _, shipInfo := range shipInfos{
		byteShipInfo, _ := json.Marshal(shipInfo)
		data := params.ParseLogisticsParams(string(byteShipInfo))
		bid := data.Bid
		bids = append(bids, bid)
		bid2shipInfo[bid] = data
	}
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	invoices := b_order.NewOrderRepository(bCtx).GetInvoicesByBidsForCorp(corp, bids)
	if len(invoices) == 0{
		ctx.Response.Error("order:ship_fail", "无效的bids")
	} else {
		for _, invoice := range invoices{
			invoice.Ship(bid2shipInfo[invoice.Bid])
		}
		ctx.Response.JSON(eel.Map{})
	}
}
