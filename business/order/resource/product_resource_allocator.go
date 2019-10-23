package resource

import (
	"context"
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
	
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
)

type ProductResourceAllocator struct {
	eel.ServiceBase
}

func NewProductResourceAllocator(ctx context.Context) business.IResourceAllocator {
	service := new(ProductResourceAllocator)
	service.Ctx = ctx
	return service
}

//Allocate 申请商品资源，减少库存
func (this *ProductResourceAllocator) Allocate(resource business.IResource, newOrder business.IOrder) error {
	productResource := resource.(*ProductResource)
	
	o := eel.GetOrmFromContext(this.Ctx)
	sku := productResource.GetPoolProduct().GetSku(productResource.Sku)
	
	db := o.Model(&m_product.ProductSku{}).Where("id", sku.Id).Update("stocks", gorm.Expr("stocks - ?", productResource.Count))
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}

//Release 释放商品资源，恢复库存
func (this *ProductResourceAllocator) Release(resource business.IResource) {
	productResource := resource.(*ProductResource)
	
	o := eel.GetOrmFromContext(this.Ctx)
	sku := productResource.GetPoolProduct().GetSku(productResource.Sku)
	
	db := o.Model(&m_product.ProductSku{}).Where("id", sku.Id).Update("stocks", gorm.Expr("stocks + ?", productResource.Count))
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}


func init() {
}
