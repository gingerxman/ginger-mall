package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ProductLabels struct {
	eel.RestResource
}

func (this *ProductLabels) Resource() string {
	return "product.labels"
}

func (this *ProductLabels) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *ProductLabels) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := product.NewProductLabelRepository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	productLabels, nextPageInfo := repository.GetEnabledProductLabelsForCorp(corp, page, filters)

	fillService := product.NewFillProductLabelService(bCtx)
	fillService.Fill(productLabels, eel.FillOption{
	})

	encodeService := product.NewEncodeProductLabelService(bCtx)
	rows := encodeService.EncodeMany(productLabels)
	
	ctx.Response.JSON(eel.Map{
		"labels": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}
