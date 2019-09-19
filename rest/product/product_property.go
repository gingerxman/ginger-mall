package product

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ProductProperty struct {
	eel.RestResource
}

func (this *ProductProperty) Resource() string {
	return "product.property"
}

func (this *ProductProperty) GetParameters() map[string][]string {
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

func (this *ProductProperty) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductPropertyRepository(bCtx)
	productProperty := repository.GetProductProperty(id)
	if productProperty == nil {
		ctx.Response.Error("product_property:invalid_product_property", fmt.Sprintf("id=%d", id))
		return
	}

	fillService := product.NewFillProductPropertyService(bCtx)
	fillService.Fill([]*product.ProductProperty{ productProperty }, eel.FillOption{
		"with_product_property_value": true,
	})

	encodeService := product.NewEncodeProductPropertyService(bCtx)
	respData := encodeService.Encode(productProperty)
	
	ctx.Response.JSON(respData)
}

func (this *ProductProperty) Put(ctx *eel.Context) {
	req := ctx.Request
	name := req.GetString("name")

	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productProperty := product.NewProductProperty(
		bCtx, 
		corp,
		name,
	)
	
	ctx.Response.JSON(eel.Map{
		"id": productProperty.Id,
	})
}

func (this *ProductProperty) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	name := req.GetString("name")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductPropertyRepository(bCtx)
	productProperty := repository.GetProductProperty(id)
	if productProperty == nil {
		ctx.Response.Error("product_property:invalid_product_property", fmt.Sprintf("id=%d", id))
		return
	}

	_ = productProperty.Update(
		name,
	)
	
	ctx.Response.JSON(eel.Map{})
}

func (this *ProductProperty) Delete(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductPropertyRepository(bCtx)
	productProperty := repository.GetProductProperty(id)
	if productProperty == nil {
		ctx.Response.Error("product_property:invalid_product_property", fmt.Sprintf("id=%d", id))
		return
	}

	_ = productProperty.Delete()
	
	ctx.Response.JSON(eel.Map{})
}
