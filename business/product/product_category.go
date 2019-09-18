package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)

type ProductCategory struct {
	eel.EntityBase
	Id int
	CorpId int
	Name string
	NodeType string
	ProductCount int
	DisplayIndex int
	IsEnabled bool
	CreatedAt time.Time
	
	//foreign key
	FatherId int
	Corp business.ICorp
}

//Update 更新对象
func (this *ProductCategory) Update(
	name string,
) error {
	var model m_product.ProductCategory
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"name": name,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_category:update_fail")
	}
	
	return nil
}

func (this *ProductCategory) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.ProductCategory{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *ProductCategory) Enable() {
	this.enable(true)
}

func (this *ProductCategory) Disable() {
	this.enable(false)
}

func (this *ProductCategory) GetSubCategories() []*ProductCategory {
	return NewProductCategoryRepository(this.Ctx).GetProductCategories(eel.Map{
		"father_id": this.Id,
		"corp_id": this.CorpId,
	})
}

func (this *ProductCategory) GetPagedProducts(filters eel.Map, page *eel.PageInfo, orderExprs string) ([]*PoolProduct, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	qs := o.Model(&m_product.Product{})
	
	var models []*m_product.Product
	qs = qs.Where(eel.Map{
		"category_id": this.Id,
	})
	if len(orderExprs) > 0 {
		qs = qs.Order(orderExprs)
	}
	paginateResult, err := eel.Paginate(qs, page, &models)
	
	if err != nil {
		eel.Logger.Error(err)
		return nil, paginateResult
	}
	
	//获取PoolProduct对象集合
	productIds := make([]int, 0)
	for _, model := range models {
		productIds = append(productIds, model.Id)
	}
	
	if len(productIds) == 0 {
		return make([]*PoolProduct, 0), eel.MockPaginate(0, page)
	}
	
	productPool := GetProductPoolForCorp(this.Ctx, this.Corp)
	filters["product_id__in"] = productIds
	//filters["corp_id"] = this.CorpId
	poolProducts := productPool.GetPoolProducts(filters, "id desc")
	return poolProducts, paginateResult
}

//工厂方法
func NewProductCategory(
	ctx context.Context,
	corp business.ICorp,
	name string,
	fatherId int,
) *ProductCategory {
	o := eel.GetOrmFromContext(ctx)
	
	if fatherId != 0 {
		//需要将fatherId指定的node的type从"leaf"修改为"node"
		db := o.Model(&m_product.ProductCategory{}).Where("id", fatherId).Update(gorm.Params{
			"node_type": m_product.PRODUCT_CATEGORY_NODE_TYPE_NODE,
		})
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("product_category:update_father_node_type_fail", fmt.Sprintf("更新父节点类型失败")))
		}
	}
	
	model := m_product.ProductCategory{}
	model.CorpId = corp.GetId()
	model.IsEnabled = true
	model.Name = name
	model.FatherId = fatherId
	model.NodeType = m_product.PRODUCT_CATEGORY_NODE_TYPE_LEAF
	
	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_category:create_fail", fmt.Sprintf("创建失败")))
	}
	
	return NewProductCategoryFromModel(ctx, &model)
}

//根据model构建对象
func NewProductCategoryFromModel(ctx context.Context, model *m_product.ProductCategory) *ProductCategory {
	instance := new(ProductCategory)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Name = model.Name
	instance.NodeType = m_product.PRODUCTCATEGORYNODE2STR[model.NodeType]
	instance.FatherId = model.FatherId
	instance.ProductCount = model.ProductCount
	instance.DisplayIndex = model.DisplayIndex
	instance.IsEnabled = model.IsEnabled
	instance.CreatedAt = model.CreatedAt
	
	return instance
}

//根据model构建对象
func NewRootProductCategory(ctx context.Context, corp business.ICorp) *ProductCategory {
	instance := new(ProductCategory)
	instance.Ctx = ctx
	instance.Id = 0
	instance.FatherId = 0
	if corp != nil {
		instance.CorpId = corp.GetId()
		instance.Corp = corp
	}
	
	return instance
}

func init() {
}
