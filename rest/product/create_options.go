package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type CreateOptions struct {
	eel.RestResource
}

func (this *CreateOptions) Resource() string {
	return "product.create_options"
}

func (this *CreateOptions) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *CreateOptions) Get(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductCategoryRepository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	category := repository.GetProductCategoryInCorp(corp, 0)
	categories := category.GetSubCategories()

	fillService := product.NewFillProductCategoryService(bCtx)
	fillService.Fill(categories, eel.FillOption{
	})

	encodeService := product.NewEncodeProductCategoryService(bCtx)
	rows := encodeService.EncodeMany(categories)
	
	ctx.Response.JSON(eel.Map{
		"categories": rows,
	})
}
