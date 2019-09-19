package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ProductProperties struct {
	eel.RestResource
}

func (this *ProductProperties) Resource() string {
	return "product.product_properties"
}

func (this *ProductProperties) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *ProductProperties) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := product.NewProductPropertyRepository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	productProperties, nextPageInfo := repository.GetEnabledProductPropertiesForCorp(corp, page, filters)

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
