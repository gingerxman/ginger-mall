package product

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
)

type ProductPropertyValue struct {
	eel.RestResource
}

func (this *ProductPropertyValue) Resource() string {
	return "product.property_value"
}

func (this *ProductPropertyValue) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{
			"property_id:int",
			"text:string",
			"image:string",
		},
		"DELETE": []string{
			"property_id:int",
			"id:int",
		},
	}
}

func (this *ProductPropertyValue) Put(ctx *eel.Context) {
	req := ctx.Request
	propertyId, _ := req.GetInt("property_id")
	text := req.GetString("text")
	image := req.GetString("image")

	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	property := product.NewProductPropertyRepository(bCtx).GetProductPropertyInCorp(corp, propertyId)
	propertyValue := property.AddNewValue(text, image)
	
	ctx.Response.JSON(eel.Map{
		"id": propertyValue.Id,
	})
}

func (this *ProductPropertyValue) Delete(ctx *eel.Context) {
	req := ctx.Request
	propertyId, _ := req.GetInt("property_id")
	valueId, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	property := product.NewProductPropertyRepository(bCtx).GetProductPropertyInCorp(corp, propertyId)
	property.DeleteValue(valueId)
	
	ctx.Response.JSON(eel.Map{})
}
