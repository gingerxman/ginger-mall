package order
import (
	"context"
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"time"
)

type OrderStatusLog struct {
	eel.EntityBase

	Id int
	OrderBid string
	FromStatus string
	ToStatus string
	Remark string
	Operator string
	CreatedAt time.Time
}

func NewOrderStatusLogFromModel(ctx context.Context, dbModel *m_order.OrderStatusLog) *OrderStatusLog{
	instance := new(OrderStatusLog)
	instance.Ctx = ctx
	instance.Id = dbModel.Id
	instance.OrderBid = dbModel.OrderBid
	instance.FromStatus = m_order.STATUS2STR[dbModel.FromStatus]
	instance.ToStatus = m_order.STATUS2STR[dbModel.ToStatus]
	instance.Remark = dbModel.Remark
	instance.Operator = dbModel.Operator
	instance.CreatedAt = dbModel.CreatedAt
	return instance
}