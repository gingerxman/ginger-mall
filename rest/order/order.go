package order

import (
	"encoding/json"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_order "github.com/gingerxman/ginger-mall/business/order"
	"github.com/gingerxman/ginger-mall/business/order/resource"
)

type Order struct {
	eel.RestResource
}

func (this *Order) Resource() string {
	return "order.order"
}

func (this *Order) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"bid", "?with_options:json"},
		"PUT": []string{
			"?corp_id:int",
			"products:json-array",
			"?ship_info:json",
			"?custom_order_type",
			"?callback_resource",
			"?biz_code",
			"?imoney_usages:json-array",
			"?shopping_cart_item_ids:json-array",
			"?coupon_usage:json",
			"?settlement_rule:json",
			"?salesman_id:int",
			"?extra_data:json",
			"?no_resettle:bool",
		},
	}
}

func (this *Order) GetLockKey() string {
	// TODO: 实现GetLockKey逻辑
	//method := this.Ctx.Input.Method()
	//if method == "GET" {
	//	return ""
	//} else {
	//	bCtx := req.GetBusinessContext()
	//	user := account.GetUserFromContext(bCtx)
	//
	//	return fmt.Sprintf("order_create_%d", user.GetId())
	//}
	return ""
}

func (this *Order) Get(ctx *eel.Context) {
	//get order
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	bid := req.GetString("bid")
	order := b_order.NewOrderRepository(bCtx).GetOrderByBid(bid)
	
	//处理fill options
	fillOptions := req.GetJSON("with_options")
	fillService := b_order.NewFillOrderService(bCtx)
	if fillOptions == nil{
		fillOptions = eel.Map{}
	}
	//处理invoice的fill options
	invoiceFillOptions := eel.Map{
		"with_products": true,
	}
	//if option, ok := fillOptions["with_settlements"]; ok {
	//	invoiceFillOptions["with_settlements"] = option
	//} else {
	//	invoiceFillOptions["with_settlements"] = true
	//}
	fillOptions["with_invoice"] = invoiceFillOptions
	fillService.Fill([]*b_order.Order{order}, fillOptions)
	
	//encode
	data := b_order.NewEncodeOrderService(bCtx).Encode(order)
	
	ctx.Response.JSON(data)
}

func (this *Order) Put(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	purchaseInfo := this.parsePurchaseInfo(ctx)
	
	order, err := b_order.NewOrderFactory(bCtx).CreateOrder(purchaseInfo)
	if err != nil {
		ctx.Response.Error(err.Error(), err.Error())
	} else {
		ctx.Response.JSON(eel.Map{
			"id": order.Id,
			"bid": order.Bid,
			"money": order.Money.FinalMoney,
			"status": order.GetStatusText(),
		})
	}
}

func (this *Order) parseShipInfo(ctx *eel.Context) *b_order.ShipInfo {
	req := ctx.Request
	strShipInfo := req.GetString("ship_info")
	shipInfo := b_order.ShipInfo{}
	
	if strShipInfo != "" {
		err := json.Unmarshal([]byte(strShipInfo), &shipInfo)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("order:parse_ship_info_fail", "解析ShipInfo出错"))
		}
	}
	
	return &shipInfo
}

func (this *Order) parseCouponUsage(ctx *eel.Context) *resource.CouponUsage {
	req := ctx.Request
	strCouponUsage := req.GetString("coupon_usage", "")
	if strCouponUsage == "" {
		return nil
	}
	couponUsage := resource.CouponUsage{}
	
	if strCouponUsage != "" {
		err := json.Unmarshal([]byte(strCouponUsage), &couponUsage)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("order:parse_coupon_usage_fail", "解析CouponUsage出错"))
		}
	}
	
	return &couponUsage
}

// parseExtraData extra_data中的参数做兼容
// callback_resource、source_service、settlement_rule、no_resettle应该放在extra_data中
// 如果外部有参数，则使用外部参数覆盖extra_data中的对应值
func (this *Order) parseExtraData(ctx *eel.Context) map[string]interface{}{
	req := ctx.Request
	extraData := req.GetJSON("extra_data")
	if extraData == nil{
		extraData = make(map[string]interface{})
	}
	callbackResource := req.GetString("callback_resource", "")
	if callbackResource != ""{
		extraData["callback_resource"] = callbackResource
	}
	bizCode := req.GetString("biz_code", "")
	if bizCode != ""{
		extraData["source_service"] = bizCode
	}

	settlementRule := req.GetJSON("settlement_rule")
	if settlementRule != nil && len(settlementRule)>0{
		extraData["settlement_rule"] = settlementRule
	}
	noResettle, _ := req.GetBool("no_resettle", false)
	if _, ok := extraData["no_resettle"]; !ok{
		extraData["no_resettle"] = noResettle
	}

	return extraData
}

func (this *Order) parsePurchaseInfo(ctx *eel.Context) *b_order.PurchaseInfo {
	purchaseInfo := &b_order.PurchaseInfo{}
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	
	//确定corpId
	corp := account.GetCorpFromContext(bCtx)
	purchaseInfo.CorpId = corp.GetId()
	
	extraData := this.parseExtraData(ctx)
	bizCode := req.GetString("biz_code", "")
	if bizCode == ""{
		if v, ok := extraData["source_service"]; ok && v != nil{
			bizCode = v.(string)
		}else{
			bizCode = "app"
		}
	}
	purchaseInfo.CustomerMessage = req.GetString("message", "")
	purchaseInfo.OrderType = req.GetString("custom_order_type", "product")
	purchaseInfo.ShipInfo = this.parseShipInfo(ctx)
	purchaseInfo.CouponUsage = this.parseCouponUsage(ctx)
	purchaseInfo.ExtraData = extraData
	purchaseInfo.BizCode = bizCode
	purchaseInfo.SalesmanId, _ = req.GetInt("salesman_id", 0)
	purchaseInfo.ShoppingCartItemIds = req.GetIntArray("shopping_cart_item_ids")
	
	productsArray := req.GetJSONArray("products")
	imoneysArray := req.GetJSONArray("imoney_usages")
	purchaseInfo.Resources = resource.NewParseResourceService(bCtx).Parse(purchaseInfo.SalesmanId, productsArray, imoneysArray, purchaseInfo.CouponUsage)
	
	return purchaseInfo
}
