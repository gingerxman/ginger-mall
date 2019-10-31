package order

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/order/resource"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type ShipInfo struct {
	Phone string `json:"phone"`//收货电话
	Address string `json:"address"`//详细地址
	Name string `json:"name"`//收件人
	AreaCode string `json:"area_code"`
}

func (this *ShipInfo) GetArea() *eel.Area {
	return eel.NewAreaService().GetAreaByCode(this.AreaCode)
}

func (this *ShipInfo) IsValid() bool {
	return true
}

type PurchaseInfo struct {
	User *account.User
	Resources []business.IResource
	ShipInfo *ShipInfo
	CouponUsage *resource.CouponUsage
	CustomerMessage string
	OrderType string
	CorpId int
	BizCode string
	SalesmanId int
	ShoppingCartItemIds []int
	ExtraData map[string]interface{}
}

func (this *PurchaseInfo) Check() error {
	return nil
}

func (this *PurchaseInfo) IsCustomTypeOrder() bool {
	return this.OrderType == m_order.ORDERTYPE2STR[m_order.ORDER_TYPE_CUSTOM]
}

func (this *PurchaseInfo) IsFromShoppingCart() bool {
	return len(this.ShoppingCartItemIds) > 0
}

//IsFromSalesman 是否是通过分销员分销进行购买
func (this *PurchaseInfo) IsFromSalesman() bool {
	return this.SalesmanId > 0
}

func init() {
}
