package mall

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
	return "mall.product"
}

func (this *Product) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"id:int"},
	}
}

func (this *Product) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id", 0)
	
	if id == 0 {
		ctx.Response.Error( "mall.product:invalid_product", "id == 0")
		return
	}

	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	productPool := b_product.GetProductPoolForCorp(bCtx, corp)
	poolProduct := productPool.GetPoolProduct(id)
	
	if poolProduct == nil {
		ctx.Response.Error( "mall.product:invalid_product", fmt.Sprintf("无效的id(%d)", id))
		return
	}

	fillService := b_product.NewFillPoolProductService(bCtx)
	fillService.Fill([]*b_product.PoolProduct{ poolProduct }, eel.FillOption{
		"with_category": true,
		"with_description": true,
		"with_logistics": true,
		"with_media": true,
		"with_sku": true,
	})
	
	encodeService := b_product.NewEncodePoolProductService(bCtx)
	respData := encodeService.Encode(poolProduct)

	ctx.Response.JSON(respData)
}
