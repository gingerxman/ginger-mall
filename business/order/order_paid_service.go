package order

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/event"
	
	"github.com/gingerxman/ginger-mall/business/events"
)

type OrderPaidService struct {
	eel.ServiceBase
}

// AfterPaid 订单完成支付后
func (this *OrderPaidService) AfterPaid(order *Order){

	if order.CustomType == "waiting_confirmation"{
		// 异步消息
		event.AsyncEvent.Send(events.ORDER_WAITING_CONFIRM, map[string]interface{}{
			"bid": order.Bid,
		})
	}

	// 处理充值订单
	if order.IsDepositOrder(){
		extraData := order.GetExtraData()
		if v, ok := extraData["deposit_imoney"]; ok{
			depositData := v.(map[string]interface{})
			resp, err := eel.NewResource(this.Ctx).Put("gplutus", "imoney.deposit", eel.Map{
				"bid": order.Bid,
				"target_user_id": order.UserId,
				"imoney_code": depositData["code"].(string),
				"amount": depositData["amount"].(float64),
			})
			invoice := NewInvoiceFromOrder(this.Ctx, order)
			if err != nil || !resp.IsSuccess(){
				// 充值失败，订单取消
				invoice.Cancel()
			}else{
				invoice.ForceFinish()
			}
		}
	}
}

func NewOrderPaidService(ctx context.Context) *OrderPaidService {
	service := new(OrderPaidService)
	service.Ctx = ctx
	return service
}