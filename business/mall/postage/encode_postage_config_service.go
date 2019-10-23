package postage

import (
	"context"
	"encoding/json"
	
	"strconv"
	"strings"
	"sync"
	
	"github.com/gingerxman/eel"
)

type EncodePostageConfigService struct {
	eel.ServiceBase
	sync.RWMutex
	id2province map[int]*RProvince
}

func NewEncodePostageConfigService(ctx context.Context) *EncodePostageConfigService {
	service := new(EncodePostageConfigService)
	service.Ctx = ctx
	return service
}

func (this *EncodePostageConfigService) getNameProvinceMap() map[int]*RProvince {
	if this.id2province != nil {
		return this.id2province
	}
	
	this.Lock()
	defer this.Unlock()
	
	resp, err := eel.NewResource(this.Ctx).Get("gaia", "area.provinces", eel.Map{
	})
	
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("encode_postage_config_service:get_gaia_provinces_fail", "获取gaia的area.provinces失败"))
	}
	
	this.id2province = make(map[int]*RProvince)
	data, _ := resp.Data().Map()
	provinceDatas := data["provinces"].([]interface{})
	for _, provinceData := range provinceDatas {
		data := provinceData.(map[string]interface{})
		id, err := data["id"].(json.Number).Int64()
		if err != nil {
			eel.Logger.Error(err)
			continue
		}
		
		name := data["name"].(string)
		this.id2province[int(id)] = &RProvince{
			Id: int(id),
			Name: name,
		}
	}
	
	return this.id2province
}

func (this *EncodePostageConfigService) parseDestinations(destination string) []*RProvince {
	rProvinces := make([]*RProvince, 0)
	if destination == "" {
		return rProvinces
	}
	
	id2province := this.getNameProvinceMap()
	items := strings.Split(destination, ",")
	for _, item := range items {
		id, _ := strconv.Atoi(item)
		if province, ok := id2province[id]; ok {
			rProvinces = append(rProvinces, province)
		}
	}
	return rProvinces
}

//Encode 对单个实体对象进行编码
func (this *EncodePostageConfigService) Encode(postageConfig *PostageConfig) *RPostageConfig {
	if postageConfig == nil {
		return nil
	}
	
	rDefaultPostageConfig := &RDefaultPostageConfig{
		FirstWeight: postageConfig.FirstWeight,
		FirstWeightPrice: postageConfig.FirstWeightPrice,
		AddedWeight: postageConfig.AddedWeight,
		AddedWeightPrice: postageConfig.AddedWeightPrice,
	}
	
	rSpecialAreaPostageConfigs := make([]*RSpecialAreaPostageConfig, 0)
	for _, config := range postageConfig.SpecialAreaPostageConfigs {
		rSpecialAreaPostageConfigs = append(rSpecialAreaPostageConfigs, &RSpecialAreaPostageConfig{
			FirstWeight: config.FirstWeight,
			FirstWeightPrice: config.FirstWeightPrice,
			AddedWeight: config.AddedWeight,
			AddedWeightPrice: config.AddedWeightPrice,
			Destination: config.Destination,
			DestinationProvinces: this.parseDestinations(config.Destination),
		})
	}
	
	rFreePostageConfigs := make([]*RFreePostageConfig, 0)
	for _, config := range postageConfig.FreePostageConfigs {
		rFreePostageConfigs = append(rFreePostageConfigs, &RFreePostageConfig{
			Condition: config.Condition,
			ConditionValue: config.ConditionValue,
			Destination: config.Destination,
			DestinationProvinces: this.parseDestinations(config.Destination),
		})
	}

	return &RPostageConfig{
		Id: postageConfig.Id,
		CorpId: postageConfig.CorpId,
		Name: postageConfig.Name,
		IsUsed: postageConfig.IsUsed,
		IsSystemLevelConfig: postageConfig.IsSystemLevelConfig,
		IsEnableSpecialConfig: postageConfig.IsEnableSpecialConfig,
		IsEnableFreeConfig: postageConfig.IsEnableFreeConfig,
		DefaultPostageConfig: rDefaultPostageConfig,
		SpecialAreaPostageConfigs: rSpecialAreaPostageConfigs,
		FreePostageConfigs: rFreePostageConfigs,
		IsEnabled: postageConfig.IsEnabled,
		IsDeleted: postageConfig.IsDeleted,
		CreatedAt: postageConfig.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodePostageConfigService) EncodeMany(postageConfigs []*PostageConfig) []*RPostageConfig {
	rDatas := make([]*RPostageConfig, 0)
	for _, postageConfig := range postageConfigs {
		rDatas = append(rDatas, this.Encode(postageConfig))
	}
	
	return rDatas
}

func init() {
}
