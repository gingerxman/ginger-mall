package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductCategoryRepository struct {
	eel.RepositoryBase
}

func NewProductCategoryRepository(ctx context.Context) *ProductCategoryRepository {
	repository := new(ProductCategoryRepository)
	repository.Ctx = ctx
	return repository
}

func (this *ProductCategoryRepository) GetProductCategories(filters eel.Map, orderExprs ...string) []*ProductCategory {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductCategory{})
	
	var models []*m_product.ProductCategory
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
	
	ProductCategories := make([]*ProductCategory, 0)
	for _, model := range models {
		ProductCategories = append(ProductCategories, NewProductCategoryFromModel(this.Ctx, model))
	}
	return ProductCategories
}

func (this *ProductCategoryRepository) GetPagedProductCategories(filters eel.Map, page *eel.PageInfo, orderExprs string) ([]*ProductCategory, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.ProductCategory{})
	
	var models []*m_product.ProductCategory
	if len(filters) > 0 {
		db = db.Where(filters)
	}	
	if len(orderExprs) > 0 {
		db = db.Order(orderExprs)
	}
	paginateResult, db := eel.Paginate(db, page, &models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil, paginateResult
	}
	
	ProductCategories := make([]*ProductCategory, 0)
	for _, model := range models {
		ProductCategories = append(ProductCategories, NewProductCategoryFromModel(this.Ctx, model))
	}
	return ProductCategories, paginateResult
}

//GetEnabledProductCategoriesForCorp 获得启用的ProductCategory对象集合
func (this *ProductCategoryRepository) GetEnabledProductCategoriesForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*ProductCategory, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_enabled"] = true
	
	return this.GetPagedProductCategories(filters, page, "id desc")
}

//GetCorpProductCategories 获得corp的ProductCategory对象集合
func (this *ProductCategoryRepository) GetCorpProductCategories(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*ProductCategory, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedProductCategories(filters, page, "id desc")
}

//GetAllProductCategoriesForCorp 获得corp的ProductCategory对象集合
func (this *ProductCategoryRepository) GetAllProductCategories() []*ProductCategory {
	filters := eel.Map{
		"is_enabled": true,
	}
	
	return this.GetProductCategories(filters, "father_id")
}

//GetProductCategoryInCorp 根据id和corp获得ProductCategory对象
func (this *ProductCategoryRepository) GetProductCategoryInCorp(corp business.ICorp, id int) *ProductCategory {
	if id == 0 {
		return NewRootProductCategory(this.Ctx, corp)
	}
	
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	productCategories := this.GetProductCategories(filters)
	
	if len(productCategories) == 0 {
		return nil
	} else {
		productCategory := productCategories[0]
		productCategory.Corp = corp
		return productCategory
	}
}

//GetProductCategory 根据id和corp获得ProductCategory对象
func (this *ProductCategoryRepository) GetProductCategory(id int) *ProductCategory {
	if id == 0 {
		return NewRootProductCategory(this.Ctx, nil)
	}
	
	filters := eel.Map{
		"id": id,
	}
	
	ProductCategories := this.GetProductCategories(filters)
	
	if len(ProductCategories) == 0 {
		return nil
	} else {
		return ProductCategories[0]
	}
}

//GetProductCategoryPathEndWithId 获得以id结尾的category tree path
func (this *ProductCategoryRepository) GetProductCategoryPathEndWithId(endId int) []*ProductCategory {
	categoryId := endId
	categories := make([]*ProductCategory, 0)
	
	for {
		if categoryId == 0 {
			break
		}
		
		category := this.GetProductCategory(categoryId)
		if category == nil {
			break
		}
		
		categories = append(categories, category)
		categoryId = category.FatherId
	}
	
	//reverse it
	reversDatas := make([]*ProductCategory, 0)
	for i := len(categories)-1; i >= 0; i-- {
		reversDatas = append(reversDatas, categories[i])
	}
	return reversDatas
}

func init() {
}
