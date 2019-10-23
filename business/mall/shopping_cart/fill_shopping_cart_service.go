package shopping_cart

import (
	"context"
	"github.com/gingerxman/eel"
	
)

type FillShoppingCartService struct {
	eel.ServiceBase
}

func NewFillShoppingCartService(ctx context.Context) *FillShoppingCartService {
	service := new(FillShoppingCartService)
	service.Ctx = ctx
	return service
}

//func (this *FillShoppingCartService) FillOne(shipInfo *ShipInfo, option eel.FillOption) {
//	this.Fill([]*ShipInfo{ shipInfo }, option)
//}
//
//func (this *FillShoppingCartService) Fill(shipInfos []*ShipInfo, option eel.FillOption) {
//	if len(shipInfos) == 0 {
//		return
//	}
//
//	ids := make([]int, 0)
//	for _, shipInfo := range shipInfos {
//		ids = append(ids, shipInfo.Id)
//	}
//	return
//}

func init() {
}
