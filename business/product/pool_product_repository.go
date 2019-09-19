package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"strings"
)

type PoolProductRepository struct {
	eel.RepositoryBase
}

func NewPoolProductRepository(ctx context.Context) *PoolProductRepository {
	repository := new(PoolProductRepository)
	repository.Ctx = ctx
	return repository
}

// fillFiltersWithFilteredProducts 从filters中找出商品属性的查询参数，填充合适的对应参数到filters中
func (this *PoolProductRepository) fillFiltersWithFilteredProducts(filters eel.Map){
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.Product{})
	productFilters := make(eel.Map)
	if filters == nil{
		filters = eel.Map{}
	}else{
		if productName, ok := filters["name"]; ok{
			productFilters["name"] = productName
			delete(filters, "name")
		}
	}

	if len(productFilters) > 0{
		productIds := make([]int, 0)
		var dbModels []*m_product.Product
		db := db.Where(productFilters).Find(&dbModels)
		if db.Error != nil{
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("product:fetch_failed", "查询商品信息失败"))
		}
		for _, dbModel := range dbModels{
			productIds = append(productIds, dbModel.Id)
		}

		if len(productIds) > 0{
			filters["product_id__in"] = productIds
		}
	}
}

func (this *PoolProductRepository) GetPoolProducts(filters eel.Map, orderExprs ...string) []*PoolProduct {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.PoolProduct{})
	
	var models []*m_product.PoolProduct
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
	
	PoolProducts := make([]*PoolProduct, 0)
	for _, model := range models {
		PoolProducts = append(PoolProducts, NewPoolProductFromModel(this.Ctx, model))
	}
	return PoolProducts
}

func checkFieldNames(fields []string, fieldName string) bool {
	ok := false
	for _, field := range fields {
		if fieldName == field {
			ok = true
		}
	}
	return ok
}

func (this *PoolProductRepository) parseFilters(filters map[string]interface{}) map[string]interface{}{
	productFilters := eel.Map{}
	productFields := []string{"name"}

	poolFilters := eel.Map{}

	type2filters := eel.Map{}

	for key, value := range filters {
		keyString := strings.Split(key, "__")
		fieldName := keyString[0]
		if ok := checkFieldNames(productFields, fieldName); ok{
			productFilters[key] = value
			type2filters["productFilters"] = productFilters
		} else {
			poolFilters[key] = value
			type2filters["poolFilters"] = poolFilters
		}
	}
	return type2filters
}

func (this *PoolProductRepository) GetPagedPoolProducts(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*PoolProduct, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.PoolProduct{})

	//this.fillFiltersWithFilteredProducts(filters)

	type2filters := this.parseFilters(filters)
	if poolFilters, ok := type2filters["poolFilters"]; ok && poolFilters != nil {
		db = db.Where(poolFilters)
	}

	if productFilters, ok := type2filters["productFilters"]; ok && productFilters != nil {
		products := NewProductRepository(this.Ctx).GetProducts(productFilters.(eel.Map))
		productIds := []int{0}
		for _, product := range products {
			productIds = append(productIds, product.Id)
		}
		db = db.Where("product_id__in", productIds)
	}

	var models []*m_product.PoolProduct
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, db := eel.Paginate(db, page, &models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil, paginateResult
	}
	
	PoolProducts := make([]*PoolProduct, 0)
	for _, model := range models {
		PoolProducts = append(PoolProducts, NewPoolProductFromModel(this.Ctx, model))
	}
	return PoolProducts, paginateResult
}

//GetEnabledPoolProductsForCorp 获得启用的PoolProduct对象集合
func (this *PoolProductRepository) GetEnabledPoolProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*PoolProduct, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	//filters["is_enabled"] = true
	
	return this.GetPagedPoolProducts(filters, page, "id desc")
}

//GetAllPoolProductsForCorp 获得所有PoolProduct对象集合
func (this *PoolProductRepository) GetAllPoolProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*PoolProduct, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedPoolProducts(filters, page, "id desc")
}

//GetPoolProductInCorp 根据id和corp获得PoolProduct对象
func (this *PoolProductRepository) GetPoolProductInCorp(corp business.ICorp, id int) *PoolProduct {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	PoolProducts := this.GetPoolProducts(filters)
	
	if len(PoolProducts) == 0 {
		return nil
	} else {
		return PoolProducts[0]
	}
}

//GetPoolProduct 根据id和corp获得PoolProduct对象
func (this *PoolProductRepository) GetPoolProduct(id int) *PoolProduct {
	filters := eel.Map{
		"id": id,
	}
	
	PoolProducts := this.GetPoolProducts(filters)
	
	if len(PoolProducts) == 0 {
		return nil
	} else {
		return PoolProducts[0]
	}
}

//GetPoolProductsByProductIds 根据id和corp获得PoolProduct对象
func (this *PoolProductRepository) GetPoolProductsByProductIdsForCorp(corp business.ICorp, productIds []int) []*PoolProduct {
	filters := eel.Map{
		"product_id__in": productIds,
		"corp_id": corp.GetId(),
	}
	return this.GetPoolProducts(filters)
}

// getSelectablePlatformPoolProducts 获得可被选择的平台商品集合
func (this *PoolProductRepository) getSelectablePlatformPoolProducts(platformCorpId int, page *eel.PageInfo, filters eel.Map) ([]*PoolProduct, eel.INextPageInfo) {
	filters["corp_id"] = platformCorpId
	filters["product_type"] = "product"
	filters["status"] = m_product.PP_STATUS_ON
	//filters["status__in"] = []int{m_product.PP_STATUS_ON, m_product.PP_STATUS_ON_POOL, m_product.PP_STATUS_OFF}
	poolProducts, paginateResult := this.GetPagedPoolProducts(filters, page, "id desc")
	
	return poolProducts, paginateResult
}

// GetSelectablePlatformPoolProductsForCorp 获得可供corp选择的平台商品集合
func (this *PoolProductRepository) GetSelectablePlatformPoolProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*PoolProduct, eel.INextPageInfo) {
	platformCorpId := corp.GetPlatformId()
	poolProducts, paginateResult := this.getSelectablePlatformPoolProducts(platformCorpId, page, filters)

	// todo 之后优化
	// 若不是平台则重新获取商品池的状态
	if !corp.IsPlatform(){
		productIds := make([]int, 0)
		for _, poolProduct := range poolProducts{
			productIds = append(productIds, poolProduct.ProductId)
		}
		// 获取商品池根据productIds
		corpProductId2Status := make(map[int]string, 0)
		if len(productIds) > 0 {
			corpPoolProducts := this.GetPoolProductsByProductIdsForCorp(corp, productIds)
			if corpPoolProducts != nil{
				for _, corpPoolProduct := range corpPoolProducts{
					corpProductId2Status[corpPoolProduct.ProductId] = corpPoolProduct.Status
				}
			}
		}
		for _, poolProduct := range poolProducts{
			if status, ok := corpProductId2Status[poolProduct.ProductId]; ok{
				poolProduct.Status = status
			} else{
				poolProduct.Status = m_product.PPSTATUS2STR[m_product.PP_STATUS_ON_POOL]
			}
		}
	}

	return poolProducts, paginateResult
}


func init() {
}
