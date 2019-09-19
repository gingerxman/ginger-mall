package product

import (
	"context"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductDescription struct {
	eel.EntityBase
	Id int
	ProductId int
	Introduction string
	Detail string
	Remark string
}

//根据model构建对象
func NewProductDescriptionFromModel(ctx context.Context, model *m_product.ProductDescription) *ProductDescription {
	instance := new(ProductDescription)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.ProductId = model.ProductId
	instance.Introduction = model.Introduction
	instance.Detail = model.Detail
	instance.Remark = model.Remark

	return instance
}

func init() {
}
