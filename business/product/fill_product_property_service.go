package product

import (
	"context"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type FillProductPropertyService struct {
	eel.ServiceBase
}

func NewFillProductPropertyService(ctx context.Context) *FillProductPropertyService {
	service := new(FillProductPropertyService)
	service.Ctx = ctx
	return service
}

func (this *FillProductPropertyService) Fill(productProperties []*ProductProperty, option eel.FillOption) {
	if len(productProperties) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, productProperty := range productProperties {
		ids = append(ids, productProperty.Id)
	}

	if enableOption, ok := option["with_value"]; ok && enableOption {
		this.fillProductPropertyValue(productProperties, ids)
	}
	return
}


func (this *FillProductPropertyService) fillProductPropertyValue(productProperties []*ProductProperty, ids []int) {
	//从db中获取数据集合
	var models []*m_product.ProductPropertyValue
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductPropertyValue{}).Where("property_id__in", ids).Where("is_deleted", false).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	//构建<id, productProperty>
	id2property := make(map[int]*ProductProperty)
	for _, productProperty := range productProperties {
		id2property[productProperty.Id] = productProperty
	}
	
	//填充product_property的ProductPropertyValue对象
	for _, model := range models {
		if productProperty, ok := id2property[model.PropertyId]; ok {
			productProperty.AppendValue(NewProductPropertyValueFromModel(this.Ctx, model))
		}
	}
}


func init() {
}
