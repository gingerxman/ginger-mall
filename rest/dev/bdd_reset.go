package dev

import (
	"github.com/gingerxman/eel"
)

type BDDReset struct {
	eel.RestResource
}

func (this *BDDReset) Resource() string {
	return "dev.bdd_reset"
}

func (this *BDDReset) SkipAuthCheck() bool {
	return true
}

func (r *BDDReset) IsForDevTest() bool {
	return true
}

func (this *BDDReset) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT":  []string{},
	}
}

func (this *BDDReset) Put(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	o := eel.GetOrmFromContext(bCtx)
	
	o.Exec("delete from mall_shopping_cart")
	o.Exec("delete from mall_ship_info")
	
	o.Exec("delete from order_user_consumption_record")
	o.Exec("delete from order_has_product")
	o.Exec("delete from order_has_logistics")
	o.Exec("delete from order_operation_log")
	o.Exec("delete from order_status_log")
	o.Exec("delete from order_order")
	
	o.Exec("delete from product_has_label")
	o.Exec("delete from product_label")
	o.Exec("delete from product_sku_has_property")
	o.Exec("delete from product_sku_property_value")
	o.Exec("delete from product_sku_property")
	o.Exec("delete from product_category_has_product")
	o.Exec("delete from product_category")
	o.Exec("delete from product_media")
	o.Exec("delete from product_usable_imoney")
	o.Exec("delete from product_description")
	o.Exec("delete from product_sku")
	o.Exec("delete from product_logistics")
	o.Exec("delete from product_pool_product")
	o.Exec("delete from product_product")
	
	ctx.Response.JSON(eel.Map{})
}

