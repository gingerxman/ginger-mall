package product

import (
	"github.com/gingerxman/eel"
)

type Count struct {
	eel.Model
	Type string `gorm:"size:125"`
	CorpId int
	Count int
}
func (this *Count) TableName() string {
	return "mall_count"
}
func (this *Count) TableIndex() [][]string {
	return [][]string{
		[]string{"Type", "CorpId"},
	}
}

//PostageConfig Model
type PostageConfig struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	SpecialAreaPostageConfigId int //foreign key special_area_postage_config
	FreePostageConfigId int //foreign key free_postage_config
	Name string
	FirstWeight float64
	FirstWeightPrice float64
	AddedWeight float64
	AddedWeightPrice float64
	IsUsed bool `gorm:"default:false"`
	IsSystemLevelConfig bool `gorm:"default:false"`
	IsEnableSpecialConfig bool `gorm:"default:false"`
	IsEnableFreeConfig bool `gorm:"default:false"`
	IsEnabled bool
	IsDeleted bool
}
func (self *PostageConfig) TableName() string {
	return "mall_postage_config"
}
func (this *PostageConfig) TableIndex() [][]string {
	return [][]string{
		[]string{"CorpId", "IsDeleted", "IsEnabled"},
	}
}


//LimitZone Model
type LimitZone struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	Name string `gorm:"size:64"`
	Areas string `gorm:"type:text"`
}
func (self *LimitZone) TableName() string {
	return "mall_limit_zone"
}


//SpecialAreaPostageConfig Model
type SpecialAreaPostageConfig struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	PostageConfigId int `gorm:"index"`
	FirstWeight float64
	FirstWeightPrice float64
	AddedWeight float64
	AddedWeightPrice float64
	Destination string
}
func (self *SpecialAreaPostageConfig) TableName() string {
	return "mall_special_area_postage_config"
}


//FreePostageConfig Model
type FreePostageConfig struct {
	eel.Model
	CorpId          int `gorm:"index"` //foreign key for corp
	PostageConfigId int `gorm:"index"`
	Destination     string
	Condition       string
	ConditionValue  string
}
func (self *FreePostageConfig) TableName() string {
	return "mall_free_postage_config"
}


//ShipInfo Model
type ShipInfo struct {
	eel.Model
	UserId int
	Name string
	Phone string
	Area string `gorm:"size:256"`
	AreaCode string `gorm:"size:256"`
	AreaJson string `gorm:"size:2048"`
	Address string
	IsDefault bool `gorm:"default:false"`
	IsEnabled bool
	IsDeleted bool
}
func (self *ShipInfo) TableName() string {
	return "mall_ship_info"
}
func (this *ShipInfo) TableIndex() [][]string {
	return [][]string{
		[]string{"UserId", "IsDeleted", "IsEnabled"},
		[]string{"UserId", "IsDefault"},
	}
}


//ShoppingCartItem Model
type ShoppingCartItem struct {
	eel.Model
	UserId int
	CorpId int
	ProductId int
	ProductSkuName string `gorm:"size:256"`
	ProductSkuDisplayName string `gorm:"size:256"`
	Count int
}
func (self *ShoppingCartItem) TableName() string {
	return "mall_shopping_cart"
}
func (this *ShoppingCartItem) TableIndex() [][]string {
	return [][]string{
		[]string{"UserId", "CorpId", "ProductId"},
	}
}





func init() {
	eel.RegisterModel(new(Count))
	eel.RegisterModel(new(PostageConfig))
	eel.RegisterModel(new(FreePostageConfig))
	eel.RegisterModel(new(SpecialAreaPostageConfig))
	eel.RegisterModel(new(ShoppingCartItem))
	eel.RegisterModel(new(ShipInfo))
	eel.RegisterModel(new(LimitZone))
}
