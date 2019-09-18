package product

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type DisabledCategory struct {
	eel.RestResource
}

func (this *DisabledCategory) Resource() string {
	return "product.disabled_category"
}

func (this *DisabledCategory) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"id:int"},
		"DELETE": []string{"id:int"},
	}
}

func (this *DisabledCategory) Put(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	repository := product.NewProductCategoryRepository(bCtx)
	productCategory := repository.GetProductCategoryInCorp(corp, id)
	
	if productCategory == nil {
		ctx.Response.Error("product_category:invalid_category", fmt.Sprintf("id=%d", id))
		return
	}
	
	productCategory.Disable()
	
	ctx.Response.JSON(eel.Map{})
}

func (this *DisabledCategory) Delete(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	repository := product.NewProductCategoryRepository(bCtx)
	productCategory := repository.GetProductCategoryInCorp(corp, id)
	
	if productCategory == nil {
		ctx.Response.Error("product_category:invalid_category", fmt.Sprintf("id=%d", id))
		return
	}
	
	productCategory.Enable()
	
	ctx.Response.JSON(eel.Map{})
}
