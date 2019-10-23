package postage

type RProvince struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type RPostageConfig struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	Name string `json:"name"`
	IsUsed bool `json:"is_used"`
	IsSystemLevelConfig bool `json:"is_system_level_config"`
	IsEnableSpecialConfig bool `json:"is_enable_special_config"`
	IsEnableFreeConfig bool `json:"is_enable_free_config"`
	DefaultPostageConfig *RDefaultPostageConfig `json:"default_config"`
	SpecialAreaPostageConfigs []*RSpecialAreaPostageConfig `json:"special_configs"`
	FreePostageConfigs []*RFreePostageConfig `json:"free_configs"`
	IsEnabled bool `json:"is_enabled"`
	IsDeleted bool `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
}

type RDefaultPostageConfig struct {
	FirstWeight float64 `json:"first_weight"`
	FirstWeightPrice float64 `json:"first_weight_price"`
	AddedWeight float64 `json:"added_weight"`
	AddedWeightPrice float64 `json:"added_weight_price"`
}

type RSpecialAreaPostageConfig struct {
	FirstWeight float64 `json:"first_weight"`
	FirstWeightPrice float64 `json:"first_weight_price"`
	AddedWeight float64 `json:"added_weight"`
	AddedWeightPrice float64 `json:"added_weight_price"`
	Destination string `json:"destination"`
	DestinationProvinces []*RProvince `json:"destination_provinces"`
}

type RFreePostageConfig struct {
	Destination string `json:"destination"`
	DestinationProvinces []*RProvince `json:"destination_provinces"`
	Condition string `json:"condition"`
	ConditionValue string `json:"value"`
}


func init() {
}
