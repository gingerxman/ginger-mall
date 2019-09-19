package product

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	b_product "github.com/gingerxman/ginger-mall/business/product"
)

type Product struct {
	eel.RestResource
}

func (this *Product) Resource() string {
	return "product.product"
}

func (this *Product) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?id:int", "?name", "?corp_id:int", "?raw_product_id:int"},
		"PUT": []string{
			"base_info:json",
			"skus_info:json-array",
			"?media_info:json",
			"?logistics_info:json",
			"?imoney_codes:json-array",
			"?auto_on_shelf:bool",
		},
		"POST": []string{
			"id:int",
			"base_info:json",
			"skus_info:json-array",
			"?media_info:json",
			"?logistics_info:json",
			"?imoney_codes:json-array",
		},
		"DELETE": []string{"id:int"},
	}
}

func (this *Product) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id", 0)
	rawProductId, _ := req.GetInt("raw_product_id", 0)
	name := req.GetString("name")
	corpId, _ := req.GetInt("corp_id", 0)
	
	if rawProductId == 0 && id == 0 && name == "" {
		ctx.Response.Error( "product:invalid_product", "name或id或raw_product_id必须有效")
		return
	}

	bCtx := ctx.GetBusinessContext()
	var corp *account.Corp
	if corpId != 0 {
		corp = account.NewCorpFromOnlyId(bCtx, corpId)
	} else {
		corp = account.GetCorpFromContext(bCtx)
	}
	productPool := b_product.GetProductPoolForCorp(bCtx, corp)
	var poolProduct *b_product.PoolProduct
	if id != 0 {
		poolProduct = productPool.GetPoolProduct(id)
	} else if rawProductId != 0 {
		poolProduct = productPool.GetPoolProductByProductId(rawProductId)
	} else {
		product := b_product.NewProductRepository(bCtx).GetProductByName(name)
		poolProduct = b_product.GetGlobalProductPool(bCtx).GetPoolProductByProductId(product.Id)
	}
	
	if poolProduct == nil {
		ctx.Response.Error( "product:invalid_product", fmt.Sprintf("无效的id(%d)", id))
		return
	}

	fillService := b_product.NewFillPoolProductService(bCtx)
	fillService.Fill([]*b_product.PoolProduct{ poolProduct }, eel.FillOption{
		"with_category": true,
		"with_description": true,
		"with_logistics": true,
		"with_label": true,
		"with_media": true,
		"with_sku": true,
	})
	
	encodeService := b_product.NewEncodePoolProductService(bCtx)
	respData := encodeService.Encode(poolProduct)

	ctx.Response.JSON(respData)
}

func (this *Product) Put(ctx *eel.Context) {
	req := ctx.Request
	baseInfo := req.GetString("base_info")
	strSkuInfos := req.GetString("skus_info")
	strLogisticsInfo := req.GetString("logistics_info")
	mediaInfo := req.GetString("media_info")

	bCtx := ctx.GetBusinessContext()
	//corp := account.GetCorpFromContext(bCtx)
	productFactory := b_product.NewProductFactory(bCtx)
	poolProduct := productFactory.CreateProduct(baseInfo, strSkuInfos, mediaInfo, strLogisticsInfo)
	
	ctx.Response.JSON(eel.Map{
		"id": poolProduct.Id,
		"raw_product_id": poolProduct.ProductId,
	})
}

func (this *Product) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	baseInfo := req.GetString("base_info")
	strSkuInfos := req.GetString("skus_info")
	mediaInfo := req.GetString("media_info")
	strLogisticsInfo := req.GetString("logistics_info")

	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	poolProduct := b_product.GetProductPoolForCorp(bCtx, corp).GetPoolProduct(id)
	updateService := b_product.NewUpdateProductService(bCtx)
	updateService.Update(poolProduct.ProductId, baseInfo, strSkuInfos, mediaInfo, strLogisticsInfo)
	
	ctx.Response.JSON(eel.Map{})
}
