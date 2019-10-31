package ship_info

import (
	"context"
	"errors"
	"fmt"
	"github.com/gingerxman/ginger-mall/business"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
	"time"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/gorm"
)

type sAreaItem struct {
	Id int `json:"int"`
	Name string `json:"name"`
}

type sAreaInfo struct {
	Province sAreaItem `json:"province"`
	City sAreaItem `json:"city"`
	District sAreaItem `json:"district"`
}

type ShipInfo struct {
	eel.EntityBase
	Id int
	UserId int
	Name string
	Phone string
	Area string
	AreaCode string
	AreaJson string
	Address string
	IsDefault bool
	IsEnabled bool
	IsDeleted bool
	CreatedAt time.Time

	//foreign key
}

//Update 更新对象
func (this *ShipInfo) Update(
	name string,
	phone string,
	areaCode string,
	address string,
) error {
	var model m_mall.ShipInfo
	o := eel.GetOrmFromContext(this.Ctx)
	
	area := eel.NewAreaService().GetAreaByCode(areaCode)
	
	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"name": name,
		"phone": phone,
		"area": fmt.Sprintf("%s %s %s", area.Province.Name, area.City.Name, area.District.Name),
		"area_code": areaCode,
		"address": address,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("ship_info:update_fail")
	}

	return nil
}

func (this *ShipInfo) SetDefault(user business.IUser) error {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.ShipInfo{}).Where(eel.Map{
		"user_id": user.GetId(),
		"is_default": true,
	}).Update(gorm.Params{
		"is_default": false,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	db = o.Model(&m_mall.ShipInfo{}).Where(eel.Map{
		"id": this.Id,
	}).Update(gorm.Params{
		"is_default": true,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}

func (this *ShipInfo) enable(isEnabled bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.ShipInfo{}).Where(eel.Map{
		"id": this.Id,
	}).Update(gorm.Params{
		"is_enabled": isEnabled,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *ShipInfo) Enable() {
	this.enable(true)
}

func (this *ShipInfo) Disable() {
	this.enable(false)
}

func (this *ShipInfo) Delete() error {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.ShipInfo{}).Where("id", this.Id).Update(gorm.Params{
		"is_deleted": true,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}

	return nil
}



//工厂方法
func NewShipInfo(
	ctx context.Context, 
	user business.IUser,
	name string,
	phone string,
	areaCode string,
	address string,
) *ShipInfo {
	o := eel.GetOrmFromContext(ctx)
	
	area := eel.NewAreaService().GetAreaByCode(areaCode)

	//保存数据
	model := m_mall.ShipInfo{}
	model.IsEnabled = true
	model.IsDeleted = false
	model.UserId = user.GetId()
	model.Name = name
	model.Phone = phone
	model.Area = fmt.Sprintf("%s %s %s", area.Province.Name, area.City.Name, area.District.Name)
	model.AreaCode = areaCode
	model.AreaJson = ""
	model.Address = address
	model.IsDefault = false
	
	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("ship_info:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewShipInfoFromModel(ctx, &model)
}

//根据model构建对象
func NewShipInfoFromModel(ctx context.Context, model *m_mall.ShipInfo) *ShipInfo {
	instance := new(ShipInfo)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.UserId = model.UserId
	instance.Name = model.Name
	instance.Phone = model.Phone
	instance.Area = model.Area
	instance.AreaCode = model.AreaCode
	instance.AreaJson = model.AreaJson
	instance.Address = model.Address
	instance.IsDefault = model.IsDefault
	instance.IsEnabled = model.IsEnabled
	instance.IsDeleted = model.IsDeleted
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
