package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type CorpProductProperties struct {
	eel.RestResource
}

func (this *CorpProductProperties) Resource() string {
	return "product.corp_product_properties"
}

func (this *CorpProductProperties) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *CorpProductProperties) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := product.NewProductPropertyRepository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	productProperties, nextPageInfo := repository.GetAllProductPropertiesForCorp(corp, page, filters)

	fillService := product.NewFillProductPropertyService(bCtx)
	fillService.Fill(productProperties, eel.FillOption{
		"with_value": true,
	})

	encodeService := product.NewEncodeProductPropertyService(bCtx)
	rows := encodeService.EncodeMany(productProperties)
	
	ctx.Response.JSON(eel.Map{
		"product_properties": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}
