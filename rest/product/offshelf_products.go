package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type OffshelfProducts struct {
	eel.RestResource
}

func (this *OffshelfProducts) Resource() string {
	return "product.offshelf_products"
}

func (this *OffshelfProducts) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
		"PUT": []string{"product_ids:json-array"},
	}
}

func (this *OffshelfProducts) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	corp := account.GetCorpFromContext(bCtx)
	poolProducts, nextPageInfo := product.GetForSaleProductShelfForCorp(bCtx, corp).GetPagedProducts(filters, page)

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

func (this *OffshelfProducts) Put(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productIds := req.GetIntArray("product_ids")
	
	poolProducts := product.GetProductPoolForCorp(bCtx, corp).GetPoolProductsByIds(productIds)
	product.GetForSaleProductShelfForCorp(bCtx, corp).AddProducts(poolProducts)
	
	ctx.Response.JSON(eel.Map{})
}
