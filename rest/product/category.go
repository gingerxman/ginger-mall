package product

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type Category struct {
	eel.RestResource
}

func (this *Category) Resource() string {
	return "product.category"
}

func (this *Category) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"id:int"},
		"PUT": []string{
			"name:string",
			"father_id:int",
		},
		"POST": []string{
			"id:int",
			"name:string",
		},
		"DELETE": []string{"id:int"},
	}
}

func (this *Category) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductCategoryRepository(bCtx)
	productCategory := repository.GetProductCategory(id)
	if productCategory == nil {
		ctx.Response.Error("product_category:invalid_category", fmt.Sprintf("id=%d", id))
		return
	}

	fillService := product.NewFillProductCategoryService(bCtx)
	fillService.Fill([]*product.ProductCategory{ productCategory }, eel.FillOption{
	})

	encodeService := product.NewEncodeProductCategoryService(bCtx)
	respData := encodeService.Encode(productCategory)

	ctx.Response.JSON(respData)
}

func (this *Category) Put(ctx *eel.Context) {
	req := ctx.Request
	name := req.GetString("name")
	fatherId, _ := req.GetInt("father_id", 0)
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productCategory := product.NewProductCategory(
		bCtx, 
		corp,
		name,
		fatherId,
	)

	ctx.Response.JSON(eel.Map{
		"id": productCategory.Id,
	})
}

func (this *Category) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	name := req.GetString("name")

	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductCategoryRepository(bCtx)
	productCategory := repository.GetProductCategory(id)
	if productCategory == nil {
		ctx.Response.Error("product_category:invalid_category", fmt.Sprintf("id=%d", id))
		return
	}

	_ = productCategory.Update(
		name,
	)
	
	ctx.Response.JSON(eel.Map{})
}
