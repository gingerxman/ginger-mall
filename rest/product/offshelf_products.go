package product

import (
	"gpeanut/business/account"
	"gpeanut/business/common"
	"gpeanut/business/product"
	
	"github.com/kfchen81/beego/vanilla"
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

func (this *OffshelfProducts) Get() {
	bCtx := this.GetBusinessContext()
	page := eel.ExtractPageInfoFromRequest(this.Ctx)
	filters := common.ConvertToBeegoOrmFilter(this.GetFilters())
	corp := account.GetCorpFromContext(bCtx)
	poolProducts, nextPageInfo := product.GetForSaleProductShelfForCorp(bCtx, corp).GetPagedProducts(filters, page)

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

func (this *OffshelfProducts) Put() {
	bCtx := this.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productIds := this.GetIntArray("product_ids")
	
	poolProducts := product.GetProductPoolForCorp(bCtx, corp).GetPoolProductsByIds(productIds)
	product.GetForSaleProductShelfForCorp(bCtx, corp).AddProducts(poolProducts)
	
	product.NewOffshelfApplication(bCtx, poolProducts)
	
	response := eel.MakeResponse(eel.Map{})
	this.ReturnJSON(response)
}
