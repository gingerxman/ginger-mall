package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ProductLabel struct {
	eel.RestResource
}

func (this *ProductLabel) Resource() string {
	return "product.label"
}

func (this *ProductLabel) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"id:int"},
		"PUT": []string{
			"name:string",
		},
		"POST": []string{
			"id:int",
			"name:string",
		},
		"DELETE": []string{"id:int"},
	}
}

func (this *ProductLabel) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductLabelRepository(bCtx)
	productLabel := repository.GetProductLabel(id)

	fillService := product.NewFillProductLabelService(bCtx)
	fillService.FillOne(productLabel, eel.FillOption{
	})

	encodeService := product.NewEncodeProductLabelService(bCtx)
	respData := encodeService.Encode(productLabel)
	
	ctx.Response.JSON(respData)
}

func (this *ProductLabel) Put(ctx *eel.Context) {
	req := ctx.Request
	name := req.GetString("name")

	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productLabel := product.NewProductLabel(
		bCtx, 
		corp,
		name,
	)
	
	ctx.Response.JSON(eel.Map{
		"id": productLabel.Id,
	})
}

func (this *ProductLabel) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	name := req.GetString("name")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductLabelRepository(bCtx)
	productLabel := repository.GetProductLabel(id)

	_ = productLabel.Update(
		name,
	)
	
	ctx.Response.JSON(eel.Map{})
}
