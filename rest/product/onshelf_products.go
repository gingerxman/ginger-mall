package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type OnshelfProducts struct {
	eel.RestResource
}

func (this *OnshelfProducts) Resource() string {
	return "product.onshelf_products"
}

func (this *OnshelfProducts) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?corp_id:int", "?filters:json", "?type"},
		"PUT": []string{"product_ids:json-array"},
	}
}

func (this *OnshelfProducts) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := this.GetBusinessContext()
	page := eel.ExtractPageInfoFromRequest(this.Ctx)
	filters := common.ConvertToBeegoOrmFilter(this.GetFilters())
	
	productType := this.GetString("type")
	if productType == "" {
		filters["product_type"] = "product"
	}
	
	corpId, _ := this.GetInt("corp_id", 0)
	var corp *account.Corporation
	if corpId != 0{
		corp = account.NewCorpFromOnlyId(bCtx, corpId)
	}else{
		corp = account.GetCorpFromContext(bCtx)
	}

	poolProducts, nextPageInfo := product.GetInSaleProductShelfForCorp(bCtx, corp).GetPagedProducts(filters, page)

	//fillService := product.NewFillPoolProductService(bCtx)
	fillService := product.NewFillPoolProductServiceForCorp(bCtx, corp)
	fillService.Fill(poolProducts, eel.FillOption{
		"with_category": true,
		"with_description": true,
		"with_media": true,
		"with_sku": true,
		"with_label": true,
		"with_commission": true,
	})

	encodeService := product.NewEncodePoolProductService(bCtx)
	rows := encodeService.EncodeMany(poolProducts)
	
	response := eel.MakeResponse(eel.Map{
		"products": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
	this.ReturnJSON(response)
}

func (this *OnshelfProducts) Put(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productIds := req.GetIntArray("product_ids")
	
	poolProducts := product.GetProductPoolForCorp(bCtx, corp).GetPoolProductsByIds(productIds)
	product.GetInSaleProductShelfForCorp(bCtx, corp).AddProducts(poolProducts)
	
	product.NewOffshelfApplicationSerivce(bCtx).CancelOffshelfApplicationForPoolProducts(poolProducts)
	
	response := eel.MakeResponse(eel.Map{})
	this.ReturnJSON(response)
}
