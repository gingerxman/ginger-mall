package order

import (
	"context"
	
	"github.com/gingerxman/eel"
)

type EncodeOutlineService struct {
	eel.ServiceBase
}

func NewEncodeOutlineService(ctx context.Context) *EncodeOutlineService {
	service := new(EncodeOutlineService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeOutlineService) Encode(outline *OrderOutline) *ROrderOutline {

	return &ROrderOutline{
		TotalMoney: eel.Decimal(outline.totalMoney),
		IncrementMoney: eel.Decimal(outline.incrementMoney),
		TotalOrderCount: outline.totalOrderCount,
		IncrementOrderCount: outline.incrementOrderCount,
		TotalUserCount: outline.totalUserCount,
		IncrementUserCount: outline.incrementUserCount,
	}
}

func init() {
}
