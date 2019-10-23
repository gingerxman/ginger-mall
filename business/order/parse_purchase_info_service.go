package order

import (
	"context"
	"github.com/gingerxman/eel"
	
)

type ParsePurchaseInfoService struct {
	eel.ServiceBase
}

func NewParsePurchaseInfoService(ctx context.Context) *ParsePurchaseInfoService {
	service := new(ParsePurchaseInfoService)
	service.Ctx = ctx
	return service
}

func (this *ParsePurchaseInfoService) Parse() {

}


func init() {
}
