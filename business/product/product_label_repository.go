package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductLabelRepository struct {
	eel.RepositoryBase
}

func NewProductLabelRepository(ctx context.Context) *ProductLabelRepository {
	repository := new(ProductLabelRepository)
	repository.Ctx = ctx
	return repository
}

func (this *ProductLabelRepository) GetProductLabels(filters eel.Map, orderExprs ...string) []*ProductLabel {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductLabel{})
	
	var models []*m_product.ProductLabel
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	if len(orderExprs) > 0 {
		db = db.Order(orderExprs)
	}
	db = db.Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	labels := make([]*ProductLabel, 0)
	for _, model := range models {
		labels = append(labels, NewProductLabelFromModel(this.Ctx, model))
	}
	return labels
}

func (this *ProductLabelRepository) GetPagedProductLabels(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*ProductLabel, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductLabel{})
	
	var models []*m_product.ProductLabel
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
	
	labels := make([]*ProductLabel, 0)
	for _, model := range models {
		labels = append(labels, NewProductLabelFromModel(this.Ctx, model))
	}
	return labels, paginateResult
}

//GetEnabledProductLabelsForCorp 获得启用的ProductLabel对象集合
func (this *ProductLabelRepository) GetEnabledProductLabelsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*ProductLabel, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_enabled"] = true
	
	return this.GetPagedProductLabels(filters, page, "id desc")
}

//GetAllProductLabelsForCorp 获得所有ProductLabel对象集合
func (this *ProductLabelRepository) GetAllProductLabelsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*ProductLabel, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedProductLabels(filters, page, "id desc")
}

//GetProductLabelInCorp 根据id和corp获得ProductLabel对象
func (this *ProductLabelRepository) GetProductLabelInCorp(corp business.ICorp, id int) *ProductLabel {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	labels := this.GetProductLabels(filters)
	
	if len(labels) == 0 {
		return nil
	} else {
		return labels[0]
	}
}

func (this *ProductLabelRepository) GetProductLabelsInCorp(corp business.ICorp, ids []int) []*ProductLabel {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id__in": ids,
	}
	
	return this.GetProductLabels(filters)
}

//GetProductLabel 根据id和corp获得ProductLabel对象
func (this *ProductLabelRepository) GetProductLabel(id int) *ProductLabel {
	filters := eel.Map{
		"id": id,
	}
	
	labels := this.GetProductLabels(filters)
	
	if len(labels) == 0 {
		return nil
	} else {
		return labels[0]
	}
}

func init() {
}
