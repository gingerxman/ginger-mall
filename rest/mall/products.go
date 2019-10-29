package mall

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type Products struct {
	eel.RestResource
}

func (this *Products) Resource() string {
	return "mall.products"
}

func (this *Products) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?filters:json"},
	}
}

func (this *Products) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	
	corp := account.GetCorpFromContext(bCtx)

	poolProducts, nextPageInfo := product.GetInSaleProductShelfForCorp(bCtx, corp).GetPagedProducts(filters, page)

	fillService := product.NewFillPoolProductServiceForCorp(bCtx, corp)
	fillService.Fill(poolProducts, eel.FillOption{
		"with_category": true,
		"with_description": true,
		"with_logistics": true,
		"with_media": true,
		"with_sku": true,
	})

	encodeService := product.NewEncodePoolProductService(bCtx)
	rows := encodeService.EncodeMany(poolProducts)
	
	ctx.Response.JSON(eel.Map{
		"products": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}
