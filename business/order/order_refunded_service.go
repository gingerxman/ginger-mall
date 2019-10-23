package order

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/event"
	
	"github.com/gingerxman/ginger-mall/business/events"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type OrderRefundedService struct {
	eel.ServiceBase
}

// AfterPaid 订单完成支付后
func (this *OrderRefundedService) AfterRefunded(order *Order){
	if order.Status == m_order.ORDER_STATUS_REFUNDED{
		// 处理订单交易
		resp, _ := eel.NewResource(this.Ctx).Put("gplutus", "clearance.refund_orders", eel.Map{
			"order_bids": eel.ToJsonString([]string{order.Bid}),
		})
		if resp.IsSuccess(){
			// 异步消息
			event.AsyncEvent.Send(events.ORDER_REFUNDED, map[string]interface{}{
				"bid": order.Bid,
			})
		}
	}
}

func NewOrderRefundedService(ctx context.Context) *OrderRefundedService {
	service := new(OrderRefundedService)
	service.Ctx = ctx
	return service
}