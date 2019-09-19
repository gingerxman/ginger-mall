package product

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeProductPropertyService struct {
	eel.ServiceBase
}

func NewEncodeProductPropertyService(ctx context.Context) *EncodeProductPropertyService {
	service := new(EncodeProductPropertyService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeProductPropertyService) Encode(productProperty *ProductProperty) *RProductProperty {
	if productProperty == nil {
		return nil
	}
	rValues := NewEncodeProductPropertyValueService(this.Ctx).EncodeMany(productProperty.Values)

	return &RProductProperty{
		Id: productProperty.Id,
		Name: productProperty.Name,
		IsDeleted: productProperty.IsDeleted,
		Values: rValues,
		CreatedAt: productProperty.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeProductPropertyService) EncodeMany(product_properties []*ProductProperty) []*RProductProperty {
	rDatas := make([]*RProductProperty, 0)
	for _, productProperty := range product_properties {
		rDatas = append(rDatas, this.Encode(productProperty))
	}
	
	return rDatas
}

func init() {
}
