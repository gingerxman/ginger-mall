package postage

import (
	"context"
	"encoding/json"
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/ginger-mall/business"
	"strconv"
	
	"github.com/gingerxman/eel"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
)

type DefaultPostage struct {
	FirstWeight string `json:"first_weight"`
	FirstWeightPrice string `json:"first_weight_price"`
	AddedWeight string `json:"added_weight"`
	AddedWeightPrice string `json:"added_weight_price"`
}

type SpecialAreaPostage struct {
	Destinations string `json:"destinations"`
	FirstWeight string `json:"first_weight"`
	FirstWeightPrice string `json:"first_weight_price"`
	AddedWeight string `json:"added_weight"`
	AddedWeightPrice string `json:"added_weight_price"`
}

type FreePostage struct {
	Destinations string `json:"destinations"`
	Condition string `json:"condition"`
	ConditionValue string `json:"value"`
}

func parsePostageConfigs(defaultPostageStr string, isEnableSpecialConfig bool, specialAreaPostagesStr string, isEnableFreeConfig bool, freePostagesStr string) (*DefaultPostage, []*SpecialAreaPostage, []*FreePostage, error) {
	defaultPostage := DefaultPostage{}
	err := json.Unmarshal([]byte(defaultPostageStr), &defaultPostage)
	if err != nil {
		return nil, nil, nil, err
	}
	
	specialAreaPostages := make([]*SpecialAreaPostage, 0)
	if isEnableSpecialConfig {
		err = json.Unmarshal([]byte(specialAreaPostagesStr), &specialAreaPostages)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	
	freePostages := make([]*FreePostage, 0)
	if isEnableFreeConfig {
		err = json.Unmarshal([]byte(freePostagesStr), &freePostages)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	
	return &defaultPostage, specialAreaPostages, freePostages, nil
}

func saveSpecialPostageConfigs(o *gorm.DB, corp business.ICorp, postageConfigId int, postages []*SpecialAreaPostage) error {
	//删除旧数据
	db := o.Where(eel.Map{
		"postage_config_id": postageConfigId,
	}).Delete(&m_mall.SpecialAreaPostageConfig{})
	if db.Error != nil {
		return db.Error
	}
	
	//保存新数据
	for _, postage := range postages {
		firstWeight, _ := strconv.ParseFloat(postage.FirstWeight, 64)
		firstWeightPrice, _ := strconv.ParseFloat(postage.FirstWeightPrice, 64)
		addedWeight, _ := strconv.ParseFloat(postage.AddedWeight, 64)
		addedWeightPrice, _ := strconv.ParseFloat(postage.AddedWeightPrice, 64)
		model := m_mall.SpecialAreaPostageConfig{
			CorpId: corp.GetId(),
			PostageConfigId: postageConfigId,
			Destination: postage.Destinations,
			FirstWeight: firstWeight,
			FirstWeightPrice: firstWeightPrice,
			AddedWeight: addedWeight,
			AddedWeightPrice: addedWeightPrice,
		}
		
		db = o.Create(&model)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func saveFreePostageConfigs(o *gorm.DB, corp business.ICorp, postageConfigId int, postages []*FreePostage) error {
	//删除旧数据
	db := o.Where(eel.Map{
		"postage_config_id": postageConfigId,
	}).Delete(&m_mall.FreePostageConfig{})
	if db.Error != nil {
		return db.Error
	}
	
	//保存新数据
	for _, postage := range postages {
		model := m_mall.FreePostageConfig{
			CorpId: corp.GetId(),
			PostageConfigId: postageConfigId,
			Destination: postage.Destinations,
			Condition: postage.Condition,
			ConditionValue: postage.ConditionValue,
		}
		
		db := o.Create(&model)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

type PostageConfigService struct {
	eel.ServiceBase
}

func NewPostageConfigService(ctx context.Context) *PostageConfigService {
	service := new(PostageConfigService)
	service.Ctx = ctx
	return service
}

func (this *PostageConfigService) UpdatePostageConfigForCorp(corp business.ICorp, id int, name string, defaultPostageStr string, isEnableSpecialConfig bool, specialAreaPostagesStr string, isEnableFreeConfig bool, freePostagesStr string) error {
	defaultPostage, specialAreaPostages, freePostages, err := parsePostageConfigs(defaultPostageStr, isEnableSpecialConfig, specialAreaPostagesStr, isEnableFreeConfig, freePostagesStr)
	
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	//更新default postage
	{
		firstWeight, _ := strconv.ParseFloat(defaultPostage.FirstWeight, 64)
		firstWeightPrice, _ := strconv.ParseFloat(defaultPostage.FirstWeightPrice, 64)
		addedWeight, _ := strconv.ParseFloat(defaultPostage.AddedWeight, 64)
		addedWeightPrice, _ := strconv.ParseFloat(defaultPostage.AddedWeightPrice, 64)
		db := o.Model(&m_mall.PostageConfig{}).Where("id", id).Where("corp_id", corp.GetId()).Update(gorm.Params{
			"name": name,
			"first_weight": firstWeight,
			"first_weight_price": firstWeightPrice,
			"added_weight": addedWeight,
			"added_weight_price": addedWeightPrice,
			"is_enable_special_config": isEnableSpecialConfig,
			"is_enable_free_config": isEnableFreeConfig,
		})
		if db.Error != nil {
			return db.Error
		}
	}
	
	err = saveSpecialPostageConfigs(o, corp, id, specialAreaPostages)
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	err = saveFreePostageConfigs(o, corp, id, freePostages)
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	return nil
}

func (this *PostageConfigService) SetDefaultPostageConfigForCorp(corp business.ICorp, id int) error {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.PostageConfig{}).Where(eel.Map{
		"corp_id": corp.GetId(),
		"is_used": true,
	}).Update(gorm.Params{
		"is_used": false,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	db = o.Model(&m_mall.PostageConfig{}).Where(eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}).Update(gorm.Params{
		"is_used": true,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}

func init() {
}
