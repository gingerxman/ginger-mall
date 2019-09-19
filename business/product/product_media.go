package product

import (
	"context"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductMedia struct {
	eel.EntityBase
	Id int
	ProductId int
	Type string
	Url string
}

//根据model构建对象
func NewProductMediaFromModel(ctx context.Context, model *m_product.ProductMedia) *ProductMedia {
	instance := new(ProductMedia)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.ProductId = model.ProductId
	instance.Type = m_product.PRODUCTMEDIA2STR[model.Type]
	instance.Url = model.Url

	return instance
}

func init() {
}
