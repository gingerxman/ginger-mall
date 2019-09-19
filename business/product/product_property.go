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

type ProductProperty struct {
	eel.EntityBase
	Id int
	CorpId int
	Name string
	IsDeleted bool
	CreatedAt time.Time

	//foreign key
	Values []*ProductPropertyValue
}

//Update 更新对象
func (this *ProductProperty) Update(
	name string,
) error {
	var model m_product.ProductProperty
	o := eel.GetOrmFromContext(this.Ctx)
	
	updateParams := gorm.Params{}
	if name != "" {
		updateParams["name"] = name
	}

	db := o.Model(&model).Where("id", this.Id).Update(updateParams)

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_property:update_fail")
	}

	return nil
}

func (this *ProductProperty) Delete() error{
	this.checkUsedBySku()

	var model m_product.ProductProperty
	o := eel.GetOrmFromContext(this.Ctx)

	updateParams := gorm.Params{}
	updateParams["is_deleted"] = true

	db := o.Model(&model).Where("id", this.Id).Update(updateParams)

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_property:deleted_fail")
	}
	err := this.deleteValues()
	if err != nil{
		return err
	}

	return nil
}
func (this *ProductProperty) deleteValues() error{
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&m_product.ProductPropertyValue{}).Where(eel.Map{
		"is_deleted": false,
		"property_id": this.Id,
	}).Update(gorm.Params{"is_deleted": true})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_property:deleted_fail")
	}
	return nil
}
func (this *ProductProperty) checkUsedBySku(){
	o := eel.GetOrmFromContext(this.Ctx)

	var models []*m_product.ProductSkuHasPropertyValue
	db := o.Model(&m_product.ProductSkuHasPropertyValue{}).Where(eel.Map{
		"property_id": this.Id,
	}).Find(&models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	if len(models) != 0{
		panic(eel.NewBusinessError("product_property:cannot deleted", "删除商品属性失败"))
	}
}
func (this *ProductProperty) checkValueUsedBySku(propertyValueId int){
	o := eel.GetOrmFromContext(this.Ctx)

	var models []*m_product.ProductSkuHasPropertyValue
	db := o.Model(&m_product.ProductSkuHasPropertyValue{}).Where(eel.Map{
		"property_id": this.Id,
		"property_value_id": propertyValueId,
	}).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	if len(models) != 0{
		panic(eel.NewBusinessError("product_property_value:cannot deleted", "删除商品属性值失败"))
	}
}

func (this *ProductProperty) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.ProductProperty{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_deleted": !isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *ProductProperty) Enable() {
	this.enable(true)
}

func (this *ProductProperty) Disable() {
	this.enable(false)
}

func (this *ProductProperty) AddNewValue(text string, image string) *ProductPropertyValue {
	return NewProductPropertyValue(
		this.Ctx,
		this.Id,
		text,
		image,
	)
}

func (this *ProductProperty) DeleteValue(valueId int) {
	this.checkValueUsedBySku(valueId)
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.ProductPropertyValue{}).Where(eel.Map{
		"id": valueId,
		"property_id": this.Id,
	}).Update(gorm.Params{
		"is_deleted": true,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_property:delete_value_fail", "删除商品属性值失败"))
	}
}

func (this *ProductProperty) AppendValue(value *ProductPropertyValue) {
	this.Values = append(this.Values, value)
}

//工厂方法
func NewProductProperty(
	ctx context.Context, 
	corp business.ICorp,
	name string,
) *ProductProperty {
	o := eel.GetOrmFromContext(ctx)
	model := m_product.ProductProperty{}
	model.CorpId = corp.GetId()
	model.Name = name
	model.IsDeleted = false
	
	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_property:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewProductPropertyFromModel(ctx, &model)
}

//根据model构建对象
func NewProductPropertyFromModel(ctx context.Context, model *m_product.ProductProperty) *ProductProperty {
	instance := new(ProductProperty)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Name = model.Name
	instance.IsDeleted = model.IsDeleted
	instance.CreatedAt = model.CreatedAt
	
	instance.Values = make([]*ProductPropertyValue, 0)

	return instance
}

func init() {
}
