package product

import (
	"fmt"
	"errors"
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"github.com/gingerxman/gorm"
	"time"
)

type ProductLabel struct {
	eel.EntityBase
	Id int
	CorpId int
	Name string
	IsEnabled bool
	CreatedAt time.Time

	//foreign key
}

//Update 更新对象
func (this *ProductLabel) Update(
	name string,
) error {
	var model m_product.ProductLabel
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"name": name,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("product_label:update_fail")
	}

	return nil
}

func (this *ProductLabel) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_product.ProductLabel{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *ProductLabel) Enable() {
	this.enable(true)
}

func (this *ProductLabel) Disable() {
	this.enable(false)
}

//工厂方法
func NewProductLabel(
	ctx context.Context, 
	corp business.ICorp,
	name string,
) *ProductLabel {
	o := eel.GetOrmFromContext(ctx)
	model := m_product.ProductLabel{}
	model.CorpId = corp.GetId()
	model.IsEnabled = true
	model.Name = name

	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("product_label:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewProductLabelFromModel(ctx, &model)
}

//根据model构建对象
func NewProductLabelFromModel(ctx context.Context, model *m_product.ProductLabel) *ProductLabel {
	instance := new(ProductLabel)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Name = model.Name
	instance.IsEnabled = model.IsEnabled
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
