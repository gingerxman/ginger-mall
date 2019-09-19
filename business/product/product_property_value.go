package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)

type ProductPropertyValue struct {
	eel.EntityBase
	Id int
	PropertyId int
	PropertyName string
	Text string
	Image string
	IsDeleted bool
	CreatedAt time.Time
}

//Update 更新对象
func (this *ProductPropertyValue) Update(
	text string,
	picUrl string,
) error {
	var model m_product.ProductPropertyValue
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"text": text,
		"image": picUrl,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_property_value:update_fail")
	}

	return nil
}

func (this *ProductPropertyValue) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.ProductPropertyValue{}).Where(eel.Map{
		"id": this.Id,
		"property_id": this.PropertyId,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *ProductPropertyValue) Enable() {
	this.enable(true)
}

func (this *ProductPropertyValue) Disable() {
	this.enable(false)
}

func (this *ProductPropertyValue) GetFullId() string {
	return fmt.Sprintf("%d:%d", this.PropertyId, this.Id)
}

//工厂方法
func NewProductPropertyValue(
	ctx context.Context, 
	propertyId int,
	text string,
	image string,
) *ProductPropertyValue {
	o := eel.GetOrmFromContext(ctx)
	model := m_product.ProductPropertyValue{}
	model.PropertyId = propertyId
	model.Text = text
	model.Image = image
	model.IsDeleted = false

	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_property_value:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewProductPropertyValueFromModel(ctx, &model)
}

//根据model构建对象
func NewProductPropertyValueFromModel(ctx context.Context, model *m_product.ProductPropertyValue) *ProductPropertyValue {
	instance := new(ProductPropertyValue)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.PropertyId = model.PropertyId
	instance.Text = model.Text
	instance.Image = model.Image
	instance.IsDeleted = model.IsDeleted
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
