package postage

import (
	"context"
	
	"github.com/gingerxman/ginger-mall/business"
	"strconv"
	
	"github.com/gingerxman/eel"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
)

type PostageConfigFactory struct {
	eel.ServiceBase
}

func NewPostageConfigFactory(ctx context.Context) *PostageConfigFactory {
	service := new(PostageConfigFactory)
	service.Ctx = ctx
	return service
}

func (this *PostageConfigFactory) CreatePostageConfigForCorp(corp business.ICorp, name string, defaultPostageStr string, isEnableSpecialConfig bool, specialAreaPostagesStr string, isEnableFreeConfig bool, freePostagesStr string) (*PostageConfig, error) {
	defaultPostage, specialAreaPostages, freePostages, err := parsePostageConfigs(defaultPostageStr, isEnableSpecialConfig, specialAreaPostagesStr, isEnableFreeConfig, freePostagesStr)
	
	if err != nil {
		eel.Logger.Error(err)
		return nil, err
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	var postageConfig *PostageConfig
	//存储default postage
	firstWeight, _ := strconv.ParseFloat(defaultPostage.FirstWeight, 64)
	firstWeightPrice, _ := strconv.ParseFloat(defaultPostage.FirstWeightPrice, 64)
	addedWeight, _ := strconv.ParseFloat(defaultPostage.AddedWeight, 64)
	addedWeightPrice, _ := strconv.ParseFloat(defaultPostage.AddedWeightPrice, 64)
	model := m_mall.PostageConfig{
		CorpId: corp.GetId(),
		Name: name,
		FirstWeight: firstWeight,
		FirstWeightPrice: firstWeightPrice,
		AddedWeight: addedWeight,
		AddedWeightPrice: addedWeightPrice,
		IsEnableSpecialConfig: isEnableSpecialConfig,
		IsEnableFreeConfig: isEnableFreeConfig,
		IsEnabled: true,
		IsUsed: false,
		IsSystemLevelConfig: false,
	}
	
	db := o.Create(&model)
	if db.Error != nil {
		return nil, db.Error
	}
	postageConfig = NewPostageConfigFromModel(this.Ctx, &model)
	
	err = saveSpecialPostageConfigs(o, corp, model.Id, specialAreaPostages)
	if err != nil {
		eel.Logger.Error(err)
		return nil, err
	}
	
	err = saveFreePostageConfigs(o, corp, model.Id, freePostages)
	if err != nil {
		eel.Logger.Error(err)
		return nil, err
	}
	
	return postageConfig, nil
}

func (this *PostageConfigFactory) MakeSureDefaultPostageConfigExits(corp business.ICorp) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	if !o.Model(m_mall.PostageConfig{}).Where("corp_id", corp.GetId()).Exist() {
		model := m_mall.PostageConfig{
			CorpId:              corp.GetId(),
			Name:                "免运费",
			IsUsed:              true,
			IsSystemLevelConfig: true,
			IsEnabled: true,
		}
		db := o.Create(&model)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
		}
	}
}

func init() {
}
