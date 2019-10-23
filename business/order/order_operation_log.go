package order
import (
	"context"
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"time"
)

type OrderOperationLog struct {
	eel.EntityBase

	Id int
	OrderBid string
	Type string
	Remark string
	Action string
	Operator string
	CreatedAt time.Time
}

func NewOrderOperationLogFromModel(ctx context.Context, dbModel *m_order.OrderOperationLog) *OrderOperationLog{
	instance := new(OrderOperationLog)
	instance.Ctx = ctx
	instance.Id = dbModel.Id
	instance.OrderBid = dbModel.OrderBid
	instance.Type = m_order.OPERATONTYPE2CODE[dbModel.Type]
	instance.Remark = dbModel.Remark
	instance.Action = dbModel.Action
	instance.Operator = dbModel.Operator
	instance.CreatedAt = dbModel.CreatedAt
	return instance
}