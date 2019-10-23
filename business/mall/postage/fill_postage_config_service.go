package postage

import (
	"context"
	"github.com/gingerxman/eel"
	
	
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
)

type FillPostageConfigService struct {
	eel.ServiceBase
}

func NewFillPostageConfigService(ctx context.Context) *FillPostageConfigService {
	service := new(FillPostageConfigService)
	service.Ctx = ctx
	return service
}

func (this *FillPostageConfigService) FillOne(postageConfig *PostageConfig, option eel.FillOption) {
	this.Fill([]*PostageConfig{ postageConfig }, option)
}

func (this *FillPostageConfigService) Fill(postageConfigs []*PostageConfig, option eel.FillOption) {
	if len(postageConfigs) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, postageConfig := range postageConfigs {
		ids = append(ids, postageConfig.Id)
	}

	this.fillSpecialAreaPostageConfig(postageConfigs, ids)
	this.fillFreePostageConfig(postageConfigs, ids)
}


func (this *FillPostageConfigService) fillSpecialAreaPostageConfig(postageConfigs []*PostageConfig, ids []int) {
	//从db中获取数据集合
	var models []*m_mall.SpecialAreaPostageConfig
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.SpecialAreaPostageConfig{}).Where("postage_config_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	//构建<id, postage>
	id2postage := make(map[int]*PostageConfig)
	for _, postage := range postageConfigs {
		id2postage[postage.Id] = postage
	}

	//填充postage_config的SpecialAreaPostageConfig对象
	for _, model := range models {
		if postageConfig, ok := id2postage[model.PostageConfigId]; ok {
			postageConfig.SpecialAreaPostageConfigs = append(postageConfig.SpecialAreaPostageConfigs, NewSpecialAreaPostageConfigFromModel(this.Ctx, model))
		}
	}
}



func (this *FillPostageConfigService) fillFreePostageConfig(postageConfigs []*PostageConfig, ids []int) {
	//从db中获取数据集合
	var models []*m_mall.FreePostageConfig
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.FreePostageConfig{}).Where("postage_config_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	//构建<id, postage>
	id2postage := make(map[int]*PostageConfig)
	for _, postage := range postageConfigs {
		id2postage[postage.Id] = postage
	}
	
	//填充postage_config的SpecialAreaPostageConfig对象
	for _, model := range models {
		if postageConfig, ok := id2postage[model.PostageConfigId]; ok {
			postageConfig.FreePostageConfigs = append(postageConfig.FreePostageConfigs, NewFreePostageConfigFromModel(this.Ctx, model))
		}
	}
}


func init() {
}
