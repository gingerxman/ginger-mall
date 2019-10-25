package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type CorpProducts struct {
	eel.RestResource
}

func (this *CorpProducts) Resource() string {
	return "product.corp_products"
}

func (this *CorpProducts) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"corp_id:int", "?filters:json", "?type"},
	}
}

func (this *CorpProducts) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	
	productType := req.GetString("type")
	if productType == "" {
		filters["product_type"] = "product"
	}
	
	corpId, _ := req.GetInt("corp_id", 0)
	var corp *account.Corp
	if corpId != 0{
		corp = account.NewCorpFromOnlyId(bCtx, corpId)
	}else{
		corp = account.GetCorpFromContext(bCtx)
	}

	productPool := product.GetProductPoolForCorp(bCtx, corp)
	poolProducts, nextPageInfo := productPool.SearchProducts(filters, page)// product.GetInSaleProductShelfForCorp(bCtx, corp).GetPagedProducts(filters, page)
	eel.Logger.Debug(len(poolProducts))

	fillService := product.NewFillPoolProductServiceForCorp(bCtx, corp)
	fillService.Fill(poolProducts, eel.FillOption{
		"with_category": true,
		"with_description": true,
		"with_media": true,
		"with_sku": true,
		"with_label": true,
	})

	encodeService := product.NewEncodePoolProductService(bCtx)
	rows := encodeService.EncodeMany(poolProducts)
	
	ctx.Response.JSON(eel.Map{
		"products": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}

