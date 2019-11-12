package order

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)
const ORDER_OPERATION_CREATE = "create_order"
const ORDER_OPERATION_START_PAY = "start_pay_order"
const ORDER_OPERATION_PAY = "pay_order"
const ORDER_OPERATION_CONFIRM = "confirm_order"
const ORDER_OPERATION_SHIP = "ship_order"
const ORDER_OPERATION_FINISH = "finish_order"
const ORDER_OPERATION_FINISH_BY_SYSTEM = "finish_order_by_system"
const ORDER_OPERATION_CANCEL = "cancel_order"
const ORDER_OPERATION_REFUND = "refund_order"
const ORDER_OPERATION_REQUEST_PLATFORM_REFUND = "request_platform_refund"
const ORDER_OPERATION_REJECT_REFUND = "reject_refund"
const ORDER_OPERATION_UPDATE_LOGISTICS_INFO = "update_order_logistics_info"
const ORDER_OPERATION_FINISH_REFUND = "finish_refund"

var OPERATION2STATUSPAIR = map[string] interface{}{
	ORDER_OPERATION_CREATE: map[string] interface{} {
		"from": -1,
		"to": m_order.ORDER_STATUS_WAIT_PAY,
	},
	ORDER_OPERATION_START_PAY: map[string] interface{}{
		"from": m_order.ORDER_STATUS_WAIT_PAY,
		"to": m_order.ORDER_STATUS_PAYING,
	},
	ORDER_OPERATION_PAY: map[string]interface{}{
		"from": m_order.ORDER_STATUS_PAYING,
		"to": m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM,
	},
	ORDER_OPERATION_CONFIRM: map[string]interface{}{
		"from": m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM,
		"to": m_order.ORDER_STATUS_PAYED_NOT_SHIP,
	},
	ORDER_OPERATION_SHIP: map[string]interface{}{
		"from": m_order.ORDER_STATUS_PAYED_NOT_SHIP,
		"to": m_order.ORDER_STATUS_PAYED_SHIPED,
	},
	ORDER_OPERATION_FINISH: map[string]interface{}{
		"from": m_order.ORDER_STATUS_PAYED_SHIPED,
		"to": m_order.ORDER_STATUS_SUCCESSED,
	},
	ORDER_OPERATION_FINISH_BY_SYSTEM: map[string]interface{}{
		"from": m_order.ORDER_STATUS_PAYED_SHIPED,
		"to": m_order.ORDER_STATUS_SUCCESSED,
	},
	ORDER_OPERATION_CANCEL: map[string]interface{}{
		"from": m_order.ORDER_STATUS_WAIT_PAY,
		"to": m_order.ORDER_STATUS_CANCEL,
	},
	ORDER_OPERATION_REFUND: map[string]interface{}{
		"from": nil,
		"to": m_order.ORDER_STATUS_REFUNDING,
	},
	ORDER_OPERATION_REQUEST_PLATFORM_REFUND: map[string]interface{}{
		"from": m_order.ORDER_STATUS_REFUNDING,
		"to": m_order.ORDER_STATUS_PLATFORM_REFUNDING,
	},
	ORDER_OPERATION_REJECT_REFUND: map[string]interface{}{
		"from": m_order.ORDER_STATUS_REFUNDING,
		"to": nil,
	},
	ORDER_OPERATION_FINISH_REFUND: map[string]interface{}{
		"from": m_order.ORDER_STATUS_PLATFORM_REFUNDING,
		"to": m_order.ORDER_STATUS_REFUNDED,
	},
}
type OperationData struct {
	Order interface{}
	Action string
	FromStatus interface{}
	ToStatus interface{}
	Data  *eel.Map
}

type OrderLogService struct {
	eel.ServiceBase
}

func NewOrderLogService(ctx context.Context) *OrderLogService {
	service := new(OrderLogService)
	service.Ctx = ctx
	return service
}

// 操作记录
func (this *OrderLogService) LogOperation(operationData *OperationData) {
	var orders []*Order
	var targetOrder *Order
	order := operationData.Order
	action := operationData.Action
	fromStatus := operationData.FromStatus
	toStatus := operationData.ToStatus
	dataStr, _ := json.Marshal(operationData.Data)
	switch order.(type) {
		case *Order:
			targetOrder = order.(*Order)
			orders = append(orders, targetOrder)
		case *Invoice:
			targetOrder = NewOrderRepository(this.Ctx).GetOrderByBid(order.(*Invoice).Bid)
			orders = append(orders, targetOrder)
		default:
			targetOrder = nil
	}
	if targetOrder != nil {
		if action == ORDER_OPERATION_CREATE || action == ORDER_OPERATION_START_PAY || action == ORDER_OPERATION_PAY || action == ORDER_OPERATION_CANCEL {
			//以下操作，订单，出货单，都要记录日志
			//创建订单，支付订单，取消订单
			originalOrderBid := targetOrder.Bid

			if !targetOrder.IsInvoice() {
				if len(targetOrder.Invoices) == 0 {
					fillService := NewFillOrderService(this.Ctx)
					fillService.Fill(orders, eel.Map{
						"with_invoice": eel.Map{},
					})
				}

				for _, invoice := range targetOrder.Invoices {
					if invoice.Bid != originalOrderBid {
						orders = append(orders, NewOrderRepository(this.Ctx).GetOrderByBid(invoice.Bid))
					}
				}
			}
		}

		// 确定operator与operationType
		var operator string
		var operationType int
		if action == ORDER_OPERATION_START_PAY || action == ORDER_OPERATION_PAY || action == ORDER_OPERATION_CANCEL || action == ORDER_OPERATION_REFUND {
			operator = "user"
			operationType = m_order.ORDER_OPERATION_TYPE_MEMBER
		} else if action == ORDER_OPERATION_FINISH_BY_SYSTEM {
			operator = "system"
			operationType = m_order.ORDER_OPERATION_TYPE_SYSTEM
		} else {
			user := account.GetUserFromContext(this.Ctx)
			spew.Dump(user)
			spew.Dump(user.Id)
			corpUsers := account.NewCorpUserRepository(this.Ctx).GetCorpUsers([]int{user.Id})
			operator = "operator"
			if corpUsers != nil && len(corpUsers) > 0 {
				operator = corpUsers[0].Name
			}
			operationType = m_order.ORDER_OPERATION_TYPE_MALL_OPERATOR
		}
		o := eel.GetOrmFromContext(this.Ctx)
		for _, order := range orders {
			//记录operation日志
			operationLogModel := m_order.OrderOperationLog{
				OrderBid: order.Bid,
				Operator: operator,
				Action:   action,
				Type:     operationType,
				Remark:   string(dataStr), //json字符串
			}
			db := o.Create(&operationLogModel)
			err := db.Error
			if err != nil {
				eel.Logger.Error(err)
				panic(eel.NewBusinessError("order:create_operation_log_failed", "创建操作日志失败"))
			}
			//记录status日志
			statusPair, ok := OPERATION2STATUSPAIR[action]
			if ok {
				statusPair := statusPair.(map[string]interface{})
				targetFromStatus, _ := statusPair["from"]
				if targetFromStatus == nil {
					targetFromStatus = fromStatus
					if targetFromStatus == nil {
						targetFromStatus = m_order.ORDER_STATUS_NONSENSE
					}
				}
				targetToStatus,_ := statusPair["to"]
				if targetToStatus == nil {
					targetToStatus = toStatus
					if targetToStatus == nil{
						targetToStatus = m_order.ORDER_STATUS_NONSENSE
					}
				}
				statusLogModel := m_order.OrderStatusLog{
					OrderBid:   order.Bid,
					FromStatus: targetFromStatus.(int),
					ToStatus:   targetToStatus.(int),
					Operator:   operator,
				}
				db := o.Create(&statusLogModel)
				err := db.Error
				if err != nil {
					eel.Logger.Error(err)
					panic(eel.NewBusinessError("order:create_status_log_failed", "创建状态日志失败"))
				}
			}
		}
	}

}


func init() {
}
