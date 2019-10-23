package order

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeLogisticsService struct {
	eel.ServiceBase
}

func NewEncodeOrderLogisticsService(ctx context.Context) *EncodeLogisticsService {
	service := new(EncodeLogisticsService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeLogisticsService) Encode(logistics *OrderLogistics) *ROrderLogistics {
	if logistics == nil {
		return nil
	}

	return &ROrderLogistics{
		Id: logistics.Id,
		OrderBid: logistics.OrderBid,
		EnableLogistics: logistics.EnableLogistics,
		ExpressCompanyName: logistics.ExpressCompanyName,
		ExpressNumber: logistics.ExpressNumber,
		Shipper: logistics.Shipper,
	}
}

func init() {
}
