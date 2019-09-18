package product

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type DisabledProductLabel struct {
	eel.RestResource
}

func (this *DisabledProductLabel) Resource() string {
	return "product.disabled_label"
}

func (this *DisabledProductLabel) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"id:int"},
		"DELETE": []string{"id:int"},
	}
}

func (this *DisabledProductLabel) Put(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	repository := product.NewProductLabelRepository(bCtx)
	productLabel := repository.GetProductLabelInCorp(corp, id)
	if productLabel == nil {
		ctx.Response.Error("product_label:invalid_label", fmt.Sprintf("id=%d", id))
		return
	}
	
	productLabel.Disable()
	
	ctx.Response.JSON(eel.Map{})
}

func (this *DisabledProductLabel) Delete(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	repository := product.NewProductLabelRepository(bCtx)
	productLabel := repository.GetProductLabelInCorp(corp, id)
	if productLabel == nil {
		ctx.Response.Error("product_label:invalid_label", fmt.Sprintf("id=%d", id))
		return
	}
	
	productLabel.Enable()
	
	ctx.Response.JSON(eel.Map{})
}
