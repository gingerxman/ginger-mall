package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductPropertyRepository struct {
	eel.RepositoryBase
}

func NewProductPropertyRepository(ctx context.Context) *ProductPropertyRepository {
	repository := new(ProductPropertyRepository)
	repository.Ctx = ctx
	return repository
}

func (this *ProductPropertyRepository) GetProductProperties(filters eel.Map, orderExprs ...string) []*ProductProperty {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductProperty{})
	
	var models []*m_product.ProductProperty
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	db = db.Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	productProperties := make([]*ProductProperty, 0)
	for _, model := range models {
		productProperties = append(productProperties, NewProductPropertyFromModel(this.Ctx, model))
	}
	return productProperties
}

func (this *ProductPropertyRepository) GetPagedProductProperties(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*ProductProperty, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductProperty{})
	
	var models []*m_product.ProductProperty
	filters["is_deleted"] = false
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, db := eel.Paginate(db, page, &models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil, paginateResult
	}
	
	productProperties := make([]*ProductProperty, 0)
	for _, model := range models {
		productProperties = append(productProperties, NewProductPropertyFromModel(this.Ctx, model))
	}
	return productProperties, paginateResult
}

//GetEnabledProductPropertiesForCorp 获得启用的ProductProperty对象集合
func (this *ProductPropertyRepository) GetEnabledProductPropertiesForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*ProductProperty, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_deleted"] = false
	
	return this.GetPagedProductProperties(filters, page)
}

//GetAllProductPropertiesForCorp 获得所有ProductProperty对象集合
func (this *ProductPropertyRepository) GetAllProductPropertiesForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*ProductProperty, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedProductProperties(filters, page)
}

//GetProductPropertyInCorp 根据id和corp获得ProductProperty对象
func (this *ProductPropertyRepository) GetProductPropertyInCorp(corp business.ICorp, id int) *ProductProperty {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	productProperties := this.GetProductProperties(filters)
	
	if len(productProperties) == 0 {
		return nil
	} else {
		return productProperties[0]
	}
}

//GetProductPropertiesInCorps 根据id和corp获得ProductProperty对象
func (this *ProductPropertyRepository) GetProductPropertiesInCorps(corpIds []int) []*ProductProperty {
	filters := eel.Map{
		"corp_id__in": corpIds,
		"is_deleted": false,
	}
	
	productProperties := this.GetProductProperties(filters)
	return productProperties
}

//GetProductProperty 根据id和corp获得ProductProperty对象
func (this *ProductPropertyRepository) GetProductProperty(id int) *ProductProperty {
	filters := eel.Map{
		"id": id,
	}
	
	productProperties := this.GetProductProperties(filters)
	
	if len(productProperties) == 0 {
		return nil
	} else {
		return productProperties[0]
	}
}

//GetProductPropertyValue 根据id获得ProductPropertyValue对象
func (this *ProductPropertyRepository) GetProductPropertyValue(id int) *ProductPropertyValue {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductPropertyValue{})
	
	var model m_product.ProductPropertyValue
	db = db.Where("id", id).Take(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	return NewProductPropertyValueFromModel(this.Ctx, &model)
}

func init() {
}
