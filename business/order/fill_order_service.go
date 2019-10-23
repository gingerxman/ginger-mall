package order

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type FillOrderService struct {
	eel.ServiceBase
}

func NewFillOrderService(ctx context.Context) *FillOrderService {
	service := new(FillOrderService)
	service.Ctx = ctx
	return service
}

func (this *FillOrderService) _fillInvoice(invoices []*Invoice, option map[string]interface{}) {
	if len(invoices) == 0 {
		return
	}
	
	ids := make([]int, 0)
	bids := make([]string, 0)
	for _, invoice := range invoices {
		ids = append(ids, invoice.Id)
		bids = append(bids, invoice.Bid)
	}
	
	if enableOption, ok := option["with_products"]; ok && enableOption.(bool) {
		this.fillProducts(invoices, ids, bids)
	}
	
	if enableOption, ok := option["with_logistics"]; ok && enableOption.(bool) {
		this.fillLogistics(invoices, ids, bids)
	}
}

func (this *FillOrderService) Fill(orders []*Order, option map[string]interface{}) {
	if len(orders) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, order := range orders {
		ids = append(ids, order.Id)
	}
	
	if subOption, ok := option["with_invoice"]; ok {
		this.fillInvoice(orders, ids, subOption.(map[string]interface{}))
	}

	if withOperationLog, ok := option["with_operation_log"]; ok {
		if withOperationLog.(bool){
			this.fillOperationLog(orders, ids)
		}
	}
	if withStatusLog, ok := option["with_status_log"]; ok {
		if withStatusLog.(bool){
			this.fillStatusLog(orders, ids)
		}
	}
}

func (this *FillOrderService) fillProducts(invoices []*Invoice, ids []int, bids []string) {
	if len(ids) == 0 {
		return
	}
	
	//构建<id, invoice>
	id2invoice := make(map[int]*Invoice)
	for _, invoice := range invoices {
		id2invoice[invoice.Id] = invoice
	}
	
	orderProducts := NewOrderProductRepository(this.Ctx).GetOrderProducts(ids)
	for _, orderProduct := range orderProducts {
		invoiceId := orderProduct.OrderId
		if invoice, ok := id2invoice[invoiceId]; ok {
			invoice.AppendProduct(orderProduct)
		}
	}
}

//FillInvoiceProducts 为单个invoice填充商品集合
func (this *FillOrderService) FillInvoiceProducts(invoice *Invoice) {
	this.fillProducts([]*Invoice{invoice}, []int{invoice.Id}, []string{invoice.Bid})
}

func (this *FillOrderService) fillLogistics(invoices []*Invoice, ids []int, bids []string) {
	if len(ids) == 0 {
		return
	}
	
	//构建<id, invoice>
	bid2invoice := make(map[string]*Invoice)
	for _, invoice := range invoices {
		bid2invoice[invoice.Bid] = invoice
	}
	
	//获得OrderLogistics models
	models := make([]*m_order.OrderLogistics, 0)
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.OrderLogistics{}).Where("order_bid__in", bids).Find(&models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return
	}
	
	//填充invoice.Logistics
	for _, model := range models {
		invoiceBid := model.OrderBid
		if invoice, ok := bid2invoice[invoiceBid]; ok {
			invoice.SetLogistics(model)
		}
	}
}

func (this *FillOrderService) fillInvoice(orders []*Order, ids []int, option map[string]interface{}) {
	invoices := make([]*Invoice, 0)
	customOrderIds := make([]int, 0)
	productOrderIds := make([]int, 0)
	
	for _, order := range orders {
		if order.IsInvoice() {
			invoices = append(invoices, NewInvoiceFromOrder(this.Ctx, order))
		} else if order.IsCustomOrder() {
			customOrderIds = append(customOrderIds, order.Id)
		} else {
			productOrderIds = append(productOrderIds, order.Id)
		}
	}
	
	//处理商品订单
	orderRepository := NewOrderRepository(this.Ctx)
	invoices = append(invoices, orderRepository.GetInvoicesByOrderIds(productOrderIds)...)
	invoices = append(invoices, orderRepository.GetInvoicesByIds(customOrderIds)...)
	
	this._fillInvoice(invoices, option)
	
	//构建<id, order>
	id2order := make(map[int]*Order)
	for _, order := range orders {
		id2order[order.Id] = order
	}
	
	//填充order.Invoices
	for _, invoice := range invoices {
		if !invoice.IsCustomOrder() {
			//非custom order的original order id才有意义，否则invoice.OriginalOrderId == invoice.Id
			if order, ok := id2order[invoice.OriginalOrderId]; ok {
				order.AddInvoice(invoice)
			}
		}
		
		if order, ok := id2order[invoice.Id]; ok {
			order.AddInvoice(invoice)
		}
	}
}

func (this *FillOrderService) fillOperationLog(orders []*Order, ids []int) {
	bids := make([]string, 0)
	for _, order := range orders {
		bids = append(bids, order.Bid)
	}
	//从db中获取数据集合
	var models []*m_order.OrderOperationLog
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.OrderOperationLog{}).Where("order_bid__in", bids).Find(&models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return
	}

	//构建<bid, [operationLog]>
	bid2operationLogs := make(map[string][]*OrderOperationLog)
	for _, model := range models {
		if _, ok := bid2operationLogs[model.OrderBid]; ok{
			bid2operationLogs[model.OrderBid] = append(bid2operationLogs[model.OrderBid], NewOrderOperationLogFromModel(this.Ctx, model))
		} else {
			bid2operationLogs[model.OrderBid] = []*OrderOperationLog{NewOrderOperationLogFromModel(this.Ctx, model)}
		}
	}

	//填充order的OperationLog对象
	for _, order := range orders {
		if operationLogs, ok := bid2operationLogs[order.Bid]; ok {
			order.OperationLogs = operationLogs
		}
	}
}

func (this *FillOrderService) fillStatusLog(orders []*Order, ids []int) {
	bids := make([]string, 0)
	for _, order := range orders {
		bids = append(bids, order.Bid)
	}
	//从db中获取数据集合
	var models []*m_order.OrderStatusLog
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.OrderStatusLog{}).Where("order_bid__in", bids).Find(&models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return
	}

	//构建<bid, [statusLog]>
	bid2statusLogs := make(map[string][]*OrderStatusLog)
	for _, model := range models {
		if _, ok := bid2statusLogs[model.OrderBid]; ok{
			bid2statusLogs[model.OrderBid] = append(bid2statusLogs[model.OrderBid], NewOrderStatusLogFromModel(this.Ctx, model))
		} else {
			bid2statusLogs[model.OrderBid] = []*OrderStatusLog{NewOrderStatusLogFromModel(this.Ctx, model)}
		}
	}

	//填充order的StatusLog对象
	for _, order := range orders {
		if statusLogs, ok := bid2statusLogs[order.Bid]; ok {
			order.StatusLogs = statusLogs
		}
	}
}


func init() {
}
