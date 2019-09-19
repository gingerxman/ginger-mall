package limit_zone

import (
	"context"
	"github.com/gingerxman/eel"
)

type FillLimitZoneService struct {
	eel.ServiceBase
}

func NewFillLimitZoneService(ctx context.Context) *FillLimitZoneService {
	service := new(FillLimitZoneService)
	service.Ctx = ctx
	return service
}

func (this *FillLimitZoneService) Fill(limigZones []*LimitZone, option eel.FillOption) {
	return
}


func init() {
}
