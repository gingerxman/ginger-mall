package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/eel/event"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/gorm"
	
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/events"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"strings"
	"time"
)

type orderMoneyInfo struct {
	Postage float64 //运费
	FinalMoney int //订单金额
	EditMoney float64 //修改金额
	PayMoney int //支付金额
	ProductPrice int //商品总价
}

type Order struct {
	eel.EntityBase
	
	Id int
	OriginalOrderId int
	Bid string
	Type string
	CustomType string
	Status int

	BizCode string //业务码
	IsCleared bool //是否清算
	
	Money *orderMoneyInfo
	
	UserId int
	User *account.User
	
	CorpId int
	Corp *account.Corp
	
	Remark string //备注
	Message string //消费者留言
	CancelReason string //取消订单的原因
	
	IsDeleted bool
	
	Invoices []*Invoice
	
	Resources string
	ExtraData string
	
	PaymentTime time.Time
	PaymentType string
	CreatedAt time.Time
	
	// 日志
	OperationLogs []*OrderOperationLog
	StatusLogs []*OrderStatusLog
}

func (this *Order) GetId() int {
	return this.Id
}

func (this *Order) GetBid() string {
	return this.Bid
}

func (this *Order) GetDeductableMoney() int {
	return 0
}

func (this *Order) GetStatusText() string{
	if text, ok := m_order.STATUS2STR[this.Status]; ok{
		return text
	}else{
		return "未知"
	}
}

func (this *Order) GetResources() []map[string]interface{}{
	//var resources []map[string]interface{}
	//err := json.Unmarshal([]byte(this.Resources), &resources)
	//if err != nil {
	//	eel.Logger.Error(err)
	//}
	resources := make([]map[string]interface{}, 0)
	js, err := simplejson.NewJson([]byte(this.Resources))
	if err != nil {
		eel.Logger.Error(err)
	} else {
		data, err := js.Array()
		if err != nil {
			eel.Logger.Error(err)
		} else {
			for _, d := range data{
				resources = append(resources, d.(map[string]interface{}))
			}
		}
	}

	return resources
}

func (this *Order) GetExtraData() map[string]interface{}{
	var data map[string]interface{}
	err := json.Unmarshal([]byte(this.ExtraData), &data)
	if err != nil {
		eel.Logger.Error(err)
	}
	return data
}

func (this *Order) GetCallbackResource() *OrderCallbackResource{
	extraData := this.GetExtraData()
	if extraData != nil{
		if str, ok := extraData["callback_resource"]; ok{
			return NewOrderCallbackResource(str.(string))
		}
	}
	return nil
}

func (this *Order) IsInvoice() bool {
	return this.Type == m_order.ORDERTYPE2STR[m_order.ORDER_TYPE_PRODUCT_INVOICE]
}

func (this *Order) IsCustomOrder() bool {
	return this.Type == m_order.ORDERTYPE2STR[m_order.ORDER_TYPE_CUSTOM]
}

// IsDepositOrder 是否充值订单
func (this *Order) IsDepositOrder() bool{
	return this.CustomType == "imoney:deposit"
}

func (this *Order) AddInvoice(invoice *Invoice) {
	this.Invoices = append(this.Invoices, invoice)
}

// GetCorpRelatedUserId 获取订单corp的关联user_id
func (this *Order) GetCorpRelatedUserId() int{
	if this.CorpId == 0{
		return 0
	}
	resp, err := eel.NewResource(this.Ctx).Get("gskep", "corp.corp", eel.Map{
		"id": this.CorpId,
	})

	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order:get_corp_related_user_failed", "获取corp关联的user_id失败"))
	}

	if !resp.IsSuccess(){
		errCode, _ := resp.RespData.Get("errCode").String()
		panic(eel.NewBusinessError(errCode, "获取corp关联的user_id失败"))
	}

	data, _ := resp.Data().Map()
	corpData := data["corp"].(map[string]interface{})
	id, _ := corpData["related_user"].(map[string]interface{})["id"].(json.Number).Int64()
	return int(id)
}

//IsFinished 订单是否已结束
func (this *Order) IsFinished() bool {
	return this.Status == m_order.ORDER_STATUS_SUCCESSED
}

func (this *Order) save() {
	ormParams := gorm.Params{
		"status": this.Status,
		"remark": this.Remark,
		"final_money": this.Money.FinalMoney,
		"update_at": time.Now(),
	}
	
	if !this.PaymentTime.IsZero() {
		ormParams["payment_time"] = this.PaymentTime
		ormParams["payment_type"] = this.PaymentType
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{}).Where("id", this.Id).Update(ormParams)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order:save_failed", "保存数据失败"))
	}
}

func (this *Order) updateInvoices(params gorm.Params){
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{}).Where("original_order_id", this.Id).Update(params)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order:update_invoices_failed", "保存数据失败"))
	}
}

func (this *Order) UpdateFinalMoney(finalMoney int) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{}).Where("Id", this.Id).Update(gorm.Params{
		"final_money": finalMoney,
	})
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(err)
	}
}

var autoFinishOrderSet = map[string]bool {
	"honeycomb:dating": true,
	"honeycomb:gift": true,
	"honeycomb:broadcast_blog": true,
	"mpcoin_product": true,
	"virtual_product": true,
}
func (this *Order) isAutoFinishOrder() bool {
	if this.CustomType == "" {
		return false
	}
	
	if _, ok := autoFinishOrderSet[this.CustomType]; ok {
		return true
	}
	
	if strings.Contains(this.CustomType, "order:auto_finish") {
		return true
	}
	
	return false
}

//Pay 对订单进行支付
// args[0] payment_type 支付通道
func (this *Order) Pay(args ...string) {
	switch this.Status {
	case m_order.ORDER_STATUS_CANCEL, m_order.ORDER_STATUS_PAYING, m_order.ORDER_STATUS_WAIT_PAY:
	default:
		panic(eel.NewBusinessError("order: payment_failed", fmt.Sprintf("订单状态不对(%s)", this.GetStatusText())))
	}

	this.PaymentType = "weixin"
	switch len(args) {
	case 1:
		this.PaymentType = args[0]
	}

	targetStatus := m_order.ORDER_STATUS_NONSENSE
	this.PaymentTime = time.Now()
	
	//更改invoice的状态
	this.updateInvoices(gorm.Params{
		"type": m_order.ORDER_TYPE_PRODUCT_INVOICE,
		"status": m_order.ORDER_STATUS_PAYED_NOT_SHIP,
		"payment_time": this.PaymentTime,
		"payment_type": this.PaymentType,
	})

	//更改order的状态
	this.Status = targetStatus
	this.save()
	NewOrderPaidService(this.Ctx).AfterPaid(this)
	
	//如果是自动完成订单，则进行自动完成的检查
	if this.isAutoFinishOrder() {
		//因为自动完成的订单都是虚拟商品，order与invoice合一的，所以可以直接强制转换为invoice
		invoice := NewInvoiceFromOrder(this.Ctx, this)
		invoice.ForceFinish()
	}
	
	// 记录操作日志
	orderLogService := NewOrderLogService(this.Ctx)
	orderLogService.LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_PAY,
	})
}

// Cancellable 是否可删除
func (this *Order) Cancellable() bool{
	switch this.Status {
	case m_order.ORDER_STATUS_WAIT_PAY, m_order.ORDER_STATUS_PAYING:
		return true
	}
	return false
}

// Refundable 是否可退款
func (this *Order) Refundable() bool{
	switch this.Status {
	case m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM, m_order.ORDER_STATUS_PAYED_NOT_SHIP,
			m_order.ORDER_STATUS_PAYED_SHIPED, m_order.ORDER_STATUS_SUCCESSED:
		return true
	}
	return false
}

// Refund 订单退款
func (this *Order) Refund(){
	if ! this.Refundable(){
		return
	}
	// 对所有出货单进行取消操作
	for _, invoice := range this.Invoices{
		invoice.FinishRefund()
	}

	this.Status = m_order.ORDER_STATUS_REFUNDED
	this.save()

	// 记录操作日志
	NewOrderLogService(this.Ctx).LogOperation(&OperationData{
		Order: this,
		Action: ORDER_OPERATION_REFUND,
	})

	NewOrderRefundedService(this.Ctx).AfterRefunded(this)
}

// Cancel 取消订单
func (this *Order) Cancel(reason string){
	if this.Cancellable(){
		//for _, invoince := range this.Invoices{
		//	resources := resource.NewParseResourceService(this.Ctx).ParseFromOrderResources(invoince.GetResources())
		//	resource.NewAllocateResourceService(this.Ctx).Release(resources)
		//}
		
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

		// 对所有出货单进行取消操作
		invoices := NewOrderRepository(this.Ctx).GetInvoicesForOrder(this.Id)
		for _, invoice := range invoices {
			invoice.Cancel(reason)
		}
		
		// 异步消息
		event.AsyncEvent.Send(events.ORDER_CANCLLED, map[string]interface{}{
			"bid": this.Bid,
		})
	}
}

// SetCallbackStatus 设置回调处理结果
func (this *Order) SetCallbackStatus(succeed bool){
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{}).Where("Bid", this.Bid).Update(gorm.Params{
		"CallbackSucceed": succeed,
	})
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order:set_callback_status_failed", "保存数据失败"))
	}
}

func (this *Order) finish(){
	if this.Status == m_order.ORDER_STATUS_SUCCESSED{
		return
	}
	handler := NewFinishedOrderHandler(this.Ctx)
	handler.DoSettlement(this)
	// 处理异步回调
	handler.DoCallback(this)
	// 异步消息
	event.AsyncEvent.Send(events.ORDER_FINISHED, map[string]interface{}{
		"bid": this.Bid,
	})
}

func (this *Order) ChangeStatusToNonsense() {
	this.Status = m_order.ORDER_STATUS_NONSENSE
	this.save()
}

//ShouldAutoPay 是否可以自动完成
//以下情况可以自动完成：
// 1. final_money = 0
func (this *Order) ShouldAutoPay() bool {
	return this.Money.FinalMoney == 0
}

func NewOrderFromModel(ctx context.Context, model *m_order.Order) *Order {
	instance := &Order{}
	
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.Bid = model.Bid
	instance.OriginalOrderId = model.OriginalOrderId
	instance.Type = m_order.ORDERTYPE2STR[model.Type]
	instance.CustomType = model.CustomType
	instance.CorpId = model.CorpId
	instance.BizCode = model.BizCode
	instance.UserId = model.UserId
	instance.Status = model.Status
	instance.IsCleared = model.IsCleared
	instance.IsDeleted = model.IsDeleted
	instance.Remark = model.Remark
	instance.PaymentTime = model.PaymentTime
	instance.PaymentType = model.PaymentType
	instance.CreatedAt = model.CreatedAt
	instance.Resources = model.Resources
	instance.ExtraData = model.ExtraData
	instance.Message = model.CustomerMessage
	instance.CancelReason = model.CancelReason
	
	instance.Money = &orderMoneyInfo{
		Postage: model.Postage,
		FinalMoney: model.FinalMoney,
	}
	
	return instance
}

func init() {
}
