package postage

import (
	"context"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
	"time"
	
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
)

type PostageConfig struct {
	eel.EntityBase
	Id int
	CorpId int
	Name string
	FirstWeight float64
	FirstWeightPrice float64
	AddedWeight float64
	AddedWeightPrice float64
	IsUsed bool
	IsSystemLevelConfig bool
	IsEnableSpecialConfig bool
	IsEnableFreeConfig bool
	IsEnabled bool
	IsDeleted bool
	CreatedAt time.Time

	//foreign key
	SpecialAreaPostageConfigId int //refer to special_area_postage_config
	SpecialAreaPostageConfigs []*SpecialAreaPostageConfig
	FreePostageConfigId int //refer to free_postage_config
	FreePostageConfigs []*FreePostageConfig
}



func (this *PostageConfig) enable(isEnabled bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.PostageConfig{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_enabled": isEnabled,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *PostageConfig) Enable() {
	this.enable(true)
}

func (this *PostageConfig) Disable() {
	this.enable(false)
}

func (this *PostageConfig) Delete() error {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.PostageConfig{}).Where("id", this.Id).Update(gorm.Params{
		"is_deleted": true,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}

	return nil
}

func (this *PostageConfig) SetUsed() error {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.PostageConfig{}).Where("corp_id", this.CorpId).Update(gorm.Params{
		"is_used": false,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	db = o.Model(&m_mall.PostageConfig{}).Where("id", this.Id).Update(gorm.Params{
		"is_used": true,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}


func (this *PostageConfig) UpdateDisplayIndex(action string) error {
	//model := m_mall.PostageConfig{}
	//item := itemPos{
	//	Id: this.Id,
	//	Table: model.TableName(),
	//}
	//err := NewUpdateDisplayIndexService(this.Ctx, DISPLAY_INDEX_ORDER_ASC).Update(&item, action)
	//if err != nil {
	//	eel.Logger.Error(err)
	//	return err
	//}
	
	return nil
}

//根据model构建对象
func NewPostageConfigFromModel(ctx context.Context, model *m_mall.PostageConfig) *PostageConfig {
	instance := new(PostageConfig)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Name = model.Name
	instance.FirstWeight = model.FirstWeight
	instance.FirstWeightPrice = model.FirstWeightPrice
	instance.AddedWeight = model.AddedWeight
	instance.AddedWeightPrice = model.AddedWeightPrice
	instance.IsUsed = model.IsUsed
	instance.IsSystemLevelConfig = model.IsSystemLevelConfig
	instance.IsEnableSpecialConfig = model.IsEnableSpecialConfig
	instance.IsEnableFreeConfig = model.IsEnableFreeConfig
	instance.SpecialAreaPostageConfigId = model.SpecialAreaPostageConfigId
	instance.FreePostageConfigId = model.FreePostageConfigId
	instance.IsEnabled = model.IsEnabled
	instance.IsDeleted = model.IsDeleted
	instance.CreatedAt = model.CreatedAt
	
	instance.SpecialAreaPostageConfigs = make([]*SpecialAreaPostageConfig, 0)
	instance.FreePostageConfigs = make([]*FreePostageConfig, 0)

	return instance
}

func init() {
}
