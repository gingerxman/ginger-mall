package order

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/gingerxman/eel"
)

type FinishedOrderHandler struct {
	eel.ServiceBase
}

// requestSettlement 结算
func (this *FinishedOrderHandler) requestSettlement(invoice *Invoice) error{
	v := 1
	sync := false
	if invoice.IsDepositOrder(){
		v = 2
		sync = true
	}
	_, err := eel.NewResource(this.Ctx).Put("gplutus", "clearance.order_clearance", map[string]interface{}{
		"bid": invoice.Bid,
		"_v": v,
		"sync": sync,
	})

	if err != nil{
		eel.Logger.Error(err)
		return err
	}else{
		return invoice.SetCleared(false)
	}
}

// requestConsumingFrozenRecord 消费冻结的虚拟资产
func (this *FinishedOrderHandler) requestConsumingFrozenRecord(invoice *Invoice) error{
	reqResource := eel.NewResource(this.Ctx)
	for _, resource := range invoice.GetResources(){
		if resource["type"] == "imoney"{
			if ids, ok := resource["frozen_record_ids"]; ok{
				recordIds := ids.([]interface{})
				if len(recordIds) > 0{
					rid , _ := recordIds[0].(json.Number).Int64()
					_, err := reqResource.Put("gplutus", "imoney.settled_frozen_record", eel.Map{
						"id": int(rid),
						"extra_data": eel.ToJsonString(map[string]interface{}{
							"action": fmt.Sprintf("order_imoney_usages: bid_%s", invoice.Bid),
							"amount": resource["count"],
							"imoney_code": resource["code"],
							"bid": invoice.Bid,
						}),
					})
					if err != nil{
						eel.Logger.Error(err)
						return err
					}
				}
			}
		}
	}
	return nil
}

// DoSettlement 订单完成后进行清算
func (this *FinishedOrderHandler) DoSettlement(order *Order){
	invoices := order.Invoices
	if len(invoices) == 0{
		NewFillOrderService(this.Ctx).Fill([]*Order{order}, eel.Map{
			"with_invoice": eel.Map{
				"with_products": false,
			},
		})
		invoices = order.Invoices
	}

	for _, invoice := range invoices{
		err1 := this.requestConsumingFrozenRecord(invoice)
		if err1 != nil{
			// common.UrgentMessage.Put(fmt.Sprintf("> 订单(%s)虚拟资产解冻失败 \n\n > errMsg: %s \n\n  ", invoice.Bid, err1.Error()))
		}
		err2 := this.requestSettlement(invoice)
		if err2 != nil{
			// common.UrgentMessage.Put(fmt.Sprintf("> 结算订单(%s)失败 \n\n >errMsg: %s \n\n ", invoice.Bid, err2.Error()))
		}
	}
}

// DoCallback 订单完成后回调
func (this *FinishedOrderHandler) DoCallback(order *Order){

	extraData := order.GetExtraData()
	if _, ok := extraData["callback_resource"]; ok{
		order.SetCallbackStatus(false)
	}
}

func NewFinishedOrderHandler(ctx context.Context) *FinishedOrderHandler {
	handler := new(FinishedOrderHandler)
	handler.Ctx = ctx
	return handler
}