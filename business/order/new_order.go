package order

import (
	"encoding/json"
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/snowflake"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/shopping_cart"
	"github.com/gingerxman/ginger-mall/business/order/resource"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"context"
)

var snowflakeNode, _ = snowflake.NewNode(1)

type NewOrder struct {
	eel.EntityBase
	Id int
	Bid string
	deductableMoney float64
	
	purchaseInfo *PurchaseInfo
	resourceManager *resource.ResourceManager
}

func (this *NewOrder) GetId() int {
	return 0
}

func (this *NewOrder) GetBid() string {
	return this.Bid
}

func (this *NewOrder) GetDeductableMoney() float64 {
	if this.deductableMoney == 0 {
	
	}
	
	return this.deductableMoney
}

func (this *NewOrder) generateBid() string {
	result := snowflakeNode.Generate().String()
	return result
}

func (this *NewOrder) AssignBid() {
	this.Bid = this.generateBid()
}

func (this *NewOrder) prepareExtraData(purchaseInfo *PurchaseInfo){
	extraData := purchaseInfo.ExtraData
	if extraData == nil{
		extraData = make(map[string]interface{})
		purchaseInfo.ExtraData = extraData
	}
}

func (this *NewOrder) saveOrderInDb(saveType string, resources []business.IResource, purchaseInfo *PurchaseInfo, moneyInfo *orderMoneyInfo, parentOrder *Order) *Order {
	user := account.GetUserFromContext(this.Ctx)
	
	model := &m_order.Order{}
	model.PlatformCorpId = 0
	model.UserId = user.GetId()
	model.CorpId = purchaseInfo.CorpId
	model.SupplierId = 0
	model.Status = m_order.ORDER_STATUS_WAIT_PAY
	model.CustomerMessage = purchaseInfo.CustomerMessage
	model.BizCode = purchaseInfo.BizCode
	
	if saveType == "order" {
		model.OriginalOrderId = 0
		model.SupplierId = 0
		if this.Bid == "" {
			panic("empty bid")
		}
		model.Bid = this.Bid
	} else {
		model.Bid = this.generateBid()
		model.OriginalOrderId = parentOrder.Id
		model.SupplierId = resources[0].(*resource.ProductResource).GetPoolProduct().SupplierId
	}
	
	shipInfo := purchaseInfo.ShipInfo
	model.ShipName = shipInfo.Name
	model.ShipPhone = shipInfo.Phone
	model.ShipAddress = shipInfo.Address
	model.ShipAreaCode = shipInfo.GetAreaCode()
	
	if purchaseInfo.OrderType == "product" {
		model.Type = m_order.ORDER_TYPE_PRODUCT_ORDER
		model.CustomType = "product"
	} else {
		model.Type = m_order.ORDER_TYPE_CUSTOM
		model.CustomType = purchaseInfo.OrderType
	}
	
	if saveType == "invoice" {
		model.Type = m_order.ORDER_TYPE_PRODUCT_INVOICE
	}
	
	//资金
	model.Postage = moneyInfo.Postage
	model.FinalMoney = moneyInfo.FinalMoney
	
	this.prepareExtraData(purchaseInfo)
	
	extraData, _ := json.Marshal(purchaseInfo.ExtraData)
	model.ExtraData = string(extraData)
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Create(model)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order:save_fail", "保存订单数据出错"))
	}
	
	order := NewOrderFromModel(this.Ctx, model)
	
	//处理订单资源
	resourceDatas := make([]map[string]interface{}, 0)
	for _, res := range resources {
		err := res.SaveForOrder(order)
		if err != nil {
			eel.Logger.Error(err)
		}
		resourceDatas = append(resourceDatas, res.ToMap())
	}
	
	//保存resources
	bytes, err := json.Marshal(resourceDatas)
	if err != nil {
		eel.Logger.Error(err)
	} else {
		db := o.Model(&m_order.Order{}).Where("id", order.Id).Update(gorm.Params{
			"resources": string(bytes),
		})
		err := db.Error
		
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("order:save_resource_fail", "保存订单数据的resources出错"))
		}
	}
	
	return order
}

func (this *NewOrder) Save() (*Order, error) {
	purchaseInfo := this.purchaseInfo
	resourceManager := this.resourceManager
	
	isCustomOrderType := (purchaseInfo.OrderType != "product")
	
	var order *Order
	calculateOrderMoneyService := NewCalculateOrderMoneyService(this.Ctx)
	if isCustomOrderType {
		moneyInfo := calculateOrderMoneyService.Calculate(resourceManager.Resources, purchaseInfo, this)
		order = this.saveOrderInDb("order", resourceManager.Resources, purchaseInfo, moneyInfo,nil)
	} else {
		moneyInfo := calculateOrderMoneyService.Calculate(resourceManager.Resources, purchaseInfo, this)
		order = this.saveOrderInDb("order", resourceManager.GetNonProductResources(), purchaseInfo, moneyInfo,nil)
		
		finalMoney := 0.0
		for _, resourceGroup := range resourceManager.GroupResourceBySupplier() {
			moneyInfo = calculateOrderMoneyService.Calculate(resourceGroup.Resources, purchaseInfo, this)
			finalMoney += moneyInfo.FinalMoney
			this.saveOrderInDb("invoice", resourceGroup.Resources, purchaseInfo, moneyInfo, order)
		}
		order.UpdateFinalMoney(finalMoney)
	}
	
	//如果是从购物车发起的购买行为，购买成功后需要删除购物车项
	if purchaseInfo.IsFromShoppingCart() {
		shopping_cart.NewShoppingCartService(this.Ctx).DeleteShoppingCartItems(purchaseInfo.ShoppingCartItemIds)
	}
	return order, nil
}

func GenerateNewOrder(ctx context.Context, purchaseInfo *PurchaseInfo, resourceManager *resource.ResourceManager) *NewOrder {
	newOrder := &NewOrder{}
	newOrder.Ctx = ctx
	newOrder.purchaseInfo = purchaseInfo
	newOrder.resourceManager = resourceManager
	newOrder.deductableMoney = 0
	newOrder.AssignBid()
	
	return newOrder
}

func init() {
}
