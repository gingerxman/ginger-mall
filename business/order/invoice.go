package order

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/event"
	"github.com/gingerxman/eel/log"
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
	
	"github.com/gingerxman/ginger-mall/business/events"
	b_order_params "github.com/gingerxman/ginger-mall/business/order/params"
	"github.com/gingerxman/ginger-mall/business/order/resource"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"reflect"
	"strconv"
	"time"
)

var enableSyncClearance bool = false

type invoiceLogistics struct {
	EnableLogistics bool
	ExpressCompanyName string
	ExpressNumber string
	Shipper string
}

type Invoice struct {
	eel.EntityBase
	Order
	
	SupplierId int
	PaymentType string
	Postage float64
	PostageStrategy int
	
	ShipInfo *ShipInfo //收货信息
	Products []*OrderProduct
	Logistics *invoiceLogistics //物流信息
	
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NeedSyncClearance 是否需要同步清算
func (this *Invoice) NeedSyncClearance() bool{
	extraData := this.GetExtraData()
	if v, ok := extraData["need_sync_clearance"]; ok && v.(bool){
		return true
	}
	return false
}

func (this *Invoice) AppendProduct(orderProduct *OrderProduct) {
	this.Products = append(this.Products, orderProduct)
}

func (this *Invoice) SetLogistics(logisticsModel *m_order.OrderLogistics) {
	this.Logistics = &invoiceLogistics{
		EnableLogistics: logisticsModel.EnableLogistics,
		ExpressCompanyName: logisticsModel.ExpressCompanyName,
		ExpressNumber: logisticsModel.ExpressNumber,
		Shipper: logisticsModel.Shipper,
	}
}

func (this *Invoice) save() {
	ormParams := gorm.Params{
		"remark": this.Remark,
		"status": this.Status,
		"update_at": time.Now(),
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{}).Where("id", this.Id).Update(ormParams)
	if db.Error != nil{
		log.Logger.Error(db.Error)
		panic(eel.NewBusinessError("invoice:save_failed", "保存数据失败"))
	}
}

// GetOrder 获取订单
func (this *Invoice) GetOrder() *Order{
	if this.IsInvoice(){

	}else{
		return &this.Order
	}

	return nil
}

func (this *Invoice) GetOriginalOrder() *Order {
	return NewOrderRepository(this.Ctx).GetOrderById(this.OriginalOrderId)
}

func (this *Invoice) Finish() {
	if this.Status != m_order.ORDER_STATUS_PAYED_SHIPED {
		panic(eel.NewBusinessError("invoice: finished_failed", fmt.Sprintf("出货单状态不对(%s)", this.GetStatusText())))
	}
	
	toStatus := m_order.ORDER_STATUS_SUCCESSED
	this.Status = toStatus

	this.save()

	//如果开启同步清算，则清算
	if enableSyncClearance {
		err := NewClearanceService(this.Ctx).ClearInvoice(this)
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	
	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_FINISH,
	})

	NewInvoiceFinishedService(this.Ctx).AfterFinished(this)

	// 结算, todo: 订单下所有出货单都已完成的条件下，对订单进行清算
	//order := this.GetOrder()
	//if order != nil {
	//	order.finish()
	//}
}

func (this *Invoice) ForceFinish() {
	//临时改变订单状态，以通过Finish中的状态检查
	this.Status = m_order.ORDER_STATUS_PAYED_SHIPED
	
	this.Finish()
}

// GetSettlementData 获取结算规则数据
func (this *Invoice) GetSettlementData(corpRelatedUserId int, productName string) map[string]interface{} {
	extraData := this.GetExtraData()
	destUserId := 0
	val, ok := extraData["relevant_user_id"]
	if ok{
		destUserId = int(val.(float64))
	}

	ruleName := fmt.Sprintf("%d.%s", corpRelatedUserId, productName)
	settlementAmount := this.Money.FinalMoney
	imoneyCode := "cash"
	sourceUserId := this.UserId

	if settlementRule, ok := extraData["settlement_rule"]; ok{
		ruleData := settlementRule.(map[string]interface{})
		ruleName = ruleData["name"].(string)

		if amount, ok := ruleData["amount"]; ok{
			settlementAmount = amount.(int)
		}
		if code, ok := ruleData["imoney_code"]; ok{
			imoneyCode = code.(string)
		}
		if suid, ok := ruleData["source_user_id"]; ok{
			sourceUserId = int(suid.(float64))
		}
		if val, ok = extraData["relevant_corp_id"]; ok && int(val.(float64)) == this.CorpId{
			destUserId = corpRelatedUserId
		}
		if val, ok = ruleData["dest_user_id"]; ok{
			// 类型容错
			if reflect.TypeOf(val).String()=="string"{
				destUserId, _ = strconv.Atoi(val.(string))
			}else{
				destUserId = int(val.(float64))
			}
		}
	}else{
		if corpRelatedUserId == 0{
			ruleName = fmt.Sprintf("sys.%s", productName)
		}
	}

	if settlementAmount <= 0{
		err := this.SetCleared(true)
		if err != nil{
			return nil
		}
	}

	settlementData := map[string]interface{}{
		"bid": this.Bid,
		"money": settlementAmount,
		"imoney_code": imoneyCode,
		"name": ruleName,
		"source_user_id": sourceUserId,
	}
	if destUserId != 0{
		settlementData["dest_user_id"] = destUserId
	}

	return settlementData
}

// SetCleared 设置已清算
func (this *Invoice) SetCleared(willPublishEvent bool) error{
	this.IsCleared = true
	ormParams := gorm.Params{
		"is_cleared": this.IsCleared,
		"update_at": time.Now(),
	}

	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{}).Where("Bid", this.Bid).Update(ormParams)

	if db.Error != nil{
		eel.Logger.Error(db.Error)
		// common.UrgentMessage.Put(fmt.Sprintf("> 结算订单(%s)成功, 但更新数据库字段is_cleared失败 \n\n ", this.Bid))
	}else{
		if willPublishEvent{
			// 异步消息
			event.AsyncEvent.Send(events.ORDER_SETTLED, map[string]interface{}{
				"bid": this.Bid,
			})
		}
	}
	return db.Error
}

//确认出货单
func (this *Invoice) Confirm() {
	if this.Status != m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM {
		panic(eel.NewBusinessError("invoice: confirmed_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}
	targetStatus := m_order.ORDER_STATUS_PAYED_NOT_SHIP

	//更改order的状态
	this.Status = targetStatus
	this.save()
	
	//更改original order的状态
	this.GetOriginalOrder().ChangeStatusToNonsense()
	
	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_CONFIRM,
	})
}

//出货单发货
func (this *Invoice) Ship(shipInfo *b_order_params.LogisticsParams) {
	if this.Status != m_order.ORDER_STATUS_PAYED_NOT_SHIP {
		panic(eel.NewBusinessError("invoice: ship_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}
	targetStatus := m_order.ORDER_STATUS_PAYED_SHIPED

	// 保存物流信息
	NewOrderLogisticsFactory(this.Ctx).CreateLogistics(shipInfo)

	//更改order的状态
	this.Status = targetStatus
	this.save()
	//记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_SHIP,
	})
}

// 修改订单物流信息
func (this *Invoice) UpdateLogistics(shipInfo *b_order_params.LogisticsParams) {
	// 更新物流信息
	orderLogistics := NewOrderLogisticsRepository(this.Ctx).GetOrderLogisticsByBid(this.Bid)
	orderLogistics.Update(shipInfo)

	//更新订单时间
	this.save()
}

// Cancel 取消出货单
func (this *Invoice) Cancel(reason string){
	if this.Status == m_order.ORDER_STATUS_WAIT_PAY  || this.Status == m_order.ORDER_STATUS_PAYING{
		// 释放出货单中使用的resource
		resources := resource.NewParseResourceService(this.Ctx).ParseFromOrderResources(this.GetResources())
		resource.NewAllocateResourceService(this.Ctx).Release(resources)
		
		o := eel.GetOrmFromContext(this.Ctx)
		this.Status = m_order.ORDER_STATUS_CANCEL
		db := o.Model(&m_order.Order{}).Where("id", this.Id).Update(gorm.Params{
			"status": this.Status,
			"cancel_reason": reason,
			"updated_at": time.Now(),
		})
		if db.Error != nil {
			eel.Logger.Error(db.Error)
		}

		// 记录操作日志
		NewOrderLogService(this.Ctx).LogOperation(&OperationData{
			Order: this,
			Action: ORDER_OPERATION_CANCEL,
		})
	} else {
		panic(eel.NewBusinessError("invoice: cancel_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}

}

// Refunding 用户申请操作退款
func (this *Invoice) Refunding(reason string){
	fromStatus := this.Status
	this.Status = m_order.ORDER_STATUS_REFUNDING
	this.save()
	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_REFUND,
		FromStatus: fromStatus,
		Data: &eel.Map{
			"reason":reason,
		},
	})
}

// RefundingForPlatform 向平台申请操作退款
func (this *Invoice) RefundingForPlatform(){
	if this.Status != m_order.ORDER_STATUS_REFUNDING {
		panic(eel.NewBusinessError("invoice: request_refund_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}
	this.Status = m_order.ORDER_STATUS_PLATFORM_REFUNDING
	this.save()
	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_REQUEST_PLATFORM_REFUND,
	})
}

//RejectRefund 驳回退款申请
func (this *Invoice) RejectRefund(reason string)  {
	if this.Status != m_order.ORDER_STATUS_PLATFORM_REFUNDING {
		panic(eel.NewBusinessError("invoice: reject_refund_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}
	// 获取申请退款的的记录
	var model m_order.OrderStatusLog
	o := eel.GetOrmFromContext(this.Ctx)
	err := o.Model(&m_order.OrderStatusLog{}).Where(eel.Map{
		"order_bid": this.Bid,
		"to_status": m_order.ORDER_STATUS_REFUNDING,
	}).Order("id desc").Take(&model)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("invoice: get_status_log_failed", "获取订单状态log失败"))
	}
	toStatus := model.FromStatus
	this.Status = toStatus
	this.save()

	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_REJECT_REFUND,
		ToStatus: toStatus,
		Data: &eel.Map{
			"reason":reason,
		},
	})
}

// FinishRefund 确认退款
func (this *Invoice) FinishRefund(){
	if this.Status != m_order.ORDER_STATUS_PLATFORM_REFUNDING {
		panic(eel.NewBusinessError("invoice: finish_refund_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}
	resources := resource.NewParseResourceService(this.Ctx).ParseFromOrderResources(this.GetResources())
	resource.NewAllocateResourceService(this.Ctx).Release(resources)
	this.Status = m_order.ORDER_STATUS_REFUNDED
	this.save()
	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_FINISH_REFUND,
	})
}

// UpdateRemark 修改备注
func (this *Invoice) UpdateRemark(remark string){
	fmt.Println(remark, "remark")
	if this.Remark != remark{
		this.Remark = remark
		this.save()
	}
}

func NewInvoiceFromModel(ctx context.Context, model *m_order.Order) *Invoice {
	orderInstance := NewOrderFromModel(ctx, model)
	
	return NewInvoiceFromOrder(ctx, orderInstance)
}

func NewInvoiceFromOrder(ctx context.Context, order *Order) *Invoice {
	
	instance := &Invoice{
		Order: *order,
	}
	
	instance.Ctx = ctx
	instance.Model = order.Model
	
	model := instance.Model.(*m_order.Order)
	instance.CreatedAt = model.CreatedAt
	instance.SupplierId = model.SupplierId
	instance.ShipInfo = &ShipInfo{
		Name: model.ShipName,
		Phone: model.ShipPhone,
		Address: model.ShipAddress,
		AreaCode: model.ShipAreaCode,
	}
	
	//var ship shipArea
	//if model.ShipAreaCode == "0_0_0" || model.ShipAreaCode == ""{
	//	instance.ShipInfo = &ShipInfo{}
	//}else{
	//	// area := eel.NewAreaService().GetAreaByCode(model.ShipAreaCode)
	//	instance.ShipInfo = &ShipInfo{
	//		Name: model.ShipName,
	//		Phone: model.ShipPhone,
	//		Address: model.ShipAddress,
	//		AreaCode: model.ShipAreaCode,
	//	}
	//}

	instance.Products = make([]*OrderProduct, 0)
	
	return instance
}

func init() {
	enableSyncClearance, _ = config.ServiceConfig.Bool("order::ENABLE_SYNC_CLEARANCE")
}