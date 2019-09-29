package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)

const NEW_PRODUCT_DISPLAY_INDEX = 9999999

type ProductPool struct {
	eel.RepositoryBase
	Corp business.ICorp
}


func (this *ProductPool) AddProducts(products []*Product, supplierCorpId int) []*PoolProduct {
	if !this.Corp.IsValid() {
		eel.Logger.Error("向corp=nil的product pool添加商品")
		panic(eel.NewBusinessError("product_pool:add_product_fail", "向无效的corp添加商品"))
	}
	
	productIds := make([]int, 0)
	id2product := make(map[int]*Product)
	for _, product := range products {
		productIds = append(productIds, product.Id)
		id2product[product.Id] = product
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	//活得数据库中已经存在的数据
	var models []*m_product.PoolProduct
	db := o.Model(&m_product.PoolProduct{}).Where(eel.Map{
		"corp_id": this.Corp.GetId(),
		"product_id__in": productIds,
	}).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("pool_product:create_fail_1", "创建pool product失败"))
	}
	
	//获得需要更新的数据和需要创建的数据
	tobeUpdateProductIds := make([]int, 0)
	tobeCreateProductIds := make([]int, 0)
	id2exist := make(map[int]bool)
	for _, model := range models {
		id2exist[model.Id] = true
	}
	for _, product := range products {
		if _, ok := id2exist[product.Id]; ok {
			tobeUpdateProductIds = append(tobeUpdateProductIds, product.Id)
		} else {
			tobeCreateProductIds = append(tobeCreateProductIds, product.Id)
		}
	}
	
	//首先恢复之前停售的商品状态
	if len(tobeUpdateProductIds) > 0 {
		db = o.Model(&m_product.PoolProduct{}).Where(eel.Map{
			"product_id__in": tobeUpdateProductIds,
			"corp_id":        this.Corp.GetId(),
		}).Update(gorm.Params{
			"status": m_product.PP_STATUS_ON_POOL,
		})
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("pool_product:create_fail_2", "创建pool product失败"))
		}
	}
	
	//创建新的PoolProduct记录
	newModels := make([]*m_product.PoolProduct, 0)
	corpId := this.Corp.GetId()
	nowTime := time.Now()
	for _, productId := range tobeCreateProductIds {
		if _, ok := id2product[productId]; ok {
			model := m_product.PoolProduct{}
			model.CorpId = corpId
			model.ProductId = productId
			model.ProductType = id2product[productId].Type
			model.SupplierId = supplierCorpId
			model.Status = m_product.PP_STATUS_OFF
			model.Type = m_product.PP_TYPE_CREATE
			model.DisplayIndex = NEW_PRODUCT_DISPLAY_INDEX
			model.SyncAt = nowTime
			
			newModels = append(newModels, &model)
		}
	}
	if len(newModels) > 0 {
		//TODO: 替换为o.BatchInsert方案
		for _, newModel := range newModels {
			db = o.Create(&newModel)
			if db.Error != nil {
				eel.Logger.Error(db.Error)
				panic(eel.NewBusinessError("pool_product:create_fail_3", "创建pool product失败"))
			}
		}
		
		//_, err = o.InsertMulti(len(newModels), newModels)
		//if err != nil {
		//	eel.Logger.Error(err)
		//	panic(eel.NewBusinessError("pool_product:create_fail_3", "创建pool product失败"))
		//}
	}
	
	poolProducts := make([]*PoolProduct, 0)
	for _, model := range newModels {
		poolProducts = append(poolProducts, NewPoolProductFromModel(this.Ctx, model))
	}
	
	return poolProducts
}

func (this *ProductPool) AddProduct(product *Product, supplierCorpId int) *PoolProduct {
	poolProducts := this.AddProducts([]*Product{product}, supplierCorpId)
	if len(poolProducts) > 0 {
		return poolProducts[0]
	} else {
		return nil
	}
}

func (this *ProductPool) GetPoolProducts(filters eel.Map, orderExprs ...string) []*PoolProduct {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_product.PoolProduct{})
	
	var models []*m_product.PoolProduct
	if this.Corp != nil && this.Corp.IsValid() {
		filters["corp_id"] = this.Corp.GetId()
	}
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
	
	poolProducts := make([]*PoolProduct, 0)
	for _, model := range models {
		poolProducts = append(poolProducts, NewPoolProductFromModel(this.Ctx, model))
	}
	return poolProducts
}

func (this *ProductPool) GetPagedPoolProducts(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*PoolProduct, eel.INextPageInfo) {
	filters["corp_id"] = this.Corp.GetId()
	return NewPoolProductRepository(this.Ctx).GetPagedPoolProducts(filters, page, orderExprs...)
}

func (this *ProductPool) GetPoolProduct(id int) *PoolProduct {
	filters := eel.Map{
		"id": id,
	}
	
	poolProducts := this.GetPoolProducts(filters)
	
	if len(poolProducts) == 0 {
		return nil
	} else {
		return poolProducts[0]
	}
}

func (this *ProductPool) GetPoolProductByProductId(productId int) *PoolProduct {
	filters := eel.Map{
		"product_id": productId,
	}
	
	poolProducts := this.GetPoolProducts(filters)
	
	if len(poolProducts) == 0 {
		return nil
	} else {
		return poolProducts[0]
	}
}

func (this *ProductPool) GetPoolProductsByIds(ids []int) []*PoolProduct {
	filters := eel.Map{
		"id__in": ids,
	}
	
	return this.GetPoolProducts(filters)
}

func (this *ProductPool) GetPoolProductsByProductIds(ids []int) []*PoolProduct {
	filters := eel.Map{
		"product_id__in": ids,
	}
	
	return this.GetPoolProducts(filters)
}

func (this *ProductPool) SearchProducts(filters eel.Map, page *eel.PageInfo) ([]*PoolProduct, eel.INextPageInfo) {
	//TODO: 使用es作为搜索解决方案
	productFilters := eel.Map{}
	poolProductFilters := eel.Map{}
	
	for key, value := range filters {
		switch key {
		case "name__contains":
			productFilters[key] = value
		default:
			poolProductFilters[key] = value
		}
	}
	
	var nextPageInfo eel.INextPageInfo
	if len(productFilters) > 0{
		var products []*Product
		products, nextPageInfo = NewProductRepository(this.Ctx).GetPagedProducts(productFilters, page)
		productIds := make([]int, 0)
		for _, product := range products {
			productIds = append(productIds, product.Id)
		}
		
		if len(productIds) > 0 {
			poolProductFilters["product_id__in"] = productIds
		}
	} else {
		nextPageInfo = eel.MockPaginate(0, page)
	}
	
	if len(poolProductFilters) > 0 {
		poolProductFilters["corp_id"] = this.Corp.GetId()
		poolProducts := this.GetPoolProducts(poolProductFilters, "-id")
		return poolProducts, nextPageInfo
	} else {
		return make([]*PoolProduct, 0), nextPageInfo
	}
}

//SearchCategoryProducts 获得product pool中属于指定category的商品集合
func (this *ProductPool) SearchCategoryProducts(productCategory *ProductCategory, filters eel.Map, page *eel.PageInfo) ([]*PoolProduct, eel.INextPageInfo) {
	products := NewProductRepository(this.Ctx).GetProducts(eel.Map{
		"category_id": productCategory.Id,
		"is_deleted": false,
	})
	
	if len(products) == 0 {
		poolProducts := make([]*PoolProduct, 0)
		nextPageInfo := eel.MockPaginate(0, page)
		return poolProducts, nextPageInfo
	}
	
	productIds := make([]int, 0)
	for _, product := range products {
		productIds = append(productIds, product.Id)
	}
	
	filters["product_id__in"] = productIds
	filters["status"] = m_product.PP_STATUS_ON
	return this.SearchProducts(filters, page)
}

func (this *ProductPool) AddProductsByProducts(products []*Product)  {
	o := eel.GetOrmFromContext(this.Ctx)
	//创建新的PoolProduct记录
	newModels := make([]*m_product.PoolProduct, 0)
	corpId := this.Corp.GetId()
	nowTime := time.Now()
	for _, product := range products {
		model := m_product.PoolProduct{}
		model.CorpId = corpId
		model.ProductId = product.Id
		model.ProductType = product.Type
		model.SupplierId = product.CorpId
		model.Status = m_product.PP_STATUS_OFF
		model.Type = m_product.PP_TYPE_CREATE
		model.DisplayIndex = NEW_PRODUCT_DISPLAY_INDEX
		model.SyncAt = nowTime
		newModels = append(newModels, &model)
	}
	if len(newModels) > 0 {
		//TODO: 替换为BatchInsert方案
		for _, newModel := range newModels {
			db := o.Create(&newModel)
			if db.Error != nil {
				eel.Logger.Error(db.Error)
				panic(eel.NewBusinessError("pool_product:create_fail", "创建pool product失败"))
			}
		}
		//_, err := o.InsertMulti(len(newModels), newModels)
		//if err != nil {
		//	eel.Logger.Error(err)
		//	panic(eel.NewBusinessError("pool_product:create_fail", "创建pool product失败"))
		//}
	}
}

func (this *ProductPool) SyncPoolProducts(poolProducts []*PoolProduct)  {
	o := eel.GetOrmFromContext(this.Ctx)
	//创建新的PoolProduct记录
	newModels := make([]*m_product.PoolProduct, 0)
	corpId := this.Corp.GetId()
	nowTime := time.Now()
	for _, poolProduct := range poolProducts {
		model := m_product.PoolProduct{}
		model.CorpId = corpId
		model.ProductId = poolProduct.ProductId
		model.ProductType = poolProduct.ProductType
		model.SupplierId = poolProduct.SupplierId
		model.Status = m_product.PP_STATUS_OFF
		model.Type = m_product.PP_TYPE_SYNC
		model.SourcePoolProductId = poolProduct.Id
		model.DisplayIndex = NEW_PRODUCT_DISPLAY_INDEX
		model.SyncAt = nowTime
		newModels = append(newModels, &model)
	}
	if len(newModels) > 0 {
		//TODO: 替换为BatchInsert方案
		for _, newModel := range newModels {
			db := o.Create(&newModel)
			if db.Error != nil {
				eel.Logger.Error(db.Error)
				panic(eel.NewBusinessError("pool_product:sync_fail", "同步pool product失败"))
			}
		}
		//_, err := o.InsertMulti(len(newModels), newModels)
		//if err != nil {
		//	eel.Logger.Error(err)
		//	panic(eel.NewBusinessError("pool_product:sync_fail", "同步pool product失败"))
		//}
	}
}

func (this *ProductPool) GetPagedPromotionPoolProducts(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*PoolProduct, eel.INextPageInfo) {
	filters["corp_id"] = this.Corp.GetId()
	filters["is_distribution_by_corp"] = true
	filters["status"] = m_product.PP_STATUS_ON
	return NewPoolProductRepository(this.Ctx).GetPagedPoolProducts(filters, page, orderExprs...)
}


func (this *ProductPool) GetLowStockProducts() []*PoolProduct {
	o := eel.GetOrmFromContext(this.Ctx)
	
	//获得pool products对应的product id集合
	var models []*m_product.PoolProduct
	db := o.Model(&m_product.PoolProduct{}).Where(eel.Map{
		"corp_id": this.Corp.GetId(),
		"status": m_product.PP_STATUS_ON,
	}).Find(&models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
	
	productIds := make([]int, 0)
	for _, model := range models {
		productIds = append(productIds, model.ProductId)
	}
	
	//获得stock小于阈值的product id集合
	var skuModels []*m_product.ProductSku
	db = o.Model(&m_product.ProductSku{}).Where(eel.Map{
		"product_id__in": productIds,
		"stocks__lt": 200,
		"stocks__gt": -1,
	}).Find(&skuModels)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
	
	productIds = make([]int, 0)
	for _, skuModel := range skuModels {
		productIds = append(productIds, skuModel.ProductId)
	}
	
	
	return this.GetPoolProductsByProductIds(productIds)
}

// 获得在售商品的数量
func (this *ProductPool) GetOnsaleCount() int {
	o := eel.GetOrmFromContext(this.Ctx)

	count := 0
	db := o.Raw("select count(*) as count from product_pool_product where status = ? and corp_id = ?", m_product.PP_STATUS_ON, this.Corp.GetId()).Scan(&count)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
	
	return count
}

// GetLowStockProductCount 获得库存不足的商品数量
func (this *ProductPool) GetLowStockProductCount() int {
	o := eel.GetOrmFromContext(this.Ctx)
	var models []*m_product.PoolProduct
	db := o.Model(&m_product.PoolProduct{}).Where(eel.Map{
		"corp_id": this.Corp.GetId(),
		"status": m_product.PP_STATUS_ON,
	}).Find(&models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
	
	productIds := make([]int, 0)
	for _, model := range models {
		productIds = append(productIds, model.ProductId)
	}
	
	count, err := o.Model(&m_product.ProductSku{}).Where(eel.Map{
		"product_id__in": productIds,
		"stocks__lt": 200,
		"stocks__gt": -1,
	}).Count()
	if err != nil {
		eel.Logger.Error(err)
	}
	
	return int(count)
}

func GetProductPoolForCorp(ctx context.Context, corp business.ICorp) *ProductPool {
	instance := new(ProductPool)
	instance.Ctx = ctx
	instance.Corp = corp
	return instance
}

func GetGlobalProductPool(ctx context.Context) *ProductPool {
	instance := new(ProductPool)
	instance.Ctx = ctx
	instance.Corp = account.NewInvalidCorp(ctx)
	return instance
}

func init() {
}
