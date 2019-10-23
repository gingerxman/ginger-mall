package order

import (
	"context"
	"github.com/gingerxman/eel"
	"strings"
)

type InvoiceFinishedService struct {
	eel.ServiceBase
}

func (this *InvoiceFinishedService) DoClearance(invoice *Invoice){
	if invoice.Money.FinalMoney == 0{
		invoice.SetCleared(true)
		return
	}
	if invoice.IsCustomOrder(){
		this.doCustomOrderClearance(invoice)
	}
}

// doCustomOrderClearance 虚拟订单清算
func (this *InvoiceFinishedService) doCustomOrderClearance(invoice *Invoice){
	if invoice.NeedSyncClearance(){
		resp, _ := eel.NewResource(this.Ctx).Put("gplutus", "clearance.order_clearance", eel.Map{
			"bid": invoice.Bid,
			"order_status": "all",
			"sync": true,
		})
		if resp.IsSuccess(){
			invoice.SetCleared(false)
		}
	}
}

// DoCallback 同步回调
func (this *InvoiceFinishedService) DoCallback(invoice *Invoice){
	extraData := invoice.GetExtraData()
	if v, ok := extraData["callback_resource"]; ok{
		callbackResource := v.(string)
		sps := strings.Split(callbackResource, ":")
		serviceName := sps[0]
		resource := sps[1]
		resp, _ := eel.NewResource(this.Ctx).Put(serviceName, resource, eel.Map{
			"bid": invoice.Bid,
			"settlements": eel.ToJsonString([]string{}),
		})
		invoice.SetCallbackStatus(resp.IsSuccess())
	}else{
		invoice.SetCallbackStatus(true) // 没有回调要求的则置为成功
	}
}

// AfterFinished 订单完成后
func (this *InvoiceFinishedService) AfterFinished(invoice *Invoice){
	this.DoClearance(invoice)
	this.DoCallback(invoice)
}

func NewInvoiceFinishedService(ctx context.Context) *InvoiceFinishedService {
	service := new(InvoiceFinishedService)
	service.Ctx = ctx
	return service
}