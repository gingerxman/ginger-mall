package ship_info

import (
	"context"
	"github.com/gingerxman/eel"
	
)

type FillShipInfoService struct {
	eel.ServiceBase
}

func NewFillShipInfoService(ctx context.Context) *FillShipInfoService {
	service := new(FillShipInfoService)
	service.Ctx = ctx
	return service
}

func (this *FillShipInfoService) FillOne(shipInfo *ShipInfo, option eel.FillOption) {
	this.Fill([]*ShipInfo{ shipInfo }, option)
}

func (this *FillShipInfoService) Fill(shipInfos []*ShipInfo, option eel.FillOption) {
	if len(shipInfos) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, shipInfo := range shipInfos {
		ids = append(ids, shipInfo.Id)
	}
	return
}

func init() {
}
