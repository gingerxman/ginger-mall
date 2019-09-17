package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type SubCategories struct {
	eel.RestResource
}

func (this *SubCategories) Resource() string {
	return "product.sub_categories"
}

func (this *SubCategories) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"father_id:int"},
	}
}

func (this *SubCategories) Get(ctx *eel.Context) {
	req := ctx.Request
	fatherId, _ := req.GetInt("father_id")
	
	bCtx := ctx.GetBusinessContext()
	repository := product.NewProductCategoryRepository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	category := repository.GetProductCategoryInCorp(corp, fatherId)
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
