package product

import (
	"fmt"
	"errors"
	"context"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
)

type ProductUsableImoney struct {
	eel.EntityBase
	Id int
	ProductId int
	ImoneyCode string
	IsEnabled bool

	//foreign key
}

//Update 更新对象
func (this *ProductUsableImoney) Update(
	productId int,
	imoneyCode string,
) error {
	var model m_product.ProductUsableImoney
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"product_id": productId,
		"imoney_code": imoneyCode,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_usable_imoney:update_fail")
	}

	return nil
}

func (this *ProductUsableImoney) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.ProductUsableImoney{}).Where(eel.Map{
		"id": this.Id,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *ProductUsableImoney) Enable() {
	this.enable(true)
}

func (this *ProductUsableImoney) Disable() {
	this.enable(false)
}

//工厂方法
func NewProductUsableImoney(
	ctx context.Context, 
	productId int,
	imoneyCode string,
) *ProductUsableImoney {
	o := eel.GetOrmFromContext(ctx)
	model := m_product.ProductUsableImoney{}
	model.IsEnabled = true
	model.ProductId = productId
	model.ImoneyCode = imoneyCode
	

	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_usable_imoney:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewProductUsableImoneyFromModel(ctx, &model)
}

//根据model构建对象
func NewProductUsableImoneyFromModel(ctx context.Context, model *m_product.ProductUsableImoney) *ProductUsableImoney {
	instance := new(ProductUsableImoney)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.ProductId = model.ProductId
	instance.ImoneyCode = model.ImoneyCode
	instance.IsEnabled = model.IsEnabled

	return instance
}

func init() {
}
