package product

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeProductPropertyValueService struct {
	eel.ServiceBase
}

func NewEncodeProductPropertyValueService(ctx context.Context) *EncodeProductPropertyValueService {
	service := new(EncodeProductPropertyValueService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeProductPropertyValueService) Encode(productPropertyValue *ProductPropertyValue) *RProductPropertyValue {
	if productPropertyValue == nil {
		return nil
	}

	return &RProductPropertyValue{
		Id: productPropertyValue.Id,
		Text: productPropertyValue.Text,
		Image: productPropertyValue.Image,
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeProductPropertyValueService) EncodeMany(product_property_values []*ProductPropertyValue) []*RProductPropertyValue {
	rDatas := make([]*RProductPropertyValue, 0)
	for _, productPropertyValue := range product_property_values {
		rDatas = append(rDatas, this.Encode(productPropertyValue))
	}
	
	return rDatas
}

func init() {
}
