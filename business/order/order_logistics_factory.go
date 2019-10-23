package order

import (
	"context"
	
	"github.com/gingerxman/eel"
	b_order_params "github.com/gingerxman/ginger-mall/business/order/params"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type OrderLogisticsFactory struct {
	eel.ServiceBase
}

func NewOrderLogisticsFactory(ctx context.Context) *OrderLogisticsFactory {
	service := new(OrderLogisticsFactory)
	service.Ctx = ctx
	return service
}

func (this *OrderLogisticsFactory) CreateLogistics(shipInfo *b_order_params.LogisticsParams) {
	o := eel.GetOrmFromContext(this.Ctx)
	dbModel := &m_order.OrderLogistics{
		OrderBid: shipInfo.Bid,
		EnableLogistics: shipInfo.EnableLogistics,
		ExpressCompanyName: shipInfo.ExpressCompanyName,
		ExpressNumber: shipInfo.ExpressNumber,
		Shipper: shipInfo.Shipper,
	}
	db := o.Create(dbModel)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order_logistics:insert_failed", "添加物流信息失败"))
	}
}

func init() {
}
