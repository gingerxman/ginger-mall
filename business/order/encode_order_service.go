package order

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeOrderService struct {
	eel.ServiceBase
}

func NewEncodeOrderService(ctx context.Context) *EncodeOrderService {
	service := new(EncodeOrderService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeOrderService) Encode(order *Order) *ROrder {
	if order == nil {
		return nil
	}
	
	rInvoices := NewEncodeInvoiceService(this.Ctx).EncodeMany(order.Invoices)

	rOperationLogs:= make([]*ROperationLog, 0)
	for _, operationLog := range order.OperationLogs {
		rOperationLog := &ROperationLog{
			Id: operationLog.Id,
			OrderBid: operationLog.OrderBid,
			Type: operationLog.Type,
			Remark: operationLog.Remark,
			Action: operationLog.Action,
			Operator: operationLog.Operator,
			CreatedAt: operationLog.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rOperationLogs = append(rOperationLogs, rOperationLog)
	}
	rStatusLogs:= make([]*RStatusLog, 0)
	for _, statusLog := range order.StatusLogs {
		rStatusLog := &RStatusLog{
			Id: statusLog.Id,
			OrderBid: statusLog.OrderBid,
			FromStatus: statusLog.FromStatus,
			ToStatus: statusLog.ToStatus,
			Remark: statusLog.Remark,
			Operator: statusLog.Operator,
			CreatedAt: statusLog.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rStatusLogs = append(rStatusLogs, rStatusLog)
	}
	
	var totalProductPrice = 0
	var totalPostage = 0
	for _, rInvoice := range rInvoices {
		totalPostage += rInvoice.Postage
		totalProductPrice += rInvoice.ProductPrice
	}
	
	return &ROrder{
		Id: order.Id,
		Bid: order.Bid,
		CorpId: order.CorpId,
		UserId: order.UserId,
		Status: order.GetStatusText(),
		Invoices: rInvoices,
		Resources: order.GetResources(),
		FinalMoney: order.Money.FinalMoney,
		IsDeleted: order.IsDeleted,
		OperationLogs: rOperationLogs,
		StatusLogs: rStatusLogs,
		Remark : order.Remark,
		Message: order.Message,
		ExtraData: order.GetExtraData(),
		
		ProductPrice: totalProductPrice,
		Postage: totalPostage,
		
		PaymentTime: order.PaymentTime.Format("2006-01-02 15:04:05"),
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeOrderService) EncodeMany(orders []*Order) []*ROrder {
	rDatas := make([]*ROrder, 0)
	for _, order := range orders {
		rDatas = append(rDatas, this.Encode(order))
	}
	
	return rDatas
}

func init() {
}
