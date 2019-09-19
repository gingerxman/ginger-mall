package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductRepository struct {
	eel.RepositoryBase
}

func NewProductRepository(ctx context.Context) *ProductRepository {
	repository := new(ProductRepository)
	repository.Ctx = ctx
	return repository
}

func (this *ProductRepository) GetProducts(filters eel.Map, orderExprs ...string) []*Product {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.Product{})
	
	var models []*m_product.Product
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
	
	products := make([]*Product, 0)
	for _, model := range models {
		products = append(products, NewProductFromModel(this.Ctx, model))
	}
	return products
}

func (this *ProductRepository) GetPagedProducts(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*Product, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.Product{})
	
	var models []*m_product.Product
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, err := eel.Paginate(db, page, &models)
	
	if err != nil {
		eel.Logger.Error(db.Error)
		return nil, paginateResult
	}
	
	products := make([]*Product, 0)
	for _, model := range models {
		products = append(products, NewProductFromModel(this.Ctx, model))
	}
	return products, paginateResult
}

//GetEnabledProductsForCorp 获得启用的Product对象集合
func (this *ProductRepository) GetEnabledProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*Product, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedProducts(filters, page, "id desc")
}

//GetProductsWithLabelForCorp 获得label对应的Product对象集合
func (this *ProductRepository) GetProductIdsWithLabelForCorp(corp business.ICorp, labelId int, page *eel.PageInfo, filters eel.Map) ([]int, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	models := make([]*m_product.ProductHasLabel, 0)
	db := o.Model(&m_product.ProductHasLabel{}).Where("label_id", labelId).Order("id desc")
	nextPaginateResult, db := eel.Paginate(db, page, &models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return make([]int, 0), eel.MockPaginate(0, page)
	}
	
	if len(models) == 0 {
		return make([]int, 0), eel.MockPaginate(0, page)
	}
	
	productIds := make([]int, 0)
	for _, model := range models {
		productIds = append(productIds, model.ProductId)
	}
	return productIds, nextPaginateResult
}

//GetAllProductsForCorp 获得所有Product对象集合
func (this *ProductRepository) GetAllProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*Product, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedProducts(filters, page, "id desc")
}

//GetProductInCorp 根据id和corp获得Product对象
func (this *ProductRepository) GetProductInCorp(corp business.ICorp, id int) *Product {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	products := this.GetProducts(filters)
	
	if len(products) == 0 {
		return nil
	} else {
		return products[0]
	}
}

//GetProduct 根据id和corp获得Product对象
func (this *ProductRepository) GetProduct(id int) *Product {
	filters := eel.Map{
		"id": id,
	}
	
	products := this.GetProducts(filters)
	
	if len(products) == 0 {
		return nil
	} else {
		return products[0]
	}
}

func (this *ProductRepository) GetProductByName(name string) *Product {
	filters := eel.Map{
		"name": name,
	}
	
	products := this.GetProducts(filters)
	
	if len(products) == 0 {
		return nil
	} else {
		return products[0]
	}
}


//GetProduct 根据id和corp获得Product对象
func (this *ProductRepository) GetProductsByIds(ids []int) []*Product {
	filters := eel.Map{
		"id__in": ids,
	}
	
	products := this.GetProducts(filters)
	return products
}

//GetProduct 根据id和corp获得Product对象
func (this *ProductRepository) GetProductsUseLimitZone(limitZoneId int) []*Product {
	o := eel.GetOrmFromContext(this.Ctx)
	var models []*m_product.ProductLogisticsInfo
	
	db := o.Model(&m_product.ProductLogisticsInfo{}).Where("limit_zone_id", limitZoneId).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return make([]*Product, 0)
	}
	
	if len(models) == 0 {
		return make([]*Product, 0)
	}
	
	productIds := make([]int, 0)
	for _, model := range models {
		productIds = append(productIds, model.ProductId)
	}
	filters := eel.Map{
		"id__in": productIds,
		"is_deleted": false,
	}
	
	products := this.GetProducts(filters)
	return products
}

// IsValidProductSku 判断(productId, skuName)的组合是否是可用的商品sku
func (this *ProductRepository) IsValidProductSku(productId int, skuName string) bool {
	o := eel.GetOrmFromContext(this.Ctx)
	return o.Model(&m_product.ProductSku{}).Where(eel.Map{
		"product_id": productId,
		"name": skuName,
	}).Exist()
}

func init() {
}
