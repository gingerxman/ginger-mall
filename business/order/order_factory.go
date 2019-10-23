package order

import (
	"context"
	"fmt"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/order/resource"
)

//var snowflakeNode, _ = snowflake.NewNode(1)

type OrderFactory struct {
	eel.ServiceBase
}

func NewOrderFactory(ctx context.Context) *OrderFactory {
	service := new(OrderFactory)
	service.Ctx = ctx
	return service
}

func (this *OrderFactory) CreateOrder(purchaseInfo *PurchaseInfo) (*Order, error) {
	resourceManager := resource.NewResourceManager(this.Ctx)
	resourceManager.AddResources(purchaseInfo.Resources)
	err := resourceManager.Validate()
	if err != nil {
		panic(eel.NewBusinessError(fmt.Sprintf("create_order_fail:%s", err.Error()), err.Error()))
	}
	
	//申请锁
	err = resourceManager.Lock()
	if err != nil {
		eel.Logger.Error(err)
		return nil, err
	}
	defer resourceManager.Unlock()
	
	//创建NewOrder
	newOrder := GenerateNewOrder(this.Ctx, purchaseInfo, resourceManager)
	
	//申请资源
	err = resourceManager.AllocateForOrder(newOrder)
	if err != nil {
		eel.Logger.Error(err)
		return nil, err
	}
	
	//如果没有成功创建订单，则释放已申请的资源
	isSuccess := false
	defer func() {
		if !isSuccess {
			eel.Logger.Error("Failed to create Order object or save the order! Release all resources.")
			resourceManager.ReleaseAllocatedResources()
		}
	}()
	
	//创建订单
	order, err := newOrder.Save()
	if err != nil {
		eel.Logger.Error(err)
		return nil, err
	}
	isSuccess = true
	
	if order.ShouldAutoPay() {
		eel.Logger.Debug("[order] auto pay order")
		order.Pay()
	}
	
	return order, nil
}

//func (this *OrderFactory) generateBid() string {
//	result := snowflakeNode.Generate().String()
//	return result
//}

//func (this *OrderFactory) prepareExtraData(purchaseInfo *PurchaseInfo){
//	extraData := purchaseInfo.ExtraData
//	if extraData == nil{
//		extraData = make(map[string]interface{})
//		purchaseInfo.ExtraData = extraData
//	}
//	extraData["biz_code"] = purchaseInfo.BizCode
//	if purchaseInfo.CallbackResourceUrl != ""{
//		extraData["callback_resource"] = purchaseInfo.CallbackResourceUrl
//	}
//	if purchaseInfo.SettlementRule != nil{
//		extraData["settlement_rule"] = purchaseInfo.SettlementRule
//	}
//}
//
//func (this *OrderFactory) saveOrderInDb(saveType string, resources []business.IResource, purchaseInfo *PurchaseInfo, moneyInfo *orderMoneyInfo, parentOrder *Order) *Order {
//	user := account.GetUserFromContext(this.Ctx)
//
//	model := &m_order.Order{}
//	model.PlatformCorpId = 0
//	model.UserId = user.GetId()
//	model.CorpId = purchaseInfo.CorpId
//	model.SupplierId = 0
//	model.Status = m_order.ORDER_STATUS_NOT
//	model.CustomerMessage = purchaseInfo.CustomerMessage
//	model.Bid = this.generateBid()
//	model.BizCode = purchaseInfo.BizCode
//
//	if saveType == "order" {
//		model.OriginalOrderId = 0
//		model.SupplierId = 0
//	} else {
//		model.OriginalOrderId = parentOrder.Id
//		model.SupplierId = resources[0].(*resource.ProductResource).GetPoolProduct().SupplierId
//	}
//
//	shipInfo := purchaseInfo.ShipInfo
//	model.ShipName = shipInfo.Name
//	model.ShipPhone = shipInfo.Phone
//	model.ShipAddress = shipInfo.Address
//	model.ShipAreaCode = shipInfo.GetAreaCode()
//
//	if purchaseInfo.OrderType == "product" {
//		model.Type = m_order.ORDER_TYPE_PRODUCT_ORDER
//		model.CustomType = "product"
//	} else {
//		model.Type = m_order.ORDER_TYPE_CUSTOM
//		model.CustomType = purchaseInfo.OrderType
//	}
//
//	if saveType == "invoice" {
//		model.Type = m_order.ORDER_TYPE_PRODUCT_INVOICE
//	}
//
//	//资金
//	model.Postage = moneyInfo.Postage
//	model.FinalMoney = moneyInfo.FinalMoney
//
//	this.prepareExtraData(purchaseInfo)
//
//	extraData, _ := json.Marshal(purchaseInfo.ExtraData)
//	model.ExtraData = string(extraData)
//
//	o := eel.GetOrmFromContext(this.Ctx)
//	id, err := o.Insert(model)
//	if err != nil {
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("order:save_fail", "保存订单数据出错"))
//	}
//
//	model.Id = int(id)
//	order := NewOrderFromModel(this.Ctx, model)
//
//	//处理订单资源
//	resourceDatas := make([]map[string]interface{}, 0)
//	for _, res := range resources {
//		err := res.SaveForOrder(order)
//		if err != nil {
//			eel.Logger.Error(err)
//		}
//		resourceDatas = append(resourceDatas, res.ToMap())
//	}
//
//	//保存resources
//	bytes, err := json.Marshal(resourceDatas)
//	if err != nil {
//		eel.Logger.Error(err)
//	} else {
//		_, err := o.Model(&m_order.Order{}).Where("id", order.Id).Update(gorm.Params{
//			"resources": string(bytes),
//		})
//
//		if err != nil {
//			eel.Logger.Error(err)
//			panic(eel.NewBusinessError("order:save_resource_fail", "保存订单数据的resources出错"))
//		}
//	}
//
//	return order
//}
//
//func (this *OrderFactory) saveOrder(purchaseInfo *PurchaseInfo, resourceManager *resource.ResourceManager) *Order {
//	isCustomOrderType := (purchaseInfo.OrderType != "product")
//
//	var order *Order
//	calculateOrderMoneyService := NewCalculateOrderMoneyService(this.Ctx)
//	if isCustomOrderType {
//		moneyInfo := calculateOrderMoneyService.Calculate(resourceManager.Resources, purchaseInfo)
//		order = this.saveOrderInDb("order", resourceManager.Resources, purchaseInfo, moneyInfo,nil)
//	} else {
//		moneyInfo := calculateOrderMoneyService.Calculate(resourceManager.Resources, purchaseInfo)
//		order = this.saveOrderInDb("order", resourceManager.GetNonProductResources(), purchaseInfo, moneyInfo,nil)
//
//		finalMoney := 0.0
//		for _, resourceGroup := range resourceManager.GroupResourceBySupplier() {
//			moneyInfo = calculateOrderMoneyService.Calculate(resourceGroup.Resources, purchaseInfo)
//			finalMoney += moneyInfo.FinalMoney
//			this.saveOrderInDb("invoice", resourceGroup.Resources, purchaseInfo, moneyInfo, order)
//		}
//		order.UpdateFinalMoney(finalMoney)
//	}
//
//	//如果是从购物车发起的购买行为，购买成功后需要删除购物车项
//	if purchaseInfo.IsFromShoppingCart() {
//		shopping_cart.NewShoppingCartService(this.Ctx).DeleteShoppingCartItems(purchaseInfo.ShoppingCartItemIds)
//	}
//	return order
//}

func init() {
}
