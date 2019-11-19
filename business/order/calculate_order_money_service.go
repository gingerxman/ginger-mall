package order

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	b_resource "github.com/gingerxman/ginger-mall/business/order/resource"
)

type CalculateOrderMoneyService struct {
	eel.ServiceBase
}

func NewCalculateOrderMoneyService(ctx context.Context) *CalculateOrderMoneyService {
	service := new(CalculateOrderMoneyService)
	service.Ctx = ctx
	return service
}

func (this *CalculateOrderMoneyService) Calculate(resources []business.IResource, purchaseInfo *PurchaseInfo, newOrder *NewOrder) *orderMoneyInfo {
	var finalMoney int = 0
	var productPrice int = 0
	var postageMoney int = 0
	var deductionMoney int = 0
	
	products := make([]business.IResource, 0)
	for _, resource := range resources {
		resourcePrice := resource.GetPrice()
		if resource.GetType() == b_resource.RESOURCE_TYPE_PRODUCT {
			productPrice += resourcePrice
			products = append(products, resource)
		}
		
		finalMoney += resourcePrice
		deductionMoney = resource.GetDeductionMoney(newOrder.GetDeductableMoney())
		finalMoney -= deductionMoney
		deductionMoney += deductionMoney
	}
	
	if finalMoney < 0 {
		//imoney不能抵扣运费，所以这里要将扣成负数的金额归零
		finalMoney = 0
	}
	
	postageMoney = NewCalculateOrderPostageService(this.Ctx).Calculate(products, purchaseInfo)
	//加入运费
	finalMoney += int(postageMoney)
	
	moneyInfo := &orderMoneyInfo{
		Postage: postageMoney,
		FinalMoney: finalMoney,
		EditMoney: 0.0,
		PayMoney: finalMoney,
		ProductPrice: productPrice,
	}
	return moneyInfo
}


func init() {
}
