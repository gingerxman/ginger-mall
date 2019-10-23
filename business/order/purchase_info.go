package order

import (
	"fmt"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/order/resource"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

type shipAreaItem struct {
	Id int
	Name string
}

type shipArea struct {
	Province shipAreaItem
	City shipAreaItem
	District shipAreaItem
}

type ShipInfo struct {
	Phone string //收货电话
	Address string //详细地址
	Name string //收件人
	Area shipArea
}


func (this *ShipInfo) GetAreaCode() string {
	return fmt.Sprintf("%d_%d_%d", this.Area.Province.Id, this.Area.City.Id, this.Area.District.Id)
}

func (this *ShipInfo) IsValid() bool {
	return this.Area.Province.Id != 0
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
