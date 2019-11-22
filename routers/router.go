package routers

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/handler/rest/console"
	"github.com/gingerxman/eel/handler/rest/op"
	"github.com/gingerxman/ginger-mall/rest/area"
	"github.com/gingerxman/ginger-mall/rest/consumption"
	"github.com/gingerxman/ginger-mall/rest/dev"
	"github.com/gingerxman/ginger-mall/rest/mall"
	"github.com/gingerxman/ginger-mall/rest/mall/ship_info"
	"github.com/gingerxman/ginger-mall/rest/mall/shopping_cart"
	"github.com/gingerxman/ginger-mall/rest/material"
	"github.com/gingerxman/ginger-mall/rest/order"
	"github.com/gingerxman/ginger-mall/rest/product"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&op.Health{})
	
	/*
	 product
	 */
	//category
	eel.RegisterResource(&product.Category{})
	eel.RegisterResource(&product.DisabledCategory{})
	eel.RegisterResource(&product.SubCategories{})
	//label
	eel.RegisterResource(&product.ProductLabel{})
	eel.RegisterResource(&product.ProductLabels{})
	eel.RegisterResource(&product.CorpProductLabels{})
	eel.RegisterResource(&product.DisabledCategory{})
	//property
	eel.RegisterResource(&product.ProductProperty{})
	eel.RegisterResource(&product.ProductPropertyValue{})
	eel.RegisterResource(&product.ProductProperties{})
	eel.RegisterResource(&product.CorpProductProperties{})
	//product
	eel.RegisterResource(&product.Product{})
	eel.RegisterResource(&product.OffshelfProducts{})
	eel.RegisterResource(&product.OnshelfProducts{})
	eel.RegisterResource(&product.CorpProducts{})
	eel.RegisterResource(&product.CreateOptions{})
	
	/*
	 order
	 */
	eel.RegisterResource(&order.Order{})
	eel.RegisterResource(&order.PayedOrder{})
	eel.RegisterResource(&order.CanceledOrder{})
	eel.RegisterResource(&order.ConfirmedInvoice{})
	eel.RegisterResource(&order.CanceledInvoice{})
	eel.RegisterResource(&order.ShippedInvoice{})
	eel.RegisterResource(&order.FinishedInvoice{})
	eel.RegisterResource(&order.Orders{})
	eel.RegisterResource(&order.CorpInvoices{})
	eel.RegisterResource(&order.UserOrders{})
	eel.RegisterResource(&order.OrderStatus{})
	eel.RegisterResource(&order.OrderRemark{})
	
	/*
	 consumption
	 */
	eel.RegisterResource(&consumption.UserConsumptionRecords{})
	
	/*
	 mall
	 */
	eel.RegisterResource(&mall.SubCategories{})
	eel.RegisterResource(&mall.Products{})
	eel.RegisterResource(&mall.Product{})
	eel.RegisterResource(&mall.PurchaseData{})
	//ship_info
	eel.RegisterResource(&ship_info.ShipInfo{})
	eel.RegisterResource(&ship_info.ShipInfos{})
	eel.RegisterResource(&ship_info.DefaultShipInfo{})
	//shopping_cart
	eel.RegisterResource(&shopping_cart.ShoppingCartItem{})
	eel.RegisterResource(&shopping_cart.ShoppingCart{})
	eel.RegisterResource(&shopping_cart.ProductCount{})
	//order
	
	/*
	 material
	*/
	eel.RegisterResource(&material.Image{})
	
	/*
	 area
	 */
	eel.RegisterResource(&area.Area{})
	eel.RegisterResource(&area.AreaCode{})
	eel.RegisterResource(&area.YouzanAreaList{})

	/*
	 dev
	 */
	eel.RegisterResource(&dev.BDDReset{})
}