package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type CorpProductLabels struct {
	eel.RestResource
}

func (this *CorpProductLabels) Resource() string {
	return "product.corp_labels"
}

func (this *CorpProductLabels) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *CorpProductLabels) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := product.NewProductLabelRepository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	productLabels, nextPageInfo := repository.GetAllProductLabelsForCorp(corp, page, filters)

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
