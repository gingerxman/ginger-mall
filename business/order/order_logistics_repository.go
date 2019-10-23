package order

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type OrderLogisticsRepository struct {
	eel.ServiceBase
}

func NewOrderLogisticsRepository(ctx context.Context) *OrderLogisticsRepository {
	service := new(OrderLogisticsRepository)
	service.Ctx = ctx
	return service
}

func (this *OrderLogisticsRepository) GetOrderLogisticsByBid(bid string) *OrderLogistics {
	o := eel.GetOrmFromContext(this.Ctx)

	var model m_order.OrderLogistics
	db := o.Model(&m_order.OrderLogistics{}).Where("OrderBid", bid).Take(&model)
	err := db.Error

	if err != nil {
		eel.Logger.Error(err)
		return nil
	}

	return NewOrderLogisticsFromModel(this.Ctx, &model)
}

func init() {
}
